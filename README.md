# NEPSE Go Library

A modern, type-safe Go client library for the NEPSE (Nepal Stock Exchange) API. This library provides comprehensive access to NEPSE market data with clean architecture, proper error handling, and full type safety.

## Features

- ✅ **Clean Architecture** - Modular design with clear separation of concerns
- ✅ **Broad Coverage** - Major NEPSE market data endpoints supported
- ✅ **Automatic Authentication** - Token management handled transparently
- ✅ **Retry Logic** - Built-in retry with exponential backoff
- ✅ **Context Support** - Full context.Context support for cancellation and timeouts
- ✅ **Error Handling** - Structured error types with proper error chains
- ✅ **Convenience Methods** - High-level methods for common operations

## Installation

```bash
go get github.com/voidarchive/nepseauth
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/voidarchive/nepseauth/nepse"
)

func main() {
    // Create a new NEPSE client
    client, err := nepse.NewClientWithDefaults()
    if err != nil {
        log.Fatalf("Failed to create NEPSE client: %v", err)
    }
    defer client.Close(context.Background())

    ctx := context.Background()

    // Get market summary
    summary, err := client.GetMarketSummary(ctx)
    if err != nil {
        log.Fatalf("Failed to get market summary: %v", err)
    }
    
    fmt.Printf("Total Turnover: Rs. %.2f\\n", summary.TotalTurnover)
    fmt.Printf("Total Transactions: %.0f\\n", summary.TotalTransactions)
    
    // Find a company by symbol
    security, err := client.FindSecurityBySymbol(ctx, "NABIL")
    if err != nil {
        log.Fatalf("Failed to find NABIL: %v", err)
    }
    
    // Get company details - can use either ID or symbol
    details, err := client.GetCompanyDetailsBySymbol(ctx, "NABIL")
    if err != nil {
        log.Fatalf("Failed to get company details: %v", err)
    }
    
    fmt.Printf("Company: %s, Sector: %s, LTP: Rs. %.2f\\n", 
        details.SecurityName, details.SectorName, details.LastTradedPrice)
}
```

## API Coverage

### Market Data

- `GetMarketSummary()` - Overall market statistics
- `GetMarketStatus()` - Current market open/close status
- `GetNepseIndex()` - NEPSE main index information
- `GetNepseSubIndices()` - All sector sub-indices
- `GetLiveMarket()` - Live market data
- `GetSupplyDemand()` - Supply and demand information

### Securities & Companies

- `GetSecurityList()` - All listed securities
- `GetCompanyList()` - All listed companies
- `GetCompanyDetails(securityID)` / `GetCompanyDetailsBySymbol(symbol)` - Detailed company information
- `GetSectorScrips()` - Securities grouped by sector (fast, no API calls needed)
- `FindSecurity(securityID)` / `FindSecurityBySymbol(symbol)` - Find security by ID or symbol

### Price & Trading Data

- `GetTodaysPrices(businessDate)` - Today's price data
- `GetPriceVolumeHistory(securityID, startDate, endDate)` / `GetPriceVolumeHistoryBySymbol(symbol, startDate, endDate)` - Historical prices
- `GetMarketDepth(securityID)` / `GetMarketDepthBySymbol(symbol)` - Market depth information
- `GetFloorSheet()` - Complete floor sheet data
- `GetFloorSheetOf(securityID, businessDate)` / `GetFloorSheetBySymbol(symbol, businessDate)` - Company-specific floor sheet

### Top Lists

- `GetTopGainers()` - Top gaining securities
- `GetTopLosers()` - Top losing securities  
- `GetTopTenTrade()` - Top by trade volume
- `GetTopTenTransaction()` - Top by transaction count
- `GetTopTenTurnover()` - Top by turnover

### Graph Data (Technical Analysis) - **⚠️ Currently Non-Functional**

**Note**: Graph endpoints currently return empty data due to NEPSE API backend issues.

- `GetDailyNepseIndexGraph()` - NEPSE index chart data
- `GetDailySensitiveIndexGraph()` - Sensitive index chart
- `GetDailyFloatIndexGraph()` - Float index chart
- `GetDailyScripPriceGraph(securityID)` / `GetDailyScripPriceGraphBySymbol(symbol)` - Individual security chart

### Sector Sub-Index Graphs

- `GetDailyBankSubindexGraph()` - Banking sector index
- `GetDailyFinanceSubindexGraph()` - Finance sector index
- `GetDailyHydroSubindexGraph()` - Hydro sector index
- And many more...

## Configuration Options

```go
options := &nepse.Options{
    BaseURL:         "https://www.nepalstock.com",
    TLSVerification: true,  // Note: May need to be false due to NEPSE server TLS issues
    HTTPTimeout:     30 * time.Second,
    MaxRetries:      3,
    RetryDelay:      time.Second,
}

client, err := nepse.NewClient(options)
```

### Important Security Note

The `TLSVerification: false` option exists due to TLS configuration issues on NEPSE's servers (nepalstock.com). This is a known limitation of the NEPSE API infrastructure, not the client library. When NEPSE fixes their TLS configuration, always use `TLSVerification: true` for production deployments.

## Error Handling

The library provides structured error types for better error handling:

```go
import "errors"

data, err := client.GetMarketSummary(ctx)
if err != nil {
    var nepseErr *nepse.NepseError
    if errors.As(err, &nepseErr) {
        switch nepseErr.Type {
        case nepse.ErrorTypeTokenExpired:
            // Handle token expiration
        case nepse.ErrorTypeNetworkError:
            // Handle network issues
        case nepse.ErrorTypeRateLimit:
            // Handle rate limiting
        }
    }
}
```

## Architecture

The library is organized into several packages:

- **`nepse`** - Main package with client interface and factory functions
- **`nepse/config.go`** - Configuration and static data
- **`nepse/errors.go`** - Structured error types
- **`nepse/types.go`** - All API response types
- **`nepse/client.go`** - Client interface definition
- **`nepse/http_client.go`** - HTTP client implementation
- **`nepse/market_data.go`** - GET API methods
- **`nepse/graphs.go`** - GET API methods for graph data

## Key Differences from Python Version

1. **Type Safety** - All responses are properly typed structs instead of generic maps
2. **Clean Architecture** - Separated into logical modules instead of one giant file
3. **Modern Error Handling** - Structured errors with proper error chaining
4. **Context Support** - Full context.Context support throughout
5. **Interface-based Design** - Clean interfaces for testing and extensibility
6. **Resource Management** - Proper connection pooling and cleanup

## Examples

See `cmd/examples` for usage examples:

- `cmd/examples/basic_usage.go` - Basic API usage examples
- `main.go` - Simple NABIL summary demo

## Testing

### Against Real NEPSE API

```bash
# Run the main test with NABIL data
go run main.go

# Run the basic usage example
go run cmd/examples/basic_usage.go
```

### Notes

- Real API calls require valid NEPSE tokens. The library handles token mint/refresh using an embedded WASM parser.
- A separate mock server and OpenAPI spec are not included in this repository at the moment.

## API Design Philosophy

This library follows modern Go best practices:

### Type-Safe Dual API Pattern

Most methods that work with securities provide both ID and symbol variants:

```go
// Primary methods (faster, direct)
client.GetMarketDepth(ctx, int32(131))              // By ID
client.GetCompanyDetails(ctx, int32(131))           // By ID

// Convenience methods (easier to use)
client.GetMarketDepthBySymbol(ctx, "NABIL")        // By symbol  
client.GetCompanyDetailsBySymbol(ctx, "NABIL")     // By symbol
```

### Performance Optimizations

- **Fast Sector Grouping**: `GetSectorScrips()` uses existing data, no API calls needed
- **Connection Pooling**: Efficient HTTP connection reuse
- **Smart Retry Logic**: Exponential backoff with circuit breaking

## Contributing

1. Follow Go best practices and conventions
2. Maintain type safety - no `interface{}` or `any` in public APIs  
3. Add proper error handling for all new features
4. Include examples for new functionality
5. Update documentation for API changes
6. **Security First**: All contributions must pass security review

## Go Version

- Requires Go 1.23+

## License

This project is licensed under the MIT License - see the LICENSE file for details.
