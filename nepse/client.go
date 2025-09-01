package nepse

import (
	"context"
	"time"
)

// Client defines the interface for NEPSE API operations
type Client interface {
	// Market Data Methods
	GetMarketSummary(ctx context.Context) (*MarketSummary, error)
	GetMarketStatus(ctx context.Context) (*MarketStatus, error)
	GetNepseIndex(ctx context.Context) (*NepseIndex, error)
	GetNepseSubIndices(ctx context.Context) ([]SubIndex, error)
	GetLiveMarket(ctx context.Context) ([]LiveMarketEntry, error)

	// Security and Company Methods
	GetSecurityList(ctx context.Context) ([]Security, error)
	GetCompanyList(ctx context.Context) ([]Company, error)
	GetCompanyDetails(ctx context.Context, securityID int32) (*CompanyDetails, error)
	GetSectorScrips(ctx context.Context) (SectorScrips, error)

	// Price and Trading Data
	GetTodaysPrices(ctx context.Context, businessDate string) ([]TodayPrice, error)
	GetPriceVolumeHistory(ctx context.Context, securityID int32, startDate, endDate string) ([]PriceHistory, error)
	GetSupplyDemand(ctx context.Context) ([]SupplyDemandEntry, error)
	GetMarketDepth(ctx context.Context, securityID int32) (*MarketDepth, error)

	// Top Lists
	GetTopGainers(ctx context.Context) ([]TopListEntry, error)
	GetTopLosers(ctx context.Context) ([]TopListEntry, error)
	GetTopTenTrade(ctx context.Context) ([]TopListEntry, error)
	GetTopTenTransaction(ctx context.Context) ([]TopListEntry, error)
	GetTopTenTurnover(ctx context.Context) ([]TopListEntry, error)

	// Floor Sheet
	GetFloorSheet(ctx context.Context) ([]FloorSheetEntry, error)
	GetFloorSheetOf(ctx context.Context, securityID int32, businessDate string) ([]FloorSheetEntry, error)

	// Graph Data (GET endpoints)
	GetDailyNepseIndexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailySensitiveIndexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyFloatIndexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailySensitiveFloatIndexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyScripPriceGraph(ctx context.Context, securityID int32) (*GraphResponse, error)
	GetDailyScripPriceGraphBySymbol(ctx context.Context, symbol string) (*GraphResponse, error)

	// Sub-Index Graphs
	GetDailyBankSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyDevelopmentBankSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyFinanceSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyHotelTourismSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyHydroSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyInvestmentSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyManufacturingSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyMicrofinanceSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyMutualfundSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyNonLifeInsuranceSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyOthersSubindexGraph(ctx context.Context) (*GraphResponse, error)
	GetDailyTradingSubindexGraph(ctx context.Context) (*GraphResponse, error)

	// Helper Methods
	FindSecurityBySymbol(ctx context.Context, symbol string) (*Security, error)
	FindCompanyBySymbol(ctx context.Context, symbol string) (*Company, error)
	GetFloorSheetBySymbol(ctx context.Context, symbol string, businessDate string) ([]FloorSheetEntry, error)
	GetPriceVolumeHistoryBySymbol(ctx context.Context, symbol string, startDate, endDate string) ([]PriceHistory, error)
	
	// Configuration
	SetTLSVerification(enabled bool)
	GetConfig() *Config
	
	// Lifecycle
	Close(ctx context.Context) error
}

// Options represents configuration options for creating a new NEPSE client
type Options struct {
	// BaseURL overrides the default NEPSE API base URL
	BaseURL string
	
	// TLSVerification enables/disables TLS certificate verification
	TLSVerification bool
	
	// HTTPTimeout sets the HTTP request timeout
	HTTPTimeout time.Duration
	
	// MaxRetries sets the maximum number of retries for failed requests
	MaxRetries int
	
	// RetryDelay sets the base delay between retries
	RetryDelay time.Duration
	
	// Config overrides the default configuration
	Config *Config
}

// DefaultOptions returns default options for the NEPSE client
func DefaultOptions() *Options {
	return &Options{
		BaseURL:         "https://www.nepalstock.com",
		TLSVerification: true,
		HTTPTimeout:     30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      time.Second,
		Config:          DefaultConfig(),
	}
}