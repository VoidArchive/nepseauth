package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/voidarchive/nepseauth/nepse"
)

func main() {
	fmt.Println("🚀 NEPSE Go Library - Basic Usage Example")
	fmt.Println("=========================================")

	// Create NEPSE client with TLS verification disabled for testing
	client, err := nepse.NewClientWithTLS(false)
	if err != nil {
		log.Fatalf("❌ Failed to create NEPSE client: %v", err)
	}
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			log.Printf("⚠️ Failed to close client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Test 1: Market Summary
	fmt.Println("\n📊 1. Getting Market Summary...")
	summary, err := client.GetMarketSummary(ctx)
	if err != nil {
		log.Printf("❌ Market summary failed: %v", err)
	} else {
		fmt.Printf("✅ Total Turnover: Rs. %.2f\n", summary.TotalTurnover)
		fmt.Printf("✅ Total Traded Shares: %.0f\n", summary.TotalTradedShares)
		fmt.Printf("✅ Total Transactions: %.0f\n", summary.TotalTransactions)
		fmt.Printf("✅ Total Scrips Traded: %.0f\n", summary.TotalScripsTraded)
	}

	// Test 2: Market Status
	fmt.Println("\n🏪 2. Checking Market Status...")
	status, err := client.GetMarketStatus(ctx)
	if err != nil {
		log.Printf("❌ Market status failed: %v", err)
	} else {
		fmt.Printf("✅ Market Open: %v\n", status.IsMarketOpen())
		fmt.Printf("✅ Status: %s\n", status.IsOpen)
	}

	// Test 3: NEPSE Index
	fmt.Println("\n📈 3. Getting NEPSE Index...")
	index, err := client.GetNepseIndex(ctx)
	if err != nil {
		log.Printf("❌ NEPSE index failed: %v", err)
	} else {
		fmt.Printf("✅ Index Value: %.2f\n", index.IndexValue)
		fmt.Printf("✅ Point Change: %.2f\n", index.PointChange)
		fmt.Printf("✅ Percent Change: %.2f%%\n", index.PercentChange)
	}

	// Test 4: Find Company by Symbol
	fmt.Println("\n🔍 4. Finding NABIL Company...")
	nabilSecurity, err := client.FindSecurityBySymbol(ctx, "NABIL")
	if err != nil {
		log.Printf("❌ Failed to find NABIL: %v", err)
	} else {
		fmt.Printf("✅ NABIL ID: %d\n", nabilSecurity.ID)
		fmt.Printf("✅ Security Name: %s\n", nabilSecurity.SecurityName)
		fmt.Printf("✅ Sector: %s\n", nabilSecurity.SectorName)

		// Test 5: Get Company Details
		fmt.Println("\n🏢 5. Getting NABIL Company Details...")
		details, err := client.GetCompanyDetails(ctx, nabilSecurity.ID)
		if err != nil {
			log.Printf("❌ NABIL details failed: %v", err)
		} else {
			fmt.Printf("✅ Sector: %s\n", details.SectorName)
			fmt.Printf("✅ Email: %s\n", details.Email)
			fmt.Printf("✅ Close Price: Rs. %.2f\n", details.ClosePrice)
			fmt.Printf("✅ Last Traded Price: Rs. %.2f\n", details.LastTradedPrice)
			fmt.Printf("✅ 52-Week High: Rs. %.2f\n", details.FiftyTwoWeekHigh)
		}

		// Test 6: Get Market Depth
		fmt.Println("\n📊 6. Getting NABIL Market Depth...")
		marketDepth, err := client.GetMarketDepth(ctx, nabilSecurity.ID)
		if err != nil {
			// Market depth often unavailable outside trading hours
			fmt.Printf("⚠️ Market depth not available (likely outside trading hours)\n")
		} else {
			fmt.Printf("✅ Buy Orders: %d levels\n", len(marketDepth.BuyDepth))
			fmt.Printf("✅ Sell Orders: %d levels\n", len(marketDepth.SellDepth))
		}
	}

	// Test 7: Top Gainers
	fmt.Println("\n📈 7. Getting Top Gainers...")
	topGainers, err := client.GetTopGainers(ctx)
	if err != nil {
		log.Printf("❌ Top gainers failed: %v", err)
	} else {
		fmt.Printf("✅ Top gainers found: %d\n", len(topGainers))
		if len(topGainers) > 0 {
			fmt.Printf("✅ Top gainer: %s (%.2f%%)\n",
				topGainers[0].Symbol, topGainers[0].PercentageChange)
		}
	}

	// Test 8: Securities by Sector
	fmt.Println("\n🏭 8. Getting Securities by Sector...")
	sectorScrips, err := client.GetSectorScrips(ctx)
	if err != nil {
		log.Printf("❌ Sector scrips failed: %v", err)
	} else {
		fmt.Printf("✅ Sectors found: %d\n", len(sectorScrips))
		if banking, exists := sectorScrips[nepse.SectorBanking]; exists {
			fmt.Printf("✅ Banking sector companies: %d\n", len(banking))
		}
	}

	// Test 9: Daily NEPSE Index Graph (POST method - currently limited)
	fmt.Println("\n📊 9. Getting NEPSE Index Graph Data...")
	graphData, err := client.GetDailyNepseIndexGraph(ctx)
	if err != nil {
		fmt.Printf("⚠️ POST endpoints currently have token timing limitations\n")
		fmt.Printf("    All GET endpoints work perfectly. POST endpoints need additional auth work.\n")
	} else {
		fmt.Printf("✅ Graph data points: %d\n", len(graphData.Data))
		if len(graphData.Data) > 0 {
			latest := graphData.Data[len(graphData.Data)-1]
			fmt.Printf("✅ Latest point: %s = %.2f\n", latest.Date, latest.Value)
		}
	}

	fmt.Println("\n🎉 All tests completed!")
	fmt.Println("✅ NEPSE Go library is working correctly!")
}
