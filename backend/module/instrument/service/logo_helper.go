package service

import (
	"fmt"
	"strings"
)

// IMPORTANT: Clearbit needs company DOMAIN not ticker symbol
// e.g. https://logo.clearbit.com/apple.com for AAPL

// Major company domains for logo lookup
var companyDomains = map[string]string{
	// Tech Giants
	"AAPL":  "apple.com",
	"MSFT":  "microsoft.com",
	"GOOGL": "google.com",
	"GOOG":  "google.com",
	"AMZN":  "amazon.com",
	"META":  "meta.com",
	"NVDA":  "nvidia.com",
	"TSLA":  "tesla.com",
	"AMD":   "amd.com",
	"INTC":  "intel.com",
	"CRM":   "salesforce.com",
	"ORCL":  "oracle.com",
	"IBM":   "ibm.com",
	"ADBE":  "adobe.com",
	"NFLX":  "netflix.com",
	"PYPL":  "paypal.com",
	"SHOP":  "shopify.com",
	"SQ":    "squareup.com",
	"SPOT":  "spotify.com",
	"UBER":  "uber.com",
	"LYFT":  "lyft.com",
	"ABNB":  "airbnb.com",
	"SNAP":  "snap.com",
	"PINS":  "pinterest.com",
	"TWTR":  "twitter.com",
	"X":     "twitter.com",
	"PLTR":  "palantir.com",
	"NOW":   "servicenow.com",
	"SNOW":  "snowflake.com",
	"DDOG":  "datadoghq.com",
	"ZS":    "zscaler.com",
	"NET":   "cloudflare.com",
	"CRWD":  "crowdstrike.com",
	"OKTA":  "okta.com",
	"MDB":   "mongodb.com",
	"COIN":  "coinbase.com",
	"HOOD":  "robinhood.com",
	"RBLX":  "roblox.com",
	"U":     "unity.com",
	"TEAM":  "atlassian.com",
	"ZM":    "zoom.us",
	"DOCU":  "docusign.com",
	"DBX":   "dropbox.com",
	"TWLO":  "twilio.com",
	"ROKU":  "roku.com",

	// Chinese ADRs
	"BABA": "alibaba.com",
	"JD":   "jd.com",
	"PDD":  "pinduoduo.com",
	"BIDU": "baidu.com",
	"NIO":  "nio.com",
	"XPEV": "heyxpeng.com",
	"LI":   "lixiang.com",

	// Finance
	"JPM":  "jpmorganchase.com",
	"BAC":  "bankofamerica.com",
	"WFC":  "wellsfargo.com",
	"GS":   "goldmansachs.com",
	"MS":   "morganstanley.com",
	"C":    "citi.com",
	"V":    "visa.com",
	"MA":   "mastercard.com",
	"AXP":  "americanexpress.com",
	"BLK":  "blackrock.com",
	"SCHW": "schwab.com",

	// Healthcare
	"JNJ":  "jnj.com",
	"UNH":  "unitedhealthgroup.com",
	"PFE":  "pfizer.com",
	"MRK":  "merck.com",
	"ABBV": "abbvie.com",
	"LLY":  "lilly.com",
	"TMO":  "thermofisher.com",
	"ABT":  "abbott.com",
	"BMY":  "bms.com",
	"AMGN": "amgen.com",
	"GILD": "gilead.com",
	"MRNA": "modernatx.com",
	"REGN": "regeneron.com",
	"VRTX": "vrtx.com",

	// Retail
	"WMT":   "walmart.com",
	"HD":    "homedepot.com",
	"COST":  "costco.com",
	"TGT":   "target.com",
	"LOW":   "lowes.com",
	"NKE":   "nike.com",
	"SBUX":  "starbucks.com",
	"MCD":   "mcdonalds.com",
	"DIS":   "disney.com",
	"CMCSA": "comcast.com",

	// Energy
	"XOM": "exxonmobil.com",
	"CVX": "chevron.com",
	"COP": "conocophillips.com",
	"SLB": "slb.com",
	"OXY": "oxy.com",
	"EOG": "eogresources.com",
	"PXD": "pxd.com",
	"MPC": "marathonpetroleum.com",
	"VLO": "valero.com",
	"PSX": "phillips66.com",

	// Industrials
	"CAT": "caterpillar.com",
	"BA":  "boeing.com",
	"HON": "honeywell.com",
	"UPS": "ups.com",
	"UNP": "up.com",
	"DE":  "deere.com",
	"GE":  "ge.com",
	"LMT": "lockheedmartin.com",
	"RTX": "rtx.com",
	"NOC": "northropgrumman.com",
	"MMM": "3m.com",

	// Telecom
	"T":    "att.com",
	"VZ":   "verizon.com",
	"TMUS": "t-mobile.com",

	// Consumer
	"KO":  "coca-cola.com",
	"PEP": "pepsico.com",
	"PG":  "pg.com",
	"PM":  "pmi.com",
	"MO":  "altria.com",
	"EL":  "esteelauder.com",
	"CL":  "colgatepalmolive.com",

	// Semiconductors
	"AVGO": "broadcom.com",
	"QCOM": "qualcomm.com",
	"TXN":  "ti.com",
	"MU":   "micron.com",
	"LRCX": "lamresearch.com",
	"KLAC": "kla.com",
	"ADI":  "analog.com",
	"MRVL": "marvell.com",
	"ON":   "onsemi.com",
	"SWKS": "skyworksinc.com",
	"ARM":  "arm.com",
	"TSM":  "tsmc.com",
	"ASML": "asml.com",

	// Gaming
	"EA":   "ea.com",
	"TTWO": "take2games.com",
	"ATVI": "activision.com",

	// EVs & Auto
	"F":    "ford.com",
	"GM":   "gm.com",
	"TM":   "toyota.com",
	"RIVN": "rivian.com",
	"LCID": "lucidmotors.com",

	// Crypto Mining
	"MARA": "mara.com",
	"RIOT": "riotplatforms.com",
	"CLSK": "cleanspark.com",
	"IREN": "iren.io",
	"CIFR": "cifr.com",
	"HUT":  "hut8.com",
	"BTBT": "bit-digital.com",
	"CORZ": "corescientific.com",
	"WULF": "terawulf.com",
	"ARBK": "argoblock.com",
	"BITF": "bitfarms.com",
	"HIVE": "hiveblockchaini.com",

	// ETFs (most don't have logos)
	"SPY":  "ssga.com",
	"QQQ":  "invesco.com",
	"IWM":  "ishares.com",
	"VTI":  "vanguard.com",
	"VOO":  "vanguard.com",
	"ARKK": "ark-funds.com",
	"DIA":  "ssga.com",
	"XLF":  "ssga.com",
	"XLE":  "ssga.com",
	"XLK":  "ssga.com",
	"GLD":  "ssga.com",
	"SLV":  "ishares.com",
	"TLT":  "ishares.com",
	"HYG":  "ishares.com",
	"EEM":  "ishares.com",
	"VWO":  "vanguard.com",
	"IEMG": "ishares.com",
	"LQD":  "ishares.com",
	"TQQQ": "proshares.com",
	"SQQQ": "proshares.com",
	"SPXL": "direxion.com",
	"SPXS": "direxion.com",
	"SOXL": "direxion.com",
	"SOXS": "direxion.com",
}

// Crypto logo mapping using CoinGecko IDs
var cryptoLogos = map[string]string{
	"BTC":   "bitcoin",
	"ETH":   "ethereum",
	"SOL":   "solana",
	"XRP":   "xrp",
	"DOGE":  "dogecoin",
	"ADA":   "cardano",
	"AVAX":  "avalanche",
	"DOT":   "polkadot",
	"LINK":  "chainlink",
	"UNI":   "uniswap",
	"ATOM":  "cosmos",
	"LTC":   "litecoin",
	"SHIB":  "shiba-inu",
	"PEPE":  "pepe",
	"MATIC": "polygon",
	"NEAR":  "near-protocol",
	"APT":   "aptos",
	"ARB":   "arbitrum",
	"OP":    "optimism",
	"SUI":   "sui",
	"SEI":   "sei",
	"FIL":   "filecoin",
	"ICP":   "internet-computer",
	"AAVE":  "aave",
	"MKR":   "maker-dao",
	"CRV":   "curve-dao",
	"SUSHI": "sushi",
	"COMP":  "compound",
	"SNX":   "synthetix",
	"LDO":   "lido-dao",
	"RNDR":  "render-token",
	"SAND":  "the-sandbox",
	"MANA":  "decentraland",
	"AXS":   "axie-infinity",
	"GALA":  "gala",
	"ENS":   "ethereum-name-service",
	"IMX":   "immutable-x",
	"CRO":   "crypto-com-coin",
	"BNB":   "binancecoin",
	"TRX":   "tron",
	"XLM":   "stellar",
	"ALGO":  "algorand",
	"EOS":   "eos",
	"HBAR":  "hedera",
	"VET":   "vechain",
	"EGLD":  "elrond-erd-2",
	"FTM":   "fantom",
	"INJ":   "injective-protocol",
	"TIA":   "celestia",
	"JUP":   "jupiter-exchange-solana",
	"BONK":  "bonk",
	"WIF":   "dogwifhat",
	"PYTH":  "pyth-network",
	"JTO":   "jito-governance-token",
}

// GetLogoURLForSymbol returns the best logo URL for a given symbol and type
func GetLogoURLForSymbol(symbol string, instrumentType string) string {
	// Handle crypto
	if instrumentType == "Crypto" {
		baseSymbol := strings.Split(symbol, "/")[0]
		if geckoID, ok := cryptoLogos[baseSymbol]; ok {
			return fmt.Sprintf("https://assets.coingecko.com/coins/images/%s/small/%s.png", getCoinGeckoImagePath(geckoID), geckoID)
		}
		// Fallback to CryptoCompare
		return fmt.Sprintf("https://www.cryptocompare.com/media/37746251/%s.png", strings.ToLower(baseSymbol))
	}

	// Handle stocks/ETFs with known domains
	if domain, ok := companyDomains[symbol]; ok {
		return fmt.Sprintf("https://logo.clearbit.com/%s", domain)
	}

	// Fallback: Try Google favicon (works for almost all companies)
	return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s.com&sz=128", strings.ToLower(symbol))
}

// getCoinGeckoImagePath returns the image path segment for a given CoinGecko ID
func getCoinGeckoImagePath(geckoID string) string {
	// Map of known image paths
	paths := map[string]string{
		"bitcoin":   "1/large",
		"ethereum":  "279/large",
		"solana":    "4128/large",
		"xrp":       "44/large",
		"dogecoin":  "5/large",
		"cardano":   "975/large",
		"avalanche": "12559/large",
		"polkadot":  "12171/large",
		"chainlink": "877/large",
		"uniswap":   "12504/large",
		"cosmos":    "1481/large",
		"litecoin":  "2/large",
		"shiba-inu": "11939/large",
		"pepe":      "29850/large",
		"polygon":   "4713/large",
	}
	if path, ok := paths[geckoID]; ok {
		return path
	}
	return "1/large" // default
}
