package nepse

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Market Data GET API Methods

// GetMarketSummary retrieves the overall market summary
func (h *HTTPClient) GetMarketSummary(ctx context.Context) (*MarketSummary, error) {
	var rawItems []MarketSummaryItem
	err := h.apiRequest(ctx, h.config.APIEndpoints["market_summary"], &rawItems)
	if err != nil {
		return nil, fmt.Errorf("failed to get market summary: %w", err)
	}

	// Convert array response to structured format
	summary := &MarketSummary{}
	for _, item := range rawItems {
		switch item.Detail {
		case "Total Turnover Rs:":
			summary.TotalTurnover = item.Value
		case "Total Traded Shares":
			summary.TotalTradedShares = item.Value
		case "Total Transactions":
			summary.TotalTransactions = item.Value
		case "Total Scrips Traded":
			summary.TotalScripsTraded = item.Value
		case "Total Market Capitalization Rs:":
			summary.TotalMarketCapitalization = item.Value
		case "Total Float Market Capitalization Rs:":
			summary.TotalFloatMarketCap = item.Value
		}
	}

	return summary, nil
}

// GetMarketStatus retrieves the current market status
func (h *HTTPClient) GetMarketStatus(ctx context.Context) (*MarketStatus, error) {
	var status MarketStatus
	err := h.apiRequest(ctx, h.config.APIEndpoints["market_open"], &status)
	if err != nil {
		return nil, fmt.Errorf("failed to get market status: %w", err)
	}
	return &status, nil
}

// GetNepseIndex retrieves the NEPSE index information
func (h *HTTPClient) GetNepseIndex(ctx context.Context) (*NepseIndex, error) {
	var rawIndices []NepseIndexRaw
	err := h.apiRequest(ctx, h.config.APIEndpoints["nepse_index"], &rawIndices)
	if err != nil {
		return nil, fmt.Errorf("failed to get NEPSE index: %w", err)
	}

	// Find the main NEPSE index (ID 58)
	for _, rawIndex := range rawIndices {
		if rawIndex.ID == 58 && rawIndex.Index == "NEPSE Index" {
			return &NepseIndex{
				IndexValue:       rawIndex.Close,
				PercentChange:    rawIndex.PerChange,
				PointChange:      rawIndex.Change,
				High:             rawIndex.High,
				Low:              rawIndex.Low,
				PreviousClose:    rawIndex.PreviousClose,
				FiftyTwoWeekHigh: rawIndex.FiftyTwoWeekHigh,
				FiftyTwoWeekLow:  rawIndex.FiftyTwoWeekLow,
				CurrentValue:     rawIndex.CurrentValue,
				GeneratedTime:    rawIndex.GeneratedTime,
			}, nil
		}
	}

	return nil, NewNotFoundError("NEPSE Index")
}

// GetNepseSubIndices retrieves all NEPSE sub-indices
func (h *HTTPClient) GetNepseSubIndices(ctx context.Context) ([]SubIndex, error) {
	var rawIndices []NepseIndexRaw
	err := h.apiRequest(ctx, h.config.APIEndpoints["nepse_index"], &rawIndices)
	if err != nil {
		return nil, fmt.Errorf("failed to get NEPSE sub-indices: %w", err)
	}

	// Filter out the main indices and return sub-indices
	var subIndices []SubIndex
	for _, rawIndex := range rawIndices {
		// Skip main indices (58=NEPSE, 57=Sensitive, 62=Float, 63=Sensitive Float)
		if rawIndex.ID != 58 && rawIndex.ID != 57 && rawIndex.ID != 62 && rawIndex.ID != 63 {
			subIndices = append(subIndices, SubIndex{
				ID:               rawIndex.ID,
				Index:            rawIndex.Index,
				Close:            rawIndex.Close,
				High:             rawIndex.High,
				Low:              rawIndex.Low,
				PreviousClose:    rawIndex.PreviousClose,
				Change:           rawIndex.Change,
				PerChange:        rawIndex.PerChange,
				FiftyTwoWeekHigh: rawIndex.FiftyTwoWeekHigh,
				FiftyTwoWeekLow:  rawIndex.FiftyTwoWeekLow,
				CurrentValue:     rawIndex.CurrentValue,
				GeneratedTime:    rawIndex.GeneratedTime,
			})
		}
	}

	return subIndices, nil
}

// GetLiveMarket retrieves live market data
func (h *HTTPClient) GetLiveMarket(ctx context.Context) ([]LiveMarketEntry, error) {
	var liveMarket []LiveMarketEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["live_market"], &liveMarket)
	if err != nil {
		return nil, fmt.Errorf("failed to get live market data: %w", err)
	}
	return liveMarket, nil
}

// GetSupplyDemand retrieves supply and demand data
func (h *HTTPClient) GetSupplyDemand(ctx context.Context) ([]SupplyDemandEntry, error) {
	var supplyDemand []SupplyDemandEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["supply_demand"], &supplyDemand)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply demand data: %w", err)
	}
	return supplyDemand, nil
}

// Top Lists Methods

// GetTopGainers retrieves the top gainers list
func (h *HTTPClient) GetTopGainers(ctx context.Context) ([]TopListEntry, error) {
	var topGainers []TopListEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["top_gainers"], &topGainers)
	if err != nil {
		return nil, fmt.Errorf("failed to get top gainers: %w", err)
	}
	return topGainers, nil
}

// GetTopLosers retrieves the top losers list
func (h *HTTPClient) GetTopLosers(ctx context.Context) ([]TopListEntry, error) {
	var topLosers []TopListEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["top_losers"], &topLosers)
	if err != nil {
		return nil, fmt.Errorf("failed to get top losers: %w", err)
	}
	return topLosers, nil
}

// GetTopTenTrade retrieves the top ten trade list
func (h *HTTPClient) GetTopTenTrade(ctx context.Context) ([]TopListEntry, error) {
	var topTrade []TopListEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_trade"], &topTrade)
	if err != nil {
		return nil, fmt.Errorf("failed to get top ten trade: %w", err)
	}
	return topTrade, nil
}

// GetTopTenTransaction retrieves the top ten transaction list
func (h *HTTPClient) GetTopTenTransaction(ctx context.Context) ([]TopListEntry, error) {
	var topTransaction []TopListEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_transaction"], &topTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to get top ten transaction: %w", err)
	}
	return topTransaction, nil
}

// GetTopTenTurnover retrieves the top ten turnover list
func (h *HTTPClient) GetTopTenTurnover(ctx context.Context) ([]TopListEntry, error) {
	var topTurnover []TopListEntry
	err := h.apiRequest(ctx, h.config.APIEndpoints["top_ten_turnover"], &topTurnover)
	if err != nil {
		return nil, fmt.Errorf("failed to get top ten turnover: %w", err)
	}
	return topTurnover, nil
}

// Price and Trading Data Methods

// GetTodaysPrices retrieves today's price data, optionally filtered by business date
func (h *HTTPClient) GetTodaysPrices(ctx context.Context, businessDate string) ([]TodayPrice, error) {
	endpoint := h.config.APIEndpoints["todays_price"]
	if businessDate != "" {
		endpoint += "?businessDate=" + businessDate + "&size=500"
	}

	var todayPrices []TodayPrice
	err := h.apiRequest(ctx, endpoint, &todayPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's prices: %w", err)
	}
	return todayPrices, nil
}

// GetPriceVolumeHistory retrieves price volume history for a security
func (h *HTTPClient) GetPriceVolumeHistory(ctx context.Context, securityID int32, startDate, endDate string) ([]PriceHistory, error) {
	endpoint := fmt.Sprintf("%s%d?size=500&startDate=%s&endDate=%s",
		h.config.APIEndpoints["company_price_volume_history"], securityID, startDate, endDate)

	// The API returns a paginated response with content array
	var response struct {
		Content []PriceHistory `json:"content"`
	}

	err := h.apiRequest(ctx, endpoint, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get price volume history for security %d: %w", securityID, err)
	}
	return response.Content, nil
}

// GetMarketDepth retrieves market depth information for a security
func (h *HTTPClient) GetMarketDepth(ctx context.Context, securityID int32) (*MarketDepth, error) {
	endpoint := fmt.Sprintf("%s%d/", h.config.APIEndpoints["market_depth"], securityID)

	var marketDepth MarketDepth
	err := h.apiRequest(ctx, endpoint, &marketDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to get market depth for security %d: %w", securityID, err)
	}
	return &marketDepth, nil
}

// Security and Company Methods

// GetSecurityList retrieves the list of all securities
func (h *HTTPClient) GetSecurityList(ctx context.Context) ([]Security, error) {
	var securities []Security
	err := h.apiRequest(ctx, h.config.APIEndpoints["security_list"], &securities)
	if err != nil {
		return nil, fmt.Errorf("failed to get security list: %w", err)
	}
	return securities, nil
}

// GetCompanyList retrieves the list of all companies
func (h *HTTPClient) GetCompanyList(ctx context.Context) ([]Company, error) {
	var companies []Company
	err := h.apiRequest(ctx, h.config.APIEndpoints["company_list"], &companies)
	if err != nil {
		return nil, fmt.Errorf("failed to get company list: %w", err)
	}
	return companies, nil
}

// GetCompanyDetails retrieves detailed information about a specific company/security
func (h *HTTPClient) GetCompanyDetails(ctx context.Context, securityID int32) (*CompanyDetails, error) {
	endpoint := fmt.Sprintf("%s%d", h.config.APIEndpoints["company_details"], securityID)

	var rawDetails CompanyDetailsRaw
	err := h.apiRequest(ctx, endpoint, &rawDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to get company details for security %d: %w", securityID, err)
	}

	// Convert raw nested response to flat structured response
	details := &CompanyDetails{
		ID:               rawDetails.SecurityData.ID,
		Symbol:           rawDetails.SecurityData.Symbol,
		SecurityName:     rawDetails.SecurityData.SecurityName,
		SectorName:       rawDetails.SecurityData.Sector,
		Email:            rawDetails.SecurityData.Email,
		ActiveStatus:     rawDetails.SecurityData.ActiveStatus,
		PermittedToTrade: rawDetails.SecurityData.PermittedToTrade,

		// Market data from securityMcsData
		OpenPrice:           rawDetails.SecurityMcsData.OpenPrice,
		HighPrice:           rawDetails.SecurityMcsData.HighPrice,
		LowPrice:            rawDetails.SecurityMcsData.LowPrice,
		ClosePrice:          rawDetails.SecurityMcsData.ClosePrice,
		LastTradedPrice:     rawDetails.SecurityMcsData.LastTradedPrice,
		PreviousClose:       rawDetails.SecurityMcsData.PreviousClose,
		TotalTradeQuantity:  rawDetails.SecurityMcsData.TotalTradeQuantity,
		TotalTrades:         rawDetails.SecurityMcsData.TotalTrades,
		FiftyTwoWeekHigh:    rawDetails.SecurityMcsData.FiftyTwoWeekHigh,
		FiftyTwoWeekLow:     rawDetails.SecurityMcsData.FiftyTwoWeekLow,
		BusinessDate:        rawDetails.SecurityMcsData.BusinessDate,
		LastUpdatedDateTime: rawDetails.SecurityMcsData.LastUpdatedDateTime,
	}

	return details, nil
}

// GetSectorScrips groups securities by their sector
// Note: This is a slower operation as it requires fetching detailed company info
func (h *HTTPClient) GetSectorScrips(ctx context.Context) (SectorScrips, error) {
	// Get security list
	securities, err := h.GetSecurityList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get security list: %w", err)
	}

	// Group securities by sector - get sector info from company details
	sectorScrips := make(SectorScrips)

	// Process first 50 securities to avoid overwhelming the API
	maxSecurities := len(securities)
	if maxSecurities > 50 {
		maxSecurities = 50
	}

	for i := 0; i < maxSecurities; i++ {
		security := securities[i]

		// Skip promoter shares (contain "P" suffix typically)
		if strings.Contains(security.Symbol, "P") && strings.HasSuffix(security.Symbol, "P") {
			if sectorScrips["Promoter Share"] == nil {
				sectorScrips["Promoter Share"] = make([]string, 0)
			}
			sectorScrips["Promoter Share"] = append(sectorScrips["Promoter Share"], security.Symbol)
			continue
		}

		// Get company details to find sector
		details, err := h.GetCompanyDetails(ctx, security.ID)
		if err != nil {
			// If we can't get details, put in "Others" category
			if sectorScrips["Others"] == nil {
				sectorScrips["Others"] = make([]string, 0)
			}
			sectorScrips["Others"] = append(sectorScrips["Others"], security.Symbol)
			continue
		}

		sectorName := details.SectorName
		if sectorName == "" {
			sectorName = "Others"
		}

		if sectorScrips[sectorName] == nil {
			sectorScrips[sectorName] = make([]string, 0)
		}
		sectorScrips[sectorName] = append(sectorScrips[sectorName], security.Symbol)

		// Add small delay to avoid overwhelming the API
		time.Sleep(50 * time.Millisecond)
	}

	return sectorScrips, nil
}

// Helper Methods

// FindSecurityBySymbol finds a security by its symbol
func (h *HTTPClient) FindSecurityBySymbol(ctx context.Context, symbol string) (*Security, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, NewInvalidClientRequestError("symbol cannot be empty")
	}

	securities, err := h.GetSecurityList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get security list: %w", err)
	}

	for _, security := range securities {
		if security.Symbol == symbol {
			return &security, nil
		}
	}

	return nil, NewNotFoundError("security with symbol " + symbol)
}

// FindCompanyBySymbol finds a company by its symbol
func (h *HTTPClient) FindCompanyBySymbol(ctx context.Context, symbol string) (*Company, error) {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	if symbol == "" {
		return nil, NewInvalidClientRequestError("symbol cannot be empty")
	}

	companies, err := h.GetCompanyList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get company list: %w", err)
	}

	for _, company := range companies {
		if company.Symbol == symbol {
			return &company, nil
		}
	}

	return nil, NewNotFoundError("company with symbol " + symbol)
}

// Floor Sheet Methods

// GetFloorSheet retrieves the complete floor sheet data
func (h *HTTPClient) GetFloorSheet(ctx context.Context) ([]FloorSheetEntry, error) {
	endpoint := fmt.Sprintf("%s?size=500&sort=contractId,desc", h.config.APIEndpoints["floor_sheet"])

	// Based on error messages, the API returns array directly
	var floorSheetArray []FloorSheetEntry
	err := h.apiRequest(ctx, endpoint, &floorSheetArray)
	if err != nil {
		return nil, fmt.Errorf("failed to get floor sheet: %w", err)
	}

	return floorSheetArray, nil
}

// GetFloorSheetOf retrieves floor sheet data for a specific security on a specific business date
func (h *HTTPClient) GetFloorSheetOf(ctx context.Context, securityID int32, businessDate string) ([]FloorSheetEntry, error) {
	endpoint := fmt.Sprintf("%s%d?businessDate=%s&size=500&sort=contractid,desc",
		h.config.APIEndpoints["company_floorsheet"], securityID, businessDate)

	// Get first page
	var firstPage FloorSheetResponse
	err := h.apiRequest(ctx, endpoint, &firstPage)
	if err != nil {
		return nil, fmt.Errorf("failed to get floor sheet for security %d: %w", securityID, err)
	}

	// Check if there's any data
	if len(firstPage.FloorSheets.Content) == 0 {
		return []FloorSheetEntry{}, nil
	}

	allEntries := firstPage.FloorSheets.Content
	totalPages := firstPage.FloorSheets.TotalPages

	// Get remaining pages
	for page := int32(1); page < totalPages; page++ {
		pageEndpoint := fmt.Sprintf("%s&page=%d", endpoint, page)

		var pageResponse FloorSheetResponse
		err := h.apiRequest(ctx, pageEndpoint, &pageResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to get floor sheet page %d for security %d: %w", page, securityID, err)
		}

		allEntries = append(allEntries, pageResponse.FloorSheets.Content...)
	}

	return allEntries, nil
}

// Convenience method to get floor sheet by symbol instead of ID
func (h *HTTPClient) GetFloorSheetBySymbol(ctx context.Context, symbol string, businessDate string) ([]FloorSheetEntry, error) {
	security, err := h.FindSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to find security %s: %w", symbol, err)
	}

	return h.GetFloorSheetOf(ctx, security.ID, businessDate)
}

// Convenience method to get price history by symbol instead of ID
func (h *HTTPClient) GetPriceVolumeHistoryBySymbol(ctx context.Context, symbol string, startDate, endDate string) ([]PriceHistory, error) {
	security, err := h.FindSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to find security %s: %w", symbol, err)
	}

	return h.GetPriceVolumeHistory(ctx, security.ID, startDate, endDate)
}

// Convenience method to get market depth by symbol instead of ID
func (h *HTTPClient) GetMarketDepthBySymbol(ctx context.Context, symbol string) (*MarketDepth, error) {
	security, err := h.FindSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to find security %s: %w", symbol, err)
	}

	return h.GetMarketDepth(ctx, security.ID)
}

// Convenience method to get company details by symbol instead of ID
func (h *HTTPClient) GetCompanyDetailsBySymbol(ctx context.Context, symbol string) (*CompanyDetails, error) {
	security, err := h.FindSecurityBySymbol(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("failed to find security %s: %w", symbol, err)
	}

	return h.GetCompanyDetails(ctx, security.ID)
}
