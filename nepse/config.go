package nepse

// Config holds static configuration data for the NEPSE API
type Config struct {
	BaseURL      string
	APIEndpoints map[string]string
	Headers      map[string]string
}

// DefaultConfig returns the default NEPSE API configuration
func DefaultConfig() *Config {
	return &Config{
		BaseURL: "https://www.nepalstock.com",
		APIEndpoints: map[string]string{
			"price_volume":                       "/api/nots/securityDailyTradeStat/58",
			"market_summary":                     "/api/nots/market-summary/",
			"supply_demand":                      "/api/nots/nepse-data/supplydemand",
			"top_gainers":                        "/api/nots/top-ten/top-gainer",
			"top_losers":                         "/api/nots/top-ten/top-loser",
			"top_ten_trade":                      "/api/nots/top-ten/trade",
			"top_ten_transaction":                "/api/nots/top-ten/transaction",
			"top_ten_turnover":                   "/api/nots/top-ten/turnover",
			"market_open":                        "/api/nots/nepse-data/market-open",
			"nepse_index":                        "/api/nots/nepse-index",
			"company_list":                       "/api/nots/company/list",
			"security_list":                      "/api/nots/security?nonDelisted=true",
			"nepse_index_daily_graph":            "/api/nots/graph/index/58",
			"sensitive_index_daily_graph":        "/api/nots/graph/index/57",
			"float_index_daily_graph":            "/api/nots/graph/index/62",
			"sensitive_float_index_daily_graph":  "/api/nots/graph/index/63",
			"banking_sub_index_graph":            "/api/nots/graph/index/51",
			"development_bank_sub_index_graph":   "/api/nots/graph/index/55",
			"finance_sub_index_graph":            "/api/nots/graph/index/60",
			"hotel_tourism_sub_index_graph":      "/api/nots/graph/index/52",
			"hydro_sub_index_graph":              "/api/nots/graph/index/54",
			"investment_sub_index_graph":         "/api/nots/graph/index/67",
			"life_insurance_sub_index_graph":     "/api/nots/graph/index/65",
			"manufacturing_sub_index_graph":      "/api/nots/graph/index/56",
			"microfinance_sub_index_graph":       "/api/nots/graph/index/64",
			"mutual_fund_sub_index_graph":        "/api/nots/graph/index/66",
			"non_life_insurance_sub_index_graph": "/api/nots/graph/index/59",
			"others_sub_index_graph":             "/api/nots/graph/index/53",
			"trading_sub_index_graph":            "/api/nots/graph/index/61",
			"company_daily_graph":                "/api/nots/market/graphdata/daily/",
			"company_details":                    "/api/nots/security/",
			"company_price_volume_history":       "/api/nots/market/history/security/",
			"company_floorsheet":                 "/api/nots/security/floorsheet/",
			"floor_sheet":                        "/api/nots/nepse-data/floorsheet",
			"todays_price":                       "/api/nots/nepse-data/today-price",
			"live_market":                        "/api/nots/lives-market",
			"market_depth":                       "/api/nots/nepse-data/marketdepth/",
		},
        Headers: map[string]string{
            "User-Agent":      "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
            "Accept":          "application/json, text/plain, */*",
            "Accept-Language": "en-US,en;q=0.5",
            // Note: Deliberately omit Accept-Encoding to allow net/http to auto-handle gzip.
            // Also omit Connection to enable keep-alives.
            "Pragma":          "no-cache",
            "Cache-Control":   "no-cache",
            "TE":              "Trailers",
        },
	}
}
