// Package nepse provides a modern, type-safe Go client for the NEPSE (Nepal Stock Exchange) API.
//
// This package offers comprehensive access to NEPSE market data including:
// - Market summaries and status information
// - Security and company listings
// - Real-time and historical price data
// - Trading volume and floor sheet information
// - Market indices and sub-indices
// - Top gainers, losers, and trading statistics
// - Supply and demand data
// - Market depth information
// - Graph data for various indices and securities
//
// The client is built with clean architecture principles, proper error handling,
// and type safety throughout. It provides both high-level convenience methods
// and low-level access to the underlying API endpoints.
//
// Example usage:
//
//	package main
//
//	import (
//		"context"
//		"fmt"
//		"log"
//		"time"
//
//		"github.com/voidarchive/nepseauth/nepse"
//	)
//
//	func main() {
//		// Create a new NEPSE client with default options
//		client, err := nepse.NewClient(nil)
//		if err != nil {
//			log.Fatalf("Failed to create NEPSE client: %v", err)
//		}
//		defer client.Close(context.Background())
//
//		ctx := context.Background()
//
//		// Get market summary
//		summary, err := client.GetMarketSummary(ctx)
//		if err != nil {
//			log.Fatalf("Failed to get market summary: %v", err)
//		}
//		fmt.Printf("Total Turnover: Rs. %.2f\n", summary.TotalTurnover)
//
//		// Find a company by symbol and get its details
//		security, err := client.FindSecurityBySymbol(ctx, "NABIL")
//		if err != nil {
//			log.Fatalf("Failed to find NABIL: %v", err)
//		}
//
//		details, err := client.GetCompanyDetails(ctx, security.ID)
//		if err != nil {
//			log.Fatalf("Failed to get company details: %v", err)
//		}
//		fmt.Printf("Company: %s, Sector: %s\n", details.SecurityName, details.SectorName)
//	}
package nepse

import (
	"context"
)

// NewClient creates a new NEPSE API client with the given options.
// If options is nil, default options will be used.
func NewClient(options *Options) (Client, error) {
	if options == nil {
		options = DefaultOptions()
	}
	
	return NewHTTPClient(options)
}

// NewClientWithDefaults creates a new NEPSE API client with default settings.
// This is a convenience function equivalent to NewClient(nil).
func NewClientWithDefaults() (Client, error) {
	return NewClient(nil)
}

// NewClientWithTLS creates a new NEPSE API client with TLS verification enabled/disabled.
// This is a convenience function for quick TLS configuration.
func NewClientWithTLS(tlsVerification bool) (Client, error) {
	options := DefaultOptions()
	options.TLSVerification = tlsVerification
	return NewClient(options)
}

// Version information
const (
	// Version is the current version of the nepse package
	Version = "1.0.0"
	
	// UserAgent is the default user agent string used by the client
	UserAgent = "nepse-go/" + Version
	
	// DefaultBaseURL is the default NEPSE API base URL
	DefaultBaseURL = "https://www.nepalstock.com"
)

// Predefined error instances for common error checking
var (
	// ErrTokenExpired can be used with errors.Is() to check for token expiration
	ErrTokenExpired = NewTokenExpiredError()
	
	// ErrNetworkError can be used with errors.Is() to check for network errors
	ErrNetworkError = NewNetworkError(nil)
	
	// ErrUnauthorized can be used with errors.Is() to check for authorization errors
	ErrUnauthorized = NewUnauthorizedError("unauthorized")
	
	// ErrNotFound can be used with errors.Is() to check for not found errors
	ErrNotFound = NewNotFoundError("resource")
	
	// ErrRateLimit can be used with errors.Is() to check for rate limit errors
	ErrRateLimit = NewRateLimitError()
)

// Common business date formats used by the NEPSE API
const (
	// DateFormat is the standard date format used by NEPSE API (YYYY-MM-DD)
	DateFormat = "2006-01-02"
	
	// DateTimeFormat is the standard datetime format used by NEPSE API
	DateTimeFormat = "2006-01-02 15:04:05"
)

// Sector names commonly used in the NEPSE market
const (
	SectorBanking            = "Banking"
	SectorDevelopmentBank    = "Development Bank"
	SectorFinance            = "Finance"
	SectorHotelTourism       = "Hotel Tourism"
	SectorHydro              = "Hydro"
	SectorInvestment         = "Investment"
	SectorLifeInsurance      = "Life Insurance"
	SectorManufacturing      = "Manufacturing"
	SectorMicrofinance       = "Microfinance"
	SectorMutualFund         = "Mutual Fund"
	SectorNonLifeInsurance   = "Non Life Insurance"
	SectorOthers             = "Others"
	SectorTrading            = "Trading"
	SectorPromoterShare      = "Promoter Share"
)

// BatchRequest represents a batch operation configuration
type BatchRequest struct {
	MaxConcurrency int           // Maximum number of concurrent requests
	Timeout        context.Context // Context for timeout control
}

// DefaultBatchRequest returns default batch request settings
func DefaultBatchRequest() *BatchRequest {
	return &BatchRequest{
		MaxConcurrency: 5, // Reasonable default to avoid overwhelming the server
	}
}