package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/voidarchive/nepseauth/auth" // Replace with your actual module path
)

// HTTPClient implements auth.NepseHTTP interface
type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPClient() *HTTPClient {
	// Create transport with TLS verification disabled
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING: Only for testing!
		},
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		baseURL: "https://nepalstock.com.np", // Verify this is correct
	}
}

// GetTokens implements auth.NepseHTTP interface
func (h *HTTPClient) GetTokens(ctx context.Context) (*auth.TokenResponse, error) {
	url := h.baseURL + "/api/authenticate/prove"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers to mimic browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshTokens implements auth.NepseHTTP interface
func (h *HTTPClient) RefreshTokens(ctx context.Context, refreshToken string) (*auth.TokenResponse, error) {
	url := h.baseURL + "/api/authenticate/refresh-token"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Add refresh token to Authorization header
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &tokenResp, nil
}

func main() {
	fmt.Println("Testing NEPSE Authentication...")

	// Create HTTP client
	httpClient := NewHTTPClient()

	// Create token manager
	manager, err := auth.NewManager(httpClient)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}
	defer func() {
		if err := manager.Close(context.Background()); err != nil {
			log.Printf("Failed to close manager: %v", err)
		}
	}()

	ctx := context.Background()

	fmt.Println("\n1. Testing initial token fetch...")
	accessToken, err := manager.AccessToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	fmt.Printf("✓ Access Token (first 20 chars): %s...\n", accessToken[:min(20, len(accessToken))])

	refreshToken, err := manager.RefreshToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get refresh token: %v", err)
	}

	fmt.Printf("✓ Refresh Token (first 20 chars): %s...\n", refreshToken[:min(20, len(refreshToken))])

	fmt.Println("\n2. Testing token caching (should return same tokens)...")
	accessToken2, err := manager.AccessToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get cached access token: %v", err)
	}

	if accessToken == accessToken2 {
		fmt.Println("✓ Token caching working correctly")
	} else {
		fmt.Println("⚠ Tokens don't match - caching may not be working")
	}

	fmt.Println("\n3. Testing force update...")
	if err := manager.ForceUpdate(ctx); err != nil {
		log.Fatalf("Failed to force update: %v", err)
	}

	accessToken3, err := manager.AccessToken(ctx)
	if err != nil {
		log.Fatalf("Failed to get token after force update: %v", err)
	}

	if accessToken != accessToken3 {
		fmt.Println("✓ Force update working - got new token")
	} else {
		fmt.Println("⚠ Force update may not be working - same token returned")
	}

	fmt.Printf("✓ New Access Token (first 20 chars): %s...\n", accessToken3[:min(20, len(accessToken3))])

	fmt.Println("\n4. Testing with timeout context...")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	start := time.Now()
	_, err = manager.AccessToken(timeoutCtx)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("⚠ Request failed (took %v): %v\n", elapsed, err)
	} else {
		fmt.Printf("✓ Request completed in %v\n", elapsed)
	}

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("- Use the access token to make authenticated requests to NEPSE API")
	fmt.Println("- Test token refresh mechanism")
	fmt.Println("- Integrate with your NTX market data fetching")
}

// min helper function for older Go versions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
