package nepse

// MarketSummaryItem represents a single item in the market summary response
type MarketSummaryItem struct {
	Detail string  `json:"detail"`
	Value  float64 `json:"value"`
}

// MarketSummary represents the processed market summary data
type MarketSummary struct {
	TotalTurnover             float64
	TotalTradedShares         float64
	TotalTransactions         float64
	TotalScripsTraded         float64
	TotalMarketCapitalization float64
	TotalFloatMarketCap       float64
}

// MarketStatus represents the current market status
type MarketStatus struct {
	IsOpen string `json:"isOpen"` // API returns "OPEN" or "CLOSE" string, not boolean
	AsOf   string `json:"asOf"`
	ID     int32  `json:"id"`
}

// IsMarketOpen returns true if the market is currently open
func (m *MarketStatus) IsMarketOpen() bool {
	return m.IsOpen == "OPEN"
}

// NepseIndexRaw represents the raw NEPSE index response item
type NepseIndexRaw struct {
	ID               int32   `json:"id"`
	Index            string  `json:"index"`
	Close            float64 `json:"close"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	PerChange        float64 `json:"perChange"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// NepseIndex represents the NEPSE main index (ID 58)
type NepseIndex struct {
	IndexValue       float64 `json:"close"`
	PercentChange    float64 `json:"perChange"`
	PointChange      float64 `json:"change"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// SubIndex represents a sector sub-index (uses same structure as NepseIndexRaw)
type SubIndex struct {
	ID               int32   `json:"id"`
	Index            string  `json:"index"`
	Close            float64 `json:"close"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	PreviousClose    float64 `json:"previousClose"`
	Change           float64 `json:"change"`
	PerChange        float64 `json:"perChange"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow"`
	CurrentValue     float64 `json:"currentValue"`
	GeneratedTime    string  `json:"generatedTime"`
}

// Security represents a listed security/company
type Security struct {
	ID                   int32  `json:"id"`
	Symbol               string `json:"symbol"`
	SecurityName         string `json:"securityName"`
	IsSuspended          bool   `json:"isSuspended"`
	SectorName           string `json:"sectorName"`
	Instrument           string `json:"instrument"`
	RegulatoryCategoryID int32  `json:"regulatoryCategoryId"`
	ShareGroupID         int32  `json:"shareGroupId"`
	ActiveStatus         string `json:"activeStatus"`
	ListingDate          string `json:"listingDate"`
}

// Company represents company information (different from Security)
type Company struct {
	ID                   int32   `json:"id"`
	Symbol               string  `json:"symbol"`
	SecurityName         string  `json:"securityName"`
	SectorName           string  `json:"sectorName"`
	MarketCapitalization float64 `json:"marketCapitalization"`
	ShareOutstanding     int64   `json:"shareOutstanding"`
	ShareOutstandingDate string  `json:"shareOutstandingDate"`
	ListedShares         int64   `json:"listedShares"`
	PaidUpValue          float64 `json:"paidUpValue"`
	IsPromoterListed     bool    `json:"isPromoterListed"`
	HasTradingPermission bool    `json:"hasTradingPermission"`
}

// TodayPrice represents today's price data for a security
type TodayPrice struct {
	ID                  int32   `json:"id"`
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
	TotalTrades         int32   `json:"totalTrades"`
	BusinessDate        string  `json:"businessDate"`
	SecurityID          int32   `json:"securityId"`
	LastTradedPrice     float64 `json:"lastTradedPrice"`
	MaxPrice            float64 `json:"maxPrice"`
	MinPrice            float64 `json:"minPrice"`
}

// PriceHistory represents historical price data for a security
type PriceHistory struct {
	BusinessDate        string  `json:"businessDate"`
	SecurityID          int32   `json:"securityId"`
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closingPrice"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	TotalTrades         int32   `json:"totalTrades"`
	PreviousClose       float64 `json:"previousClose"`
	DifferenceRs        float64 `json:"differenceRs"`
	PercentageChange    float64 `json:"percentageChange"`
}

// FloorSheetEntry represents a single floor sheet entry
type FloorSheetEntry struct {
	ContractID       int64   `json:"contractId"`
	StockSymbol      string  `json:"stockSymbol"`
	SecurityName     string  `json:"securityName"`
	BuyerMemberID    int32   `json:"buyerMemberId"`
	SellerMemberID   int32   `json:"sellerMemberId"`
	ContractQuantity int64   `json:"contractQuantity"`
	ContractRate     float64 `json:"contractRate"`
	BusinessDate     string  `json:"businessDate"`
	TradeTime        string  `json:"tradeTime"`
	SecurityID       int32   `json:"securityId"`
	ContractAmount   float64 `json:"contractAmount"`
	BuyerBrokerName  string  `json:"buyerBrokerName"`
	SellerBrokerName string  `json:"sellerBrokerName"`
	TradeBookID      int64   `json:"tradeBookId"`
}

// FloorSheetResponse represents the paginated floor sheet response
type FloorSheetResponse struct {
	FloorSheets struct {
		Content          []FloorSheetEntry `json:"content"`
		PageNumber       int32             `json:"number"`
		Size             int32             `json:"size"`
		TotalElements    int64             `json:"totalElements"`
		TotalPages       int32             `json:"totalPages"`
		First            bool              `json:"first"`
		Last             bool              `json:"last"`
		NumberOfElements int32             `json:"numberOfElements"`
	} `json:"floorsheets"`
}

// MarketDepth represents market depth information for a security
type MarketDepth struct {
	SecurityID   int32  `json:"securityId"`
	Symbol       string `json:"symbol"`
	SecurityName string `json:"securityName"`
	BuyDepth     []struct {
		Price    float64 `json:"price"`
		Quantity int64   `json:"quantity"`
		Orders   int32   `json:"orders"`
	} `json:"buyDepth"`
	SellDepth []struct {
		Price    float64 `json:"price"`
		Quantity int64   `json:"quantity"`
		Orders   int32   `json:"orders"`
	} `json:"sellDepth"`
}

// TopListEntry represents entries in top gainers/losers/trades lists
type TopListEntry struct {
	Symbol              string  `json:"symbol"`
	SecurityName        string  `json:"securityName"`
	ClosePrice          float64 `json:"closePrice"`
	PercentageChange    float64 `json:"percentageChange"`
	DifferenceRs        float64 `json:"differenceRs"`
	TotalTradedQuantity int64   `json:"totalTradedQuantity"`
	TotalTradedValue    float64 `json:"totalTradedValue"`
	TotalTrades         int32   `json:"totalTrades"`
	HighPrice           float64 `json:"highPrice,omitempty"`
	LowPrice            float64 `json:"lowPrice,omitempty"`
	OpenPrice           float64 `json:"openPrice,omitempty"`
	PreviousClose       float64 `json:"previousClose,omitempty"`
}

// SupplyDemandEntry represents supply and demand data
type SupplyDemandEntry struct {
	Symbol            string  `json:"symbol"`
	SecurityName      string  `json:"securityName"`
	SupplyQuantity    int64   `json:"supplyQuantity"`
	SupplyAmount      float64 `json:"supplyAmount"`
	DemandQuantity    int64   `json:"demandQuantity"`
	DemandAmount      float64 `json:"demandAmount"`
	DemandSupplyRatio float64 `json:"demandSupplyRatio"`
}

// GraphDataPoint represents a single data point in graph data
type GraphDataPoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

// GraphResponse represents graph data response
type GraphResponse struct {
	Data []GraphDataPoint `json:"data"`
}

// CompanyDetailsRaw represents the raw nested company details response
type CompanyDetailsRaw struct {
	SecurityMcsData struct {
		SecurityID          string  `json:"securityId"`
		OpenPrice           float64 `json:"openPrice"`
		HighPrice           float64 `json:"highPrice"`
		LowPrice            float64 `json:"lowPrice"`
		TotalTradeQuantity  int64   `json:"totalTradeQuantity"`
		TotalTrades         int32   `json:"totalTrades"`
		LastTradedPrice     float64 `json:"lastTradedPrice"`
		PreviousClose       float64 `json:"previousClose"`
		BusinessDate        string  `json:"businessDate"`
		ClosePrice          float64 `json:"closePrice"`
		FiftyTwoWeekHigh    float64 `json:"fiftyTwoWeekHigh"`
		FiftyTwoWeekLow     float64 `json:"fiftyTwoWeekLow"`
		LastUpdatedDateTime string  `json:"lastUpdatedDateTime"`
	} `json:"securityMcsData"`
	SecurityData struct {
		ID               int32  `json:"id"`
		Symbol           string `json:"symbol"`
		SecurityName     string `json:"securityName"`
		ActiveStatus     string `json:"activeStatus"`
		PermittedToTrade string `json:"permittedToTrade"`
		Email            string `json:"email"`
		Sector           string `json:"sector"`
	} `json:"securityData"`
}

// CompanyDetails represents processed company information
type CompanyDetails struct {
	ID               int32  `json:"id"`
	Symbol           string `json:"symbol"`
	SecurityName     string `json:"securityName"`
	SectorName       string `json:"sectorName"`
	Email            string `json:"email"`
	ActiveStatus     string `json:"activeStatus"`
	PermittedToTrade string `json:"permittedToTrade"`

	// Market data fields
	OpenPrice           float64 `json:"openPrice"`
	HighPrice           float64 `json:"highPrice"`
	LowPrice            float64 `json:"lowPrice"`
	ClosePrice          float64 `json:"closePrice"`
	LastTradedPrice     float64 `json:"lastTradedPrice"`
	PreviousClose       float64 `json:"previousClose"`
	TotalTradeQuantity  int64   `json:"totalTradeQuantity"`
	TotalTrades         int32   `json:"totalTrades"`
	FiftyTwoWeekHigh    float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow     float64 `json:"fiftyTwoWeekLow"`
	BusinessDate        string  `json:"businessDate"`
	LastUpdatedDateTime string  `json:"lastUpdatedDateTime"`
}

// LiveMarketEntry represents live market data entry
type LiveMarketEntry struct {
	Symbol           string  `json:"symbol"`
	SecurityName     string  `json:"securityName"`
	OpenPrice        float64 `json:"openPrice"`
	HighPrice        float64 `json:"highPrice"`
	LowPrice         float64 `json:"lowPrice"`
	ClosePrice       float64 `json:"closePrice"`
	PercentChange    float64 `json:"percentChange"`
	Volume           int64   `json:"volume"`
	PreviousClose    float64 `json:"previousClose"`
	LastTradedVolume int64   `json:"lastTradedVolume"`
}

// SectorScrips represents scrips grouped by sector
type SectorScrips map[string][]string

// PaginatedResponse represents a generic paginated response
type PaginatedResponse[T any] struct {
	Content          []T   `json:"content"`
	PageNumber       int32 `json:"number"`
	Size             int32 `json:"size"`
	TotalElements    int64 `json:"totalElements"`
	TotalPages       int32 `json:"totalPages"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	NumberOfElements int32 `json:"numberOfElements"`
}

