package seeder

import (
	"context"
	"log"
	"time"

	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/instrument/repository"
)

// PopularInstruments - comprehensive list of stocks, ETFs, and crypto
var PopularInstruments = []model.Instrument{
	// ============ TOP US TECH STOCKS ============
	{Symbol: "AAPL", Name: "Apple Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/apple.com"},
	{Symbol: "MSFT", Name: "Microsoft Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/microsoft.com"},
	{Symbol: "GOOGL", Name: "Alphabet Inc. Class A", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/google.com"},
	{Symbol: "GOOG", Name: "Alphabet Inc. Class C", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/google.com"},
	{Symbol: "AMZN", Name: "Amazon.com Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/amazon.com"},
	{Symbol: "NVDA", Name: "NVIDIA Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/nvidia.com"},
	{Symbol: "META", Name: "Meta Platforms Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/meta.com"},
	{Symbol: "TSLA", Name: "Tesla Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/tesla.com"},
	{Symbol: "NFLX", Name: "Netflix Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/netflix.com"},
	{Symbol: "AMD", Name: "Advanced Micro Devices", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/amd.com"},
	{Symbol: "INTC", Name: "Intel Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/intel.com"},
	{Symbol: "CRM", Name: "Salesforce Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/salesforce.com"},
	{Symbol: "ORCL", Name: "Oracle Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/oracle.com"},
	{Symbol: "ADBE", Name: "Adobe Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/adobe.com"},
	{Symbol: "CSCO", Name: "Cisco Systems", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/cisco.com"},
	{Symbol: "AVGO", Name: "Broadcom Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/broadcom.com"},
	{Symbol: "QCOM", Name: "Qualcomm Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/qualcomm.com"},
	{Symbol: "TXN", Name: "Texas Instruments", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ti.com"},
	{Symbol: "IBM", Name: "IBM Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ibm.com"},
	{Symbol: "MU", Name: "Micron Technology", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/micron.com"},
	{Symbol: "UBER", Name: "Uber Technologies", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/uber.com"},
	{Symbol: "LYFT", Name: "Lyft Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lyft.com"},
	{Symbol: "SNAP", Name: "Snap Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/snap.com"},
	{Symbol: "PINS", Name: "Pinterest Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/pinterest.com"},
	{Symbol: "SPOT", Name: "Spotify Technology", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/spotify.com"},
	{Symbol: "SQ", Name: "Block Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/block.xyz"},
	{Symbol: "PYPL", Name: "PayPal Holdings", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/paypal.com"},
	{Symbol: "SHOP", Name: "Shopify Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/shopify.com"},
	{Symbol: "ROKU", Name: "Roku Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/roku.com"},
	{Symbol: "ZM", Name: "Zoom Video Communications", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/zoom.us"},
	{Symbol: "DOCU", Name: "DocuSign Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/docusign.com"},
	{Symbol: "TWLO", Name: "Twilio Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/twilio.com"},
	{Symbol: "NOW", Name: "ServiceNow Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/servicenow.com"},
	{Symbol: "SNOW", Name: "Snowflake Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/snowflake.com"},
	{Symbol: "PLTR", Name: "Palantir Technologies", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/palantir.com"},
	{Symbol: "COIN", Name: "Coinbase Global", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/coinbase.com"},
	{Symbol: "HOOD", Name: "Robinhood Markets", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/robinhood.com"},

	// ============ FINANCIAL STOCKS ============
	{Symbol: "JPM", Name: "JPMorgan Chase & Co.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/jpmorganchase.com"},
	{Symbol: "BAC", Name: "Bank of America", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/bankofamerica.com"},
	{Symbol: "WFC", Name: "Wells Fargo", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/wellsfargo.com"},
	{Symbol: "GS", Name: "Goldman Sachs", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/goldmansachs.com"},
	{Symbol: "MS", Name: "Morgan Stanley", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/morganstanley.com"},
	{Symbol: "C", Name: "Citigroup Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/citigroup.com"},
	{Symbol: "V", Name: "Visa Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/visa.com"},
	{Symbol: "MA", Name: "Mastercard Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/mastercard.com"},
	{Symbol: "AXP", Name: "American Express", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/americanexpress.com"},
	{Symbol: "BRK.B", Name: "Berkshire Hathaway B", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/berkshirehathaway.com"},
	{Symbol: "BLK", Name: "BlackRock Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/blackrock.com"},
	{Symbol: "SCHW", Name: "Charles Schwab", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/schwab.com"},

	// ============ CONSUMER & RETAIL ============
	{Symbol: "WMT", Name: "Walmart Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/walmart.com"},
	{Symbol: "COST", Name: "Costco Wholesale", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/costco.com"},
	{Symbol: "TGT", Name: "Target Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/target.com"},
	{Symbol: "HD", Name: "Home Depot", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/homedepot.com"},
	{Symbol: "LOW", Name: "Lowe's Companies", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lowes.com"},
	{Symbol: "NKE", Name: "Nike Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/nike.com"},
	{Symbol: "SBUX", Name: "Starbucks Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/starbucks.com"},
	{Symbol: "MCD", Name: "McDonald's Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/mcdonalds.com"},
	{Symbol: "KO", Name: "Coca-Cola Company", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/coca-cola.com"},
	{Symbol: "PEP", Name: "PepsiCo Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/pepsico.com"},
	{Symbol: "DIS", Name: "Walt Disney Company", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/disney.com"},
	{Symbol: "CMCSA", Name: "Comcast Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/comcast.com"},

	// ============ HEALTHCARE & PHARMA ============
	{Symbol: "JNJ", Name: "Johnson & Johnson", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/jnj.com"},
	{Symbol: "UNH", Name: "UnitedHealth Group", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/unitedhealthgroup.com"},
	{Symbol: "PFE", Name: "Pfizer Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/pfizer.com"},
	{Symbol: "ABBV", Name: "AbbVie Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/abbvie.com"},
	{Symbol: "MRK", Name: "Merck & Co.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/merck.com"},
	{Symbol: "LLY", Name: "Eli Lilly", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lilly.com"},
	{Symbol: "TMO", Name: "Thermo Fisher Scientific", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/thermofisher.com"},
	{Symbol: "ABT", Name: "Abbott Laboratories", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/abbott.com"},
	{Symbol: "BMY", Name: "Bristol-Myers Squibb", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/bms.com"},
	{Symbol: "GILD", Name: "Gilead Sciences", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/gilead.com"},
	{Symbol: "MRNA", Name: "Moderna Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/modernatx.com"},
	{Symbol: "BIIB", Name: "Biogen Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/biogen.com"},

	// ============ ENERGY & INDUSTRIAL ============
	{Symbol: "XOM", Name: "Exxon Mobil", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/exxonmobil.com"},
	{Symbol: "CVX", Name: "Chevron Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/chevron.com"},
	{Symbol: "COP", Name: "ConocoPhillips", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/conocophillips.com"},
	{Symbol: "SLB", Name: "Schlumberger", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/slb.com"},
	{Symbol: "BA", Name: "Boeing Company", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/boeing.com"},
	{Symbol: "CAT", Name: "Caterpillar Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/caterpillar.com"},
	{Symbol: "DE", Name: "Deere & Company", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/deere.com"},
	{Symbol: "HON", Name: "Honeywell International", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/honeywell.com"},
	{Symbol: "UPS", Name: "United Parcel Service", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ups.com"},
	{Symbol: "FDX", Name: "FedEx Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/fedex.com"},
	{Symbol: "GE", Name: "General Electric", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ge.com"},
	{Symbol: "MMM", Name: "3M Company", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/3m.com"},
	{Symbol: "RTX", Name: "RTX Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/rtx.com"},
	{Symbol: "LMT", Name: "Lockheed Martin", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lockheedmartin.com"},

	// ============ EV & CLEAN ENERGY ============
	{Symbol: "RIVN", Name: "Rivian Automotive", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/rivian.com"},
	{Symbol: "LCID", Name: "Lucid Group", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lucidmotors.com"},
	{Symbol: "NIO", Name: "NIO Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/nio.com"},
	{Symbol: "XPEV", Name: "XPeng Inc.", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/xiaopeng.com"},
	{Symbol: "LI", Name: "Li Auto Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lixiang.com"},
	{Symbol: "ENPH", Name: "Enphase Energy", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/enphase.com"},
	{Symbol: "SEDG", Name: "SolarEdge Technologies", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/solaredge.com"},
	{Symbol: "FSLR", Name: "First Solar", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/firstsolar.com"},

	// ============ GAMING & ENTERTAINMENT ============
	{Symbol: "EA", Name: "Electronic Arts", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ea.com"},
	{Symbol: "TTWO", Name: "Take-Two Interactive", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/take2games.com"},
	{Symbol: "ATVI", Name: "Activision Blizzard", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/activisionblizzard.com"},
	{Symbol: "RBLX", Name: "Roblox Corporation", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/roblox.com"},
	{Symbol: "U", Name: "Unity Software", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/unity.com"},

	// ============ CHINESE TECH (ADRs) ============
	{Symbol: "BABA", Name: "Alibaba Group", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/alibaba.com"},
	{Symbol: "JD", Name: "JD.com Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/jd.com"},
	{Symbol: "PDD", Name: "PDD Holdings (Pinduoduo)", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/pinduoduo.com"},
	{Symbol: "BIDU", Name: "Baidu Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/baidu.com"},
	{Symbol: "NTES", Name: "NetEase Inc.", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/netease.com"},

	// ============ AI & SEMICONDUCTORS ============
	{Symbol: "ARM", Name: "ARM Holdings", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/arm.com"},
	{Symbol: "MRVL", Name: "Marvell Technology", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/marvell.com"},
	{Symbol: "ON", Name: "ON Semiconductor", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/onsemi.com"},
	{Symbol: "LRCX", Name: "Lam Research", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/lamresearch.com"},
	{Symbol: "AMAT", Name: "Applied Materials", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/appliedmaterials.com"},
	{Symbol: "KLAC", Name: "KLA Corporation", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/kla.com"},
	{Symbol: "ASML", Name: "ASML Holding", Type: model.InstrumentTypeStock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/asml.com"},
	{Symbol: "TSM", Name: "Taiwan Semiconductor", Type: model.InstrumentTypeStock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/tsmc.com"},

	// ============ POPULAR ETFs ============
	{Symbol: "SPY", Name: "SPDR S&P 500 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Tracks the S&P 500 Index"},
	{Symbol: "QQQ", Name: "Invesco QQQ Trust", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/invesco.com", Description: "Tracks the NASDAQ-100 Index"},
	{Symbol: "IWM", Name: "iShares Russell 2000 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "Small-cap US stocks"},
	{Symbol: "DIA", Name: "SPDR Dow Jones ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Tracks the Dow Jones Industrial Average"},
	{Symbol: "VTI", Name: "Vanguard Total Stock Market", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/vanguard.com", Description: "Total US stock market"},
	{Symbol: "VOO", Name: "Vanguard S&P 500 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/vanguard.com", Description: "Tracks the S&P 500 Index"},
	{Symbol: "VGT", Name: "Vanguard Information Technology", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/vanguard.com", Description: "US Tech sector"},
	{Symbol: "ARKK", Name: "ARK Innovation ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ark-invest.com", Description: "Disruptive innovation"},
	{Symbol: "ARKW", Name: "ARK Next Gen Internet ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ark-invest.com", Description: "Next gen internet"},
	{Symbol: "XLK", Name: "Technology Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Tech sector"},
	{Symbol: "XLF", Name: "Financial Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Financial sector"},
	{Symbol: "XLE", Name: "Energy Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Energy sector"},
	{Symbol: "XLV", Name: "Health Care Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ssga.com", Description: "Healthcare sector"},
	{Symbol: "SOXX", Name: "iShares Semiconductor ETF", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "Semiconductor stocks"},
	{Symbol: "SMH", Name: "VanEck Semiconductor ETF", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/vaneck.com", Description: "Semiconductor stocks"},
	{Symbol: "GLD", Name: "SPDR Gold Shares", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/spdrgoldshares.com", Description: "Tracks gold prices"},
	{Symbol: "SLV", Name: "iShares Silver Trust", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "Tracks silver prices"},
	{Symbol: "TLT", Name: "iShares 20+ Year Treasury Bond", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "Long-term US treasuries"},
	{Symbol: "HYG", Name: "iShares High Yield Corporate Bond", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "High yield bonds"},
	{Symbol: "VWO", Name: "Vanguard FTSE Emerging Markets", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/vanguard.com", Description: "Emerging markets"},
	{Symbol: "EEM", Name: "iShares MSCI Emerging Markets", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "Emerging markets"},
	{Symbol: "EFA", Name: "iShares MSCI EAFE ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/ishares.com", Description: "International developed markets"},

	// ============ LEVERAGED & INVERSE ETFs ============
	{Symbol: "TQQQ", Name: "ProShares UltraPro QQQ", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/proshares.com", Description: "3x NASDAQ-100"},
	{Symbol: "SQQQ", Name: "ProShares UltraPro Short QQQ", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/proshares.com", Description: "-3x NASDAQ-100"},
	{Symbol: "SPXL", Name: "Direxion Daily S&P 500 Bull 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/direxion.com", Description: "3x S&P 500"},
	{Symbol: "SPXS", Name: "Direxion Daily S&P 500 Bear 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/direxion.com", Description: "-3x S&P 500"},
	{Symbol: "SOXL", Name: "Direxion Daily Semiconductor Bull 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/direxion.com", Description: "3x Semiconductors"},
	{Symbol: "SOXS", Name: "Direxion Daily Semiconductor Bear 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/direxion.com", Description: "-3x Semiconductors"},
	{Symbol: "UVXY", Name: "ProShares Ultra VIX Short-Term", Type: model.InstrumentTypeETF, Exchange: "BATS", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.logo.dev/proshares.com", Description: "1.5x VIX futures"},

	// ============ CRYPTOCURRENCIES ============
	{Symbol: "BTC/USD", Name: "Bitcoin", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/bitcoin-btc-logo.png"},
	{Symbol: "ETH/USD", Name: "Ethereum", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/ethereum-eth-logo.png"},
	{Symbol: "SOL/USD", Name: "Solana", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/solana-sol-logo.png"},
	{Symbol: "XRP/USD", Name: "Ripple", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/xrp-xrp-logo.png"},
	{Symbol: "DOGE/USD", Name: "Dogecoin", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/dogecoin-doge-logo.png"},
	{Symbol: "ADA/USD", Name: "Cardano", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/cardano-ada-logo.png"},
	{Symbol: "AVAX/USD", Name: "Avalanche", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/avalanche-avax-logo.png"},
	{Symbol: "MATIC/USD", Name: "Polygon", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/polygon-matic-logo.png"},
	{Symbol: "DOT/USD", Name: "Polkadot", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/polkadot-new-dot-logo.png"},
	{Symbol: "LINK/USD", Name: "Chainlink", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/chainlink-link-logo.png"},
	{Symbol: "UNI/USD", Name: "Uniswap", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/uniswap-uni-logo.png"},
	{Symbol: "ATOM/USD", Name: "Cosmos", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/cosmos-atom-logo.png"},
	{Symbol: "LTC/USD", Name: "Litecoin", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/litecoin-ltc-logo.png"},
	{Symbol: "BCH/USD", Name: "Bitcoin Cash", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/bitcoin-cash-bch-logo.png"},
	{Symbol: "SHIB/USD", Name: "Shiba Inu", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/shiba-inu-shib-logo.png"},
	{Symbol: "PEPE/USD", Name: "Pepe", Type: model.InstrumentTypeCrypto, Exchange: "Crypto", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://cryptologos.cc/logos/pepe-pepe-logo.png"},

	// ============ COMMODITIES ============
	{Symbol: "XAU/USD", Name: "Gold", Type: model.InstrumentTypeCommodity, Exchange: "Forex", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/gold-bars.png", Description: "Spot Gold"},
	{Symbol: "XAG/USD", Name: "Silver", Type: model.InstrumentTypeCommodity, Exchange: "Forex", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/silver-bars.png", Description: "Spot Silver"},
	{Symbol: "WTI/USD", Name: "Crude Oil WTI", Type: model.InstrumentTypeCommodity, Exchange: "NYMEX", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/oil-industry.png", Description: "WTI Crude Oil"},

	// ============ FOREX PAIRS ============
	{Symbol: "EUR/USD", Name: "Euro / US Dollar", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/euro.png", Description: "Euro vs US Dollar"},
	{Symbol: "GBP/USD", Name: "British Pound / US Dollar", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/uk.png", Description: "British Pound vs US Dollar"},
	{Symbol: "USD/JPY", Name: "US Dollar / Japanese Yen", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "JPY", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/japan.png", Description: "US Dollar vs Japanese Yen"},
	{Symbol: "AUD/USD", Name: "Australian Dollar / US Dollar", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/australia.png", Description: "Australian Dollar vs US Dollar"},
	{Symbol: "USD/CAD", Name: "US Dollar / Canadian Dollar", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "CAD", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/canada.png", Description: "US Dollar vs Canadian Dollar"},
	{Symbol: "USD/CHF", Name: "US Dollar / Swiss Franc", Type: model.InstrumentTypeForex, Exchange: "Forex", Currency: "CHF", Status: model.InstrumentStatusActive, LogoURL: "https://img.icons8.com/color/96/switzerland.png", Description: "US Dollar vs Swiss Franc"},
}

// SeedInstruments seeds popular instruments if they don't exist
func SeedInstruments(ctx context.Context) error {
	repo := repository.NewInstrumentRepository()

	created := 0
	for _, instrument := range PopularInstruments {
		exists, err := repo.SymbolExists(ctx, instrument.Symbol)
		if err != nil {
			log.Printf("⚠️ Error checking instrument %s: %v", instrument.Symbol, err)
			continue
		}

		if !exists {
			instrument.CreatedAt = time.Now()
			instrument.UpdatedAt = time.Now()
			if err := repo.CreateInstrument(ctx, &instrument); err != nil {
				log.Printf("⚠️ Failed to create instrument %s: %v", instrument.Symbol, err)
				continue
			}
			created++
		}
	}

	log.Printf("🎉 Seeded %d new instruments (total in list: %d)", created, len(PopularInstruments))
	return nil
}
