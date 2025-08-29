// Package auth: Authenticate Nepse API
package auth

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// -----------------------------------------------------------------------------
// Embeds the WASM used to compute indices from salts (equivalent to css.wasm).
//
//go:embed data/css.wasm
var cssWasm []byte

// -----------------------------------------------------------------------------
// Public types & ports

// NepseHTTP abstracts the token GET request your infra already has.
type NepseHTTP interface {
	// GetTokens performs GET to /api/authenticate/prove and returns the token response
	GetTokens(ctx context.Context) (*TokenResponse, error)

	// RefreshTokens performs GET to /api/authenticate/refresh-token and returns new tokens
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenResponse, error)
}

// TokenResponse mirrors the JSON from /api/authenticate/prove.
type TokenResponse struct {
	Salt1        int    `json:"salt1"`
	Salt2        int    `json:"salt2"`
	Salt3        int    `json:"salt3"`
	Salt4        int    `json:"salt4"`
	Salt5        int    `json:"salt5"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ServerTime   int64  `json:"serverTime"` // ms since epoch
}

// Manager manages NEPSE auth tokens like the Python TokenManager/AsyncTokenManager.
type Manager struct {
	http NepseHTTP

	parser *tokenParser

	maxUpdatePeriod time.Duration

	mu           sync.RWMutex
	accessToken  string
	refreshToken string
	tokenTS      time.Time
	salts        [5]int

	sf singleflight.Group
}

// NewManager constructs a Manager. It loads and initializes the embedded WASM parser once.
func NewManager(httpClient NepseHTTP) (*Manager, error) {
	parser, err := newTokenParser()
	if err != nil {
		return nil, fmt.Errorf("init wasm parser: %w", err)
	}
	return &Manager{
		http:            httpClient,
		parser:          parser,
		maxUpdatePeriod: 45 * time.Second,
	}, nil
}

// Close releases WASM runtime resources.
func (m *Manager) Close(ctx context.Context) error {
	if m.parser != nil {
		return m.parser.close(ctx)
	}
	return nil
}

// AccessToken returns a valid access token, refreshing if needed.
func (m *Manager) AccessToken(ctx context.Context) (string, error) {
	if m.isValid() {
		m.mu.RLock()
		t := m.accessToken
		m.mu.RUnlock()
		return t, nil
	}
	if err := m.update(ctx); err != nil {
		return "", err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" {
		return "", errors.New("empty access token after update")
	}
	return m.accessToken, nil
}

// RefreshToken returns a valid refresh token, refreshing if needed.
func (m *Manager) RefreshToken(ctx context.Context) (string, error) {
	if m.isValid() {
		m.mu.RLock()
		t := m.refreshToken
		m.mu.RUnlock()
		return t, nil
	}
	if err := m.update(ctx); err != nil {
		return "", err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.refreshToken == "" {
		return "", errors.New("empty refresh token after update")
	}
	return m.refreshToken, nil
}

// ForceUpdate forces a token refresh (rarely needed).
func (m *Manager) ForceUpdate(ctx context.Context) error {
	return m.update(ctx)
}

// -----------------------------------------------------------------------------
// Internals

func (m *Manager) isValid() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.accessToken == "" || m.tokenTS.IsZero() {
		return false
	}
	return time.Since(m.tokenTS) < m.maxUpdatePeriod
}

type updateResult struct{} // Empty struct as we only care about success/error

func (m *Manager) update(ctx context.Context) error {
	// Deduplicate concurrent refreshes.
	_, err, _ := m.sf.Do("token_update", func() (any, error) {
		// Re-check validity inside singleflight window.
		if m.isValid() {
			return updateResult{}, nil
		}

		resp, err := m.http.GetTokens(ctx)
		if err != nil {
			return nil, fmt.Errorf("get token: %w", err)
		}

		access, refresh, salts, ts, err := m.parseResponse(*resp)
		if err != nil {
			return nil, err
		}

		m.mu.Lock()
		m.accessToken = access
		m.refreshToken = refresh
		m.salts = salts
		if ts > 0 {
			// Python used int(serverTime/1000). We'll keep seconds precision.
			m.tokenTS = time.Unix(ts, 0)
		} else {
			m.tokenTS = time.Now()
		}
		m.mu.Unlock()

		return updateResult{}, nil
	})
	return err
}

func (m *Manager) parseResponse(tr TokenResponse) (string, string, [5]int, int64, error) {
	salts := [5]int{tr.Salt1, tr.Salt2, tr.Salt3, tr.Salt4, tr.Salt5}
	// Compute indices via WASM, mirroring the Python order and functions.
	idx, err := m.parser.indicesFromSalts(salts)
	if err != nil {
		return "", "", salts, 0, fmt.Errorf("wasm parse: %w", err)
	}

	// Apply the same slicing logic as Python (assumes indices are ascending).
	parsedAccess := sliceSkipAt(tr.AccessToken, idx.access...)
	parsedRefresh := sliceSkipAt(tr.RefreshToken, idx.refresh...)

	// Server time is in ms; Python uses seconds.
	sec := tr.ServerTime / 1000
	return parsedAccess, parsedRefresh, salts, sec, nil
}

// sliceSkipAt reproduces:
// s[0:n] + s[n+1:l] + s[l+1:o] + s[o+1:p] + s[p+1:q] + s[q+1:]
// If indices are out of order, it sorts them first for safety.
func sliceSkipAt(s string, positions ...int) string {
	if len(positions) == 0 {
		return s
	}
	// Defensive copy and sort ascending
	ps := make([]int, len(positions))
	copy(ps, positions)
	for i := 1; i < len(ps); i++ {
		j := i
		for j > 0 && ps[j-1] > ps[j] {
			ps[j-1], ps[j] = ps[j], ps[j-1]
			j--
		}
	}

	// Work on bytes (tokens are ASCII/base64-like). If unsure, convert to []rune.
	b := []byte(s)
	var out []byte
	prev := 0
	for _, p := range ps {
		if p < 0 || p >= len(b) {
			continue // ignore bad indices gracefully
		}
		out = append(out, b[prev:p]...)
		prev = p + 1
	}
	out = append(out, b[prev:]...)
	return string(out)
}

// -----------------------------------------------------------------------------
// WASM plumbing

type tokenParser struct {
	rt  wazero.Runtime
	mod api.Module
	cdx api.Function
	rdx api.Function
	bdx api.Function
	ndx api.Function
	mdx api.Function
}

func newTokenParser() (*tokenParser, error) {
	ctx := context.Background()
	rt := wazero.NewRuntime(ctx)

	compiled, err := rt.CompileModule(ctx, cssWasm)
	if err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("compile wasm: %w", err)
	}
	mod, err := rt.InstantiateModule(ctx, compiled, wazero.NewModuleConfig())
	if err != nil {
		_ = rt.Close(ctx)
		return nil, fmt.Errorf("instantiate wasm: %w", err)
	}

	getExport := func(name string) (api.Function, error) {
		f := mod.ExportedFunction(name)
		if f == nil {
			return nil, fmt.Errorf("export %q not found", name)
		}
		return f, nil
	}

	cdx, err := getExport("cdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	rdx, err := getExport("rdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	bdx, err := getExport("bdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	ndx, err := getExport("ndx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}
	mdx, err := getExport("mdx")
	if err != nil {
		_ = rt.Close(ctx)
		return nil, err
	}

	return &tokenParser{
		rt:  rt,
		mod: mod,
		cdx: cdx, rdx: rdx, bdx: bdx, ndx: ndx, mdx: mdx,
	}, nil
}

func (p *tokenParser) close(ctx context.Context) error {
	return p.rt.Close(ctx)
}

type tokenIndices struct {
	access  []int // n, l, o, p, q
	refresh []int // a, b, c, d, e
}

func (p *tokenParser) indicesFromSalts(s [5]int) (tokenIndices, error) {
	ctx := context.Background()

	// Helper to call i32 functions (wazero returns uint64 slots).
	call5 := func(f api.Function, a, b, c, d, e int) (int, error) {
		res, err := f.Call(ctx,
			uint64(uint32(a)), uint64(uint32(b)),
			uint64(uint32(c)), uint64(uint32(d)), uint64(uint32(e)),
		)
		if err != nil {
			return 0, err
		}
		return int(int32(res[0])), nil
	}

	// Python:
	// n = cdx(s1,s2,s3,s4,s5)
	// l = rdx(s1,s2,s4,s3,s5)
	// o = bdx(s1,s2,s4,s3,s5)
	// p = ndx(s1,s2,s4,s3,s5)
	// q = mdx(s1,s2,s4,s3,s5)
	s1, s2, s3, s4, s5 := s[0], s[1], s[2], s[3], s[4]
	n, err := call5(p.cdx, s1, s2, s3, s4, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	l, err := call5(p.rdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	o, err := call5(p.bdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	pp, err := call5(p.ndx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	q, err := call5(p.mdx, s1, s2, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}

	// Python:
	// a = cdx(s2,s1,s3,s5,s4)
	// b = rdx(s2,s1,s3,s4,s5)
	// c = bdx(s2,s1,s4,s3,s5)
	// d = ndx(s2,s1,s4,s3,s5)
	// e = mdx(s2,s1,s4,s3,s5)
	a, err := call5(p.cdx, s2, s1, s3, s5, s4)
	if err != nil {
		return tokenIndices{}, err
	}
	b, err := call5(p.rdx, s2, s1, s3, s4, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	cc, err := call5(p.bdx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	d, err := call5(p.ndx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}
	e, err := call5(p.mdx, s2, s1, s4, s3, s5)
	if err != nil {
		return tokenIndices{}, err
	}

	return tokenIndices{
		access:  []int{n, l, o, pp, q},
		refresh: []int{a, b, cc, d, e},
	}, nil
}

// Optional: helper to inject Authorization header into your other requests.

func AuthHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", "Salter "+token)
}
