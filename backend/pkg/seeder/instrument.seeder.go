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
	{Symbol: "AAPL", Name: "Apple Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/apple.com"},
	{Symbol: "MSFT", Name: "Microsoft Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/microsoft.com"},
	{Symbol: "GOOGL", Name: "Alphabet Inc. Class A", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/google.com"},
	{Symbol: "GOOG", Name: "Alphabet Inc. Class C", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/google.com"},
	{Symbol: "AMZN", Name: "Amazon.com Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/amazon.com"},
	{Symbol: "NVDA", Name: "NVIDIA Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/nvidia.com"},
	{Symbol: "META", Name: "Meta Platforms Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/meta.com"},
	{Symbol: "TSLA", Name: "Tesla Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/tesla.com"},
	{Symbol: "NFLX", Name: "Netflix Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/netflix.com"},
	{Symbol: "AMD", Name: "Advanced Micro Devices", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/amd.com"},
	{Symbol: "INTC", Name: "Intel Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/intel.com"},
	{Symbol: "CRM", Name: "Salesforce Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/salesforce.com"},
	{Symbol: "ORCL", Name: "Oracle Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/oracle.com"},
	{Symbol: "ADBE", Name: "Adobe Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/adobe.com"},
	{Symbol: "CSCO", Name: "Cisco Systems", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/cisco.com"},
	{Symbol: "AVGO", Name: "Broadcom Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/broadcom.com"},
	{Symbol: "QCOM", Name: "Qualcomm Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/qualcomm.com"},
	{Symbol: "TXN", Name: "Texas Instruments", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ti.com"},
	{Symbol: "IBM", Name: "IBM Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ibm.com"},
	{Symbol: "MU", Name: "Micron Technology", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/micron.com"},
	{Symbol: "UBER", Name: "Uber Technologies", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/uber.com"},
	{Symbol: "LYFT", Name: "Lyft Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lyft.com"},
	{Symbol: "SNAP", Name: "Snap Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/snap.com"},
	{Symbol: "PINS", Name: "Pinterest Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/pinterest.com"},
	{Symbol: "SPOT", Name: "Spotify Technology", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/spotify.com"},
	{Symbol: "SQ", Name: "Block Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/block.xyz"},
	{Symbol: "PYPL", Name: "PayPal Holdings", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/paypal.com"},
	{Symbol: "SHOP", Name: "Shopify Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/shopify.com"},
	{Symbol: "ROKU", Name: "Roku Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/roku.com"},
	{Symbol: "ZM", Name: "Zoom Video Communications", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/zoom.us"},
	{Symbol: "DOCU", Name: "DocuSign Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/docusign.com"},
	{Symbol: "TWLO", Name: "Twilio Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/twilio.com"},
	{Symbol: "NOW", Name: "ServiceNow Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/servicenow.com"},
	{Symbol: "SNOW", Name: "Snowflake Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/snowflake.com"},
	{Symbol: "PLTR", Name: "Palantir Technologies", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/palantir.com"},
	{Symbol: "COIN", Name: "Coinbase Global", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/coinbase.com"},
	{Symbol: "HOOD", Name: "Robinhood Markets", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/robinhood.com"},

	// ============ FINANCIAL STOCKS ============
	{Symbol: "JPM", Name: "JPMorgan Chase & Co.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/jpmorganchase.com"},
	{Symbol: "BAC", Name: "Bank of America", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/bankofamerica.com"},
	{Symbol: "WFC", Name: "Wells Fargo", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/wellsfargo.com"},
	{Symbol: "GS", Name: "Goldman Sachs", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/goldmansachs.com"},
	{Symbol: "MS", Name: "Morgan Stanley", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/morganstanley.com"},
	{Symbol: "C", Name: "Citigroup Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/citigroup.com"},
	{Symbol: "V", Name: "Visa Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/visa.com"},
	{Symbol: "MA", Name: "Mastercard Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/mastercard.com"},
	{Symbol: "AXP", Name: "American Express", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/americanexpress.com"},
	{Symbol: "BRK.B", Name: "Berkshire Hathaway B", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/berkshirehathaway.com"},
	{Symbol: "BLK", Name: "BlackRock Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/blackrock.com"},
	{Symbol: "SCHW", Name: "Charles Schwab", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/schwab.com"},

	// ============ CONSUMER & RETAIL ============
	{Symbol: "WMT", Name: "Walmart Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/walmart.com"},
	{Symbol: "COST", Name: "Costco Wholesale", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/costco.com"},
	{Symbol: "TGT", Name: "Target Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/target.com"},
	{Symbol: "HD", Name: "Home Depot", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/homedepot.com"},
	{Symbol: "LOW", Name: "Lowe's Companies", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lowes.com"},
	{Symbol: "NKE", Name: "Nike Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/nike.com"},
	{Symbol: "SBUX", Name: "Starbucks Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/starbucks.com"},
	{Symbol: "MCD", Name: "McDonald's Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/mcdonalds.com"},
	{Symbol: "KO", Name: "Coca-Cola Company", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/coca-cola.com"},
	{Symbol: "PEP", Name: "PepsiCo Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/pepsico.com"},
	{Symbol: "DIS", Name: "Walt Disney Company", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/disney.com"},
	{Symbol: "CMCSA", Name: "Comcast Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/comcast.com"},

	// ============ HEALTHCARE & PHARMA ============
	{Symbol: "JNJ", Name: "Johnson & Johnson", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/jnj.com"},
	{Symbol: "UNH", Name: "UnitedHealth Group", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/unitedhealthgroup.com"},
	{Symbol: "PFE", Name: "Pfizer Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/pfizer.com"},
	{Symbol: "ABBV", Name: "AbbVie Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/abbvie.com"},
	{Symbol: "MRK", Name: "Merck & Co.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/merck.com"},
	{Symbol: "LLY", Name: "Eli Lilly", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lilly.com"},
	{Symbol: "TMO", Name: "Thermo Fisher Scientific", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/thermofisher.com"},
	{Symbol: "ABT", Name: "Abbott Laboratories", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/abbott.com"},
	{Symbol: "BMY", Name: "Bristol-Myers Squibb", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/bms.com"},
	{Symbol: "GILD", Name: "Gilead Sciences", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/gilead.com"},
	{Symbol: "MRNA", Name: "Moderna Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/modernatx.com"},
	{Symbol: "BIIB", Name: "Biogen Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/biogen.com"},

	// ============ ENERGY & INDUSTRIAL ============
	{Symbol: "XOM", Name: "Exxon Mobil", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/exxonmobil.com"},
	{Symbol: "CVX", Name: "Chevron Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/chevron.com"},
	{Symbol: "COP", Name: "ConocoPhillips", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/conocophillips.com"},
	{Symbol: "SLB", Name: "Schlumberger", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/slb.com"},
	{Symbol: "BA", Name: "Boeing Company", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/boeing.com"},
	{Symbol: "CAT", Name: "Caterpillar Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/caterpillar.com"},
	{Symbol: "DE", Name: "Deere & Company", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/deere.com"},
	{Symbol: "HON", Name: "Honeywell International", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/honeywell.com"},
	{Symbol: "UPS", Name: "United Parcel Service", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ups.com"},
	{Symbol: "FDX", Name: "FedEx Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/fedex.com"},
	{Symbol: "GE", Name: "General Electric", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ge.com"},
	{Symbol: "MMM", Name: "3M Company", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/3m.com"},
	{Symbol: "RTX", Name: "RTX Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/rtx.com"},
	{Symbol: "LMT", Name: "Lockheed Martin", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lockheedmartin.com"},

	// ============ EV & CLEAN ENERGY ============
	{Symbol: "RIVN", Name: "Rivian Automotive", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/rivian.com"},
	{Symbol: "LCID", Name: "Lucid Group", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lucidmotors.com"},
	{Symbol: "NIO", Name: "NIO Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/nio.com"},
	{Symbol: "XPEV", Name: "XPeng Inc.", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/xiaopeng.com"},
	{Symbol: "LI", Name: "Li Auto Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lixiang.com"},
	{Symbol: "ENPH", Name: "Enphase Energy", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/enphase.com"},
	{Symbol: "SEDG", Name: "SolarEdge Technologies", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/solaredge.com"},
	{Symbol: "FSLR", Name: "First Solar", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/firstsolar.com"},

	// ============ GAMING & ENTERTAINMENT ============
	{Symbol: "EA", Name: "Electronic Arts", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ea.com"},
	{Symbol: "TTWO", Name: "Take-Two Interactive", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/take2games.com"},
	{Symbol: "ATVI", Name: "Activision Blizzard", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/activisionblizzard.com"},
	{Symbol: "RBLX", Name: "Roblox Corporation", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/roblox.com"},
	{Symbol: "U", Name: "Unity Software", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/unity.com"},

	// ============ CHINESE TECH (ADRs) ============
	{Symbol: "BABA", Name: "Alibaba Group", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/alibaba.com"},
	{Symbol: "JD", Name: "JD.com Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/jd.com"},
	{Symbol: "PDD", Name: "PDD Holdings (Pinduoduo)", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/pinduoduo.com"},
	{Symbol: "BIDU", Name: "Baidu Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/baidu.com"},
	{Symbol: "NTES", Name: "NetEase Inc.", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/netease.com"},

	// ============ AI & SEMICONDUCTORS ============
	{Symbol: "ARM", Name: "ARM Holdings", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/arm.com"},
	{Symbol: "MRVL", Name: "Marvell Technology", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/marvell.com"},
	{Symbol: "ON", Name: "ON Semiconductor", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/onsemi.com"},
	{Symbol: "LRCX", Name: "Lam Research", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/lamresearch.com"},
	{Symbol: "AMAT", Name: "Applied Materials", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/appliedmaterials.com"},
	{Symbol: "KLAC", Name: "KLA Corporation", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/kla.com"},
	{Symbol: "ASML", Name: "ASML Holding", Type: model.InstrumentTypeSock, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/asml.com"},
	{Symbol: "TSM", Name: "Taiwan Semiconductor", Type: model.InstrumentTypeSock, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/tsmc.com"},

	// ============ POPULAR ETFs ============
	{Symbol: "SPY", Name: "SPDR S&P 500 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Tracks the S&P 500 Index"},
	{Symbol: "QQQ", Name: "Invesco QQQ Trust", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/invesco.com", Description: "Tracks the NASDAQ-100 Index"},
	{Symbol: "IWM", Name: "iShares Russell 2000 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "Small-cap US stocks"},
	{Symbol: "DIA", Name: "SPDR Dow Jones ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Tracks the Dow Jones Industrial Average"},
	{Symbol: "VTI", Name: "Vanguard Total Stock Market", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/vanguard.com", Description: "Total US stock market"},
	{Symbol: "VOO", Name: "Vanguard S&P 500 ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/vanguard.com", Description: "Tracks the S&P 500 Index"},
	{Symbol: "VGT", Name: "Vanguard Information Technology", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/vanguard.com", Description: "US Tech sector"},
	{Symbol: "ARKK", Name: "ARK Innovation ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ark-invest.com", Description: "Disruptive innovation"},
	{Symbol: "ARKW", Name: "ARK Next Gen Internet ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ark-invest.com", Description: "Next gen internet"},
	{Symbol: "XLK", Name: "Technology Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Tech sector"},
	{Symbol: "XLF", Name: "Financial Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Financial sector"},
	{Symbol: "XLE", Name: "Energy Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Energy sector"},
	{Symbol: "XLV", Name: "Health Care Select Sector SPDR", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ssga.com", Description: "Healthcare sector"},
	{Symbol: "SOXX", Name: "iShares Semiconductor ETF", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "Semiconductor stocks"},
	{Symbol: "SMH", Name: "VanEck Semiconductor ETF", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/vaneck.com", Description: "Semiconductor stocks"},
	{Symbol: "GLD", Name: "SPDR Gold Shares", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/spdrgoldshares.com", Description: "Tracks gold prices"},
	{Symbol: "SLV", Name: "iShares Silver Trust", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "Tracks silver prices"},
	{Symbol: "TLT", Name: "iShares 20+ Year Treasury Bond", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "Long-term US treasuries"},
	{Symbol: "HYG", Name: "iShares High Yield Corporate Bond", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "High yield bonds"},
	{Symbol: "VWO", Name: "Vanguard FTSE Emerging Markets", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/vanguard.com", Description: "Emerging markets"},
	{Symbol: "EEM", Name: "iShares MSCI Emerging Markets", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "Emerging markets"},
	{Symbol: "EFA", Name: "iShares MSCI EAFE ETF", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/ishares.com", Description: "International developed markets"},

	// ============ LEVERAGED & INVERSE ETFs ============
	{Symbol: "TQQQ", Name: "ProShares UltraPro QQQ", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/proshares.com", Description: "3x NASDAQ-100"},
	{Symbol: "SQQQ", Name: "ProShares UltraPro Short QQQ", Type: model.InstrumentTypeETF, Exchange: "NASDAQ", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/proshares.com", Description: "-3x NASDAQ-100"},
	{Symbol: "SPXL", Name: "Direxion Daily S&P 500 Bull 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/direxion.com", Description: "3x S&P 500"},
	{Symbol: "SPXS", Name: "Direxion Daily S&P 500 Bear 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/direxion.com", Description: "-3x S&P 500"},
	{Symbol: "SOXL", Name: "Direxion Daily Semiconductor Bull 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/direxion.com", Description: "3x Semiconductors"},
	{Symbol: "SOXS", Name: "Direxion Daily Semiconductor Bear 3X", Type: model.InstrumentTypeETF, Exchange: "NYSE", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/direxion.com", Description: "-3x Semiconductors"},
	{Symbol: "UVXY", Name: "ProShares Ultra VIX Short-Term", Type: model.InstrumentTypeETF, Exchange: "BATS", Currency: "USD", Status: model.InstrumentStatusActive, LogoURL: "https://logo.clearbit.com/proshares.com", Description: "1.5x VIX futures"},

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
