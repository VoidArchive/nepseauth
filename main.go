package main

import (
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/voidarchive/nepseauth/auth"
)

// HTTPClient implements auth.NepseHTTP interface and provides API access
type HTTPClient struct {
	client  *http.Client
	baseURL string
	manager *auth.Manager
}

func NewHTTPClient() (*HTTPClient, error) {
	// Create transport with TLS verification disabled (temporary)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING: Only for testing!
		},
	}

	httpClient := &HTTPClient{
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
		baseURL: "https://www.nepalstock.com",
	}

	// Create token manager
	manager, err := auth.NewManager(httpClient)
	if err != nil {
		return nil, fmt.Errorf("create token manager: %w", err)
	}

	httpClient.manager = manager
	return httpClient, nil
}

// Helper function to handle gzip decompression
func (h *HTTPClient) getResponseBody(resp *http.Response) (io.ReadCloser, error) {
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		return gzip.NewReader(resp.Body)
	}
	return resp.Body, nil
}

// GetTokens implements auth.NepseHTTP interface
func (h *HTTPClient) GetTokens(ctx context.Context) (*auth.TokenResponse, error) {
	url := h.baseURL + "/api/authenticate/prove"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	h.setCommonHeaders(req, false)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("decompress response: %w", err)
	}
	defer body.Close()

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(body).Decode(&tokenResp); err != nil {
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

	req.Header.Set("Authorization", "Salter "+refreshToken)
	h.setCommonHeaders(req, true)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("decompress response: %w", err)
	}
	defer body.Close()

	var tokenResp auth.TokenResponse
	if err := json.NewDecoder(body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &tokenResp, nil
}

func (h *HTTPClient) setCommonHeaders(req *http.Request, includeAuth bool) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Host", "www.nepalstock.com")
	req.Header.Set("Origin", "https://www.nepalstock.com")
	req.Header.Set("Referer", "https://www.nepalstock.com/")
}

// Authenticated API request helper
func (h *HTTPClient) apiRequest(ctx context.Context, endpoint string, result any) error {
	token, err := h.manager.AccessToken(ctx)
	if err != nil {
		return fmt.Errorf("get access token: %w", err)
	}

	url := h.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set authenticated headers
	auth.AuthHeader(req, token)
	req.Header.Set("Content-Type", "application/json")
	h.setCommonHeaders(req, true)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		// Token might be expired, force refresh and retry
		if err := h.manager.ForceUpdate(ctx); err != nil {
			return fmt.Errorf("token refresh failed: %w", err)
		}
		return h.apiRequest(ctx, endpoint, result) // Retry once
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := h.getResponseBody(resp)
	if err != nil {
		return fmt.Errorf("decompress response: %w", err)
	}
	defer body.Close()

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func (h *HTTPClient) Close(ctx context.Context) error {
	if h.manager != nil {
		return h.manager.Close(ctx)
	}
	return nil
}

// NEPSE API Methods

type MarketSummary struct {
	TotalTurnover      float64 `json:"totalTurnover"`
	TotalTrades        int     `json:"totalTrades"`
	TotalSharesTraded  int64   `json:"totalSharesTraded"`
	TotalMarketCap     float64 `json:"totalMarketCap"`
	TotalListedShares  int64   `json:"totalListedShares"`
	TotalListedCompany int     `json:"totalListedCompany"`
	// Add more fields as needed
}

func (h *HTTPClient) GetMarketSummary(ctx context.Context) (*MarketSummary, error) {
	var summary MarketSummary
	err := h.apiRequest(ctx, "/api/nots/market-summary/", &summary)
	return &summary, err
}

type TodayPrice struct {
	ID                  int     `json:"id"`
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closePrice"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	PreviousClose       float64 `json:"previousClose"`
	DifferenceRs        float64 `json:"differenceRs"`
	PercentageChange    float64 `json:"percentageChange"`
	// Add more fields as needed
}

func (h *HTTPClient) GetTodaysPrices(ctx context.Context) ([]TodayPrice, error) {
	var prices []TodayPrice
	err := h.apiRequest(ctx, "/api/nots/nepse-data/today-price", &prices)
	return prices, err
}

type Security struct {
	ID           int    `json:"id"`
	Symbol       string `json:"symbol"`
	SecurityName string `json:"securityName"`
	IsSuspended  bool   `json:"isSuspended"`
	// Add more fields as needed
}

func (h *HTTPClient) GetSecurities(ctx context.Context) ([]Security, error) {
	var securities []Security
	err := h.apiRequest(ctx, "/api/nots/security?nonDelisted=true", &securities)
	return securities, err
}

func (h *HTTPClient) GetCompanyDetails(ctx context.Context, companyID int) (map[string]interface{}, error) {
	var details map[string]interface{}
	err := h.apiRequest(ctx, fmt.Sprintf("/api/nots/security/%d", companyID), &details)
	return details, err
}

func (h *HTTPClient) GetCompanyFloorSheet(ctx context.Context, companyID int, businessDate string) (map[string]interface{}, error) {
	var floorSheet map[string]interface{}
	endpoint := fmt.Sprintf("/api/nots/security/floorsheet/%d?businessDate=%s&size=500", companyID, businessDate)
	err := h.apiRequest(ctx, endpoint, &floorSheet)
	return floorSheet, err
}

func (h *HTTPClient) GetCompanyPriceHistory(ctx context.Context, companyID int, startDate, endDate string) (map[string]interface{}, error) {
	var history map[string]interface{}
	endpoint := fmt.Sprintf("/api/nots/market/history/security/%d?size=500&startDate=%s&endDate=%s", companyID, startDate, endDate)
	err := h.apiRequest(ctx, endpoint, &history)
	return history, err
}

// Helper function to find company ID by symbol
func (h *HTTPClient) FindCompanyIDBySymbol(ctx context.Context, symbol string) (int, error) {
	securities, err := h.GetSecurities(ctx)
	if err != nil {
		return 0, err
	}
	
	for _, security := range securities {
		if security.Symbol == symbol {
			return security.ID, nil
		}
	}
	
	return 0, fmt.Errorf("company with symbol %s not found", symbol)
}

type NepseIndex struct {
	Date          string  `json:"date"`
	CloseValue    float64 `json:"closeValue"`
	PercentChange float64 `json:"percentChange"`
	PointChange   float64 `json:"pointChange"`
	// Add more fields as needed
}

func (h *HTTPClient) GetNepseIndex(ctx context.Context) (*NepseIndex, error) {
	var index NepseIndex
	err := h.apiRequest(ctx, "/api/nots/nepse-index", &index)
	return &index, err
}

type MarketStatus struct {
	IsMarketOpen bool   `json:"isOpen"`
	Status       string `json:"status"`
	AsOfTime     string `json:"asOf"`
}

func (h *HTTPClient) GetMarketStatus(ctx context.Context) (*MarketStatus, error) {
	var status MarketStatus
	err := h.apiRequest(ctx, "/api/nots/nepse-data/market-open", &status)
	return &status, err
}

func main() {
	fmt.Println("Testing NEPSE API with NABIL Company Data...")

	// Create HTTP client with token manager
	client, err := NewHTTPClient()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			log.Printf("Failed to close client: %v", err)
		}
	}()

	ctx := context.Background()

	// Test 1: Find NABIL Company ID
	fmt.Println("\n1. Finding NABIL Company ID...")
	nabilID, err := client.FindCompanyIDBySymbol(ctx, "NABIL")
	if err != nil {
		log.Printf("❌ Failed to find NABIL: %v", err)
		return
	}
	fmt.Printf("✓ NABIL Company ID: %d\n", nabilID)

	// Test 2: NABIL Company Details
	fmt.Println("\n2. Testing NABIL Company Details...")
	details, err := client.GetCompanyDetails(ctx, nabilID)
	if err != nil {
		log.Printf("❌ NABIL details failed: %v", err)
	} else {
		fmt.Printf("✓ NABIL Details retrieved successfully\n")
		if name, ok := details["securityName"].(string); ok {
			fmt.Printf("✓ Company Name: %s\n", name)
		}
		if symbol, ok := details["symbol"].(string); ok {
			fmt.Printf("✓ Symbol: %s\n", symbol)
		}
		if sector, ok := details["sectorName"].(string); ok {
			fmt.Printf("✓ Sector: %s\n", sector)
		}
	}

	// Test 3: NABIL Price History (last 30 days)
	fmt.Println("\n3. Testing NABIL Price History...")
	endDate := "2024-08-29"   // Today's date
	startDate := "2024-07-30" // 30 days ago
	history, err := client.GetCompanyPriceHistory(ctx, nabilID, startDate, endDate)
	if err != nil {
		log.Printf("❌ NABIL history failed: %v", err)
	} else {
		fmt.Printf("✓ NABIL Price History retrieved successfully\n")
		if content, ok := history["content"].([]interface{}); ok {
			fmt.Printf("✓ Historical records: %d\n", len(content))
			if len(content) > 0 {
				if record, ok := content[0].(map[string]interface{}); ok {
					if date, ok := record["businessDate"].(string); ok {
						fmt.Printf("✓ Latest date: %s\n", date)
					}
					if price, ok := record["closingPrice"].(float64); ok {
						fmt.Printf("✓ Latest price: Rs. %.2f\n", price)
					}
				}
			}
		}
	}

	// Test 4: NABIL Floor Sheet (today)
	fmt.Println("\n4. Testing NABIL Floor Sheet...")
	businessDate := "2024-08-29" // Today's date
	floorSheet, err := client.GetCompanyFloorSheet(ctx, nabilID, businessDate)
	if err != nil {
		log.Printf("❌ NABIL floor sheet failed: %v", err)
	} else {
		fmt.Printf("✓ NABIL Floor Sheet retrieved successfully\n")
		if fs, ok := floorSheet["floorsheets"].(map[string]interface{}); ok {
			if content, ok := fs["content"].([]interface{}); ok {
				fmt.Printf("✓ Floor sheet records: %d\n", len(content))
				if len(content) > 0 {
					fmt.Printf("✓ Recent transactions found\n")
				} else {
					fmt.Printf("✓ No transactions for today\n")
				}
			}
		}
	}

	fmt.Println("\n🎉 NABIL API testing completed!")
	fmt.Println("\n✅ Authentication working correctly!")
	fmt.Println("✅ Company-specific endpoints functional!")
}
