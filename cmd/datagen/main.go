package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"currency_eval/internal/models"
	appLogger "currency_eval/internal/pkg/logger"
	"currency_eval/internal/repository"
	"currency_eval/internal/repository/postgres"
	postgresCurrencyRepo "currency_eval/internal/repository/postgres/currency"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type ExchangeRate struct {
	Date  time.Time          `json:"date"`
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

var currencies = []string{
	"AGLD", "FJD", "LVL", "SCR", "BBD", "HNL", "UGX", "KUJI", "NEAR", "AIOZ", "AUDIO", "WLD", "HNT",
	"ETHFI", "FARM", "SDG", "DGB", "BCH", "AI", "FKP", "JST", "HOT", "AR", "CHILLGUY", "SEI", "QAI",
	"SEK", "BB", "QAR", "JTO", "WEMIX", "G", "FLR", "BIGTIME", "BDT", "T", "LYD", "W", "BDX", "BABYDOGE",
	"SFP", "DIA", "PORK", "JUP", "RSS3", "SNEK", "LYX", "SGB", "SGD", "STRD", "WOO", "BLUR", "STRK", "WZRD",
	"HRK", "DJF", "OMIKAMI", "WAVES", "FLOKI", "SHP", "BGB", "SAND", "DKK", "WCFG", "QUBIC", "BGN", "UMA", "FOX",
	"WSTETH", "TUSD", "LRDS", "RBTC", "HTG", "BHD", "OAS", "PEAQ", "COVAL", "ZETACHAIN", "CGLD", "UNI", "FX", "HUF",
	"BIF", "GF", "PEPECOIN", "WELL", "SWFTC", "BIT", "GT", "SKK", "SKL", "UOS", "WST", "SHPING", "BRETT", "MYTH", "DNT",
	"HT", "TRUAPT", "SLE", "FLOW", "UPI",

	"SLL", "SLP", "ID", "DOG", "GIGA", "DOP", "IO", "UQC", "IQ", "DOT", "OSAK",
	"1INCH", "MAD", "TURBO", "BLD", "UNFI", "APEX", "FTM", "EVER", "POWR", "FTN", "MAV", "CORGIAI", "ARKM", "ATOM",
	"SAVAX", "QUICK", "PENDLE", "BLZ",
	"BOBA", "TONE",
	"BMD", "SNT", "SNX", "SHIRO", "USD", "API3", "ROSE", "SATS", "PORTAL", "SOL", "MUMU", "SOS", "BNB",
	"OGN", "UST",
	"CELR", "BND", "CELO", "TBTC", "AUCTION", "MANTA", "BADGER", "MULTI", "AERO", "CETUS", "SPA", "BNT",
	"SHDW", "BOB",
	"MDL", "OHM", "ME", "SPX", "MDT", "ANDY", "AERGO", "KOIN", "MX", "EGLD", "KAS", "TRUMP", "MEW",
	"PUNDIX", "FXS",
	"AEVO", "DADDYCHILL", "SRD", "NU", "FLUX", "TOPIA", "BPX", "QNT", "OM", "ETHW", "MUSE", "ETHX", "OP",
	"CANTO",
	"MGA", "OKB", "WETH", "SSP", "JUPSOL", "NEON", "SSV", "OKT", "ETH2", "PAAL", "KCS", "BUSD", "ARPA",
	"BRL", "ALCX",
	"ALEX", "STD", "STG", "IOTX", "SHIB", "KDA", "$MICHI", "ZAR", "ALEO", "STN", "BSD", "STX", "QI",
	"CAIR", "ZBC",
	"BLAST", "BSV", "IOST", "SUI", "CAKE", "MSOL", "OMG", "OMI", "PYUSD", "SUN", "BTC", "UYU", "IOTA",
	"ZYPTO", "BTG",
	"OMR", "MIR", "KES", "BTN", "RONIN", "SOLVBTC", "SVC", "BTT", "ONE", "FWOG", "RENDER", "ONG", "CETH",
	"ANKR",
	"ALGO", "SYLO", "UZS", "SC", "SD", "ONT", "DYM", "DYP", "MKD", "DZD", "MKR", "KGS", "ICP", "ZEC",
	"XAF", "NEST",
	"ICX", "XAG", "KMNO", "XAI", "ZEN", "FRIEND", "DOGE", "SXP", "HBAR", "XAU", "MLN", "PEPE", "KHR",
	"IDR", "DOGS",
	"XBC", "CTSI", "BWP", "COMAI", "OUSG", "C98", "OSMO", "NTRN", "SYN", "MMK", "SYP", "CRPT", "GAJ",
	"GAL", "GAS",
	"XCD", "MMX", "BOME", "VR", "XCH", "SYNC", "CBBTC", "ORA", "THETA", "PANDORA", "XCN", "SZL", "ORN",
	"NEXO", "AAVE",
	"MNT", "GBP", "BONK", "BYN", "XDC", "PERP", "BYR", "BONE", "BOND", "MOG", "HYPE", "XDR", "TIME",
	"BICO", "BZD",
	"MOP", "MONA", "HMSTR", "XEC", "PEOPLE", "BZR", "XT", "ZIL", "XEM", "WEETH", "TRAC", "MPL", "10SET",
	"WAXL",
	"OOKI", "SWELL", "BORA", "KMF", "GEL", "ZK", "RSETH", "EETH", "KNC", "PROM", "ALEPH", "PONKE",
	"BODEN", "GFI",
	"MRO", "MRS", "MRU", "BORG", "SUPEROETHB", "WAXP", "SUKU", "GGP", "GRIN", "VEF", "BTRFLY", "ZMK",
	"KARRAT", "ALPH",
	"MOBILE", "TAO", "MTL", "VET", "VES", "ZMW", "LUCE", "USDT", "USDC.E", "OXT", "USDS", "USDP", "ILS",
	"ILV", "GHS",
	"KPW", "TITS", "USDY", "EDU", "MEDIA", "KEEP", "CAD", "CAF", "EEK", "MUR", "IMP", "USD0", "GIP",
	"BEAM", "IMX",
	"CAT", "USDE", "BARSIK", "USDD", "USDC", "USDB", "XMON", "RETH", "INJ", "KRL", "VGX", "CHEX", "MVR",
	"TRIBE",
	"INR", "INV", "KRW", "STSOL", "MWC", "XLM", "DORA", "MWK", "EIGEN", "SUPER", "KSM", "RNDR", "GALA",
	"EGP", "RAD",
	"MOVE", "MXC", "TEL", "MOVR", "RAI", "XMR", "MXN", "TET", "MCOIN", "CDF", "GLM", "RAY", "BTC.B",
	"VXOR", "CDT",
	"FDUSD", "ZRO", "KUB", "NPXS", "SSOL", "IQD", "GMD", "RBN", "ZRX", "MYR", "CEL", "XOF", "GMT",
	"SWETH", "CET",
	"GMX", "OMNI", "GNF", "MZN", "CFG", "IRR", "GNO", "GNT", "GNS", "XPD", "THB", "XPF", "VANRY",
	"BITCOIN", "ABT",
	"CFX", "CDAI", "KWD", "VELO", "BINK", "XPT", "ISK", "ACH", "MINA", "TIA", "VTHO", "DRIFT", "PAB",
	"ACS", "ACT",
	"ACX", "REI", "REN", "ELA", "REP", "ADA", "ELF", "REQ", "STORJ", "VIRTUAL", "CHF", "RARI", "ELG",
	"ULTIMA", "RARE",
	"LADYS", "PAXG", "PAX", "XRD", "CHR", "VND", "CHZ", "KYD", "XRP", "JASMY", "INDEX", "TJS", "AED",
	"FIDA", "H2O",
	"EML", "ZWD", "OCEAN", "QGOLD", "ZWL", "PCI", "ENA", "RGT", "ENJ", "TKX", "KZT", "YFII", "DIMO",
	"GRT", "HBTC",
	"AFN", "TFUEL", "ENS", "KAIA", "DEGEN", "CKB", "LUNC", "XTZ", "LUNA", "AURORA", "AGI", "EOS", "GST",
	"FORT", "RIF",
	"NAD", "FRXETH", "TMM", "SLERF", "GTC", "PEN", "SOLO", "TMT", "CLF", "TOSHI", "EUROC", "SUNDOG",
	"GTQ", "CLP",
	"TND", "CLV", "$MYRO", "XVS", "MEME", "LYXE", "SFUND", "TON", "TOP", "PGK", "TOR", "PNUT", "GYEN",
	"CNH", "NCT",
	"WLUNA", "ERN", "VENOM", "VOXEL", "RLC", "RLB", "MAGA", "CNY", "PHP", "RLY", "OX_OLD", "COQ", "COP",
	"HOPR", "AKT",
	"COW", "GLMR", "ORAI", "XYO", "ETB", "GXC", "ETC", "PIP", "VUV", "LAK", "ETH", "NEO", "MEMES",
	"STEAKEURCV", "ALL",
	"MAVIA", "HIGH", "TRB", "ALT", "ORDI", "GYD", "TRU", "AMD", "GRASS", "DREP", "ETHDYDX", "TRY", "LBP",
	"TRX", "NFT",
	"EUR", "AMP", "ORCA", "USTC", "RON", "NGN", "CRC", "PKR", "CRE", "LUSD", "ANG", "SPELL", "LCX", "CRO",
	"PLA",
	"TTD", "SFRXETH", "CRV", "MNDE", "ANT", "BAKE", "RPL", "AOA", "PLN", "AZERO", "LDO", "MAGIC", "ALICE",
	"CORECHAIN",
	"PLU", "SEAM", "AMAPT", "ZEREBRO", "CTC", "NIO", "APE", "LEO", "APL", "MCO2", "00", "MATIC", "APT",
	"APU", "CTX",
	"PNG", "TVK", "USYC", "CUC", "SOLVBTC.BBN", "PYTH", "AI16Z", "CUP", "TWD", "RSD", "FRAX", "VRSC",
	"WBETH", "METH",
	"BAND", "POL", "ASTR", "NKN", "RSR", "TWT", "ARB", "CVC", "VARA", "CVE", "ARK", "LOKA", "ARS", "CVX",
	"LBTC",
	"MPLX", "SUSHI", "WBTC", "ASM", "RUB", "AST", "BNSOL", "MANA", "CSPR", "AGENTFUN", "ATA", "NMR",
	"STAU", "JEP",
	"NMT", "ATH", "LIT", "TNSR", "POLYX", "DESO", "LOOM", "RVN", "PRO", "TZS", "PRQ", "ONDO", "AUD",
	"RWF", "KAVA",
	"NOK", "STRAX", "LKR", "NOT", "NOS", "CZK", "AVT", "NPC", "EURC", "YER", "LSETH", "MASK", "AWG",
	"NPR", "PRIME",
	"YFI", "MOODENG", "CWBTC", "POPCAT", "WBT", "LQTY", "OLAS", "ALUSD", "PUPS", "PIXEL", "ZETA", "AXL",
	"BILLY",
	"CHEEL", "YGG", "FARTCOIN", "AXS", "HONEY", "0X0", "COMP", "HFT", "WAMPL", "CMETH", "RUNE", "DEXT",
	"FORTH",
	"GHST", "MATH", "IDEX", "DEXE", "VNST", "AVAX", "AZN", "AMPL", "GOAT", "UAH", "BANANA", "WEN",
	"NEIRO", "LPT",
	"EZETH", "GODS", "EDUM", "PYG", "JMD", "XAUT", "PYR", "BTRST", "MKUSD", "DAG", "LQR", "DAI", "DAO",
	"AVAIL", "DAR",
	"FET", "CBETH", "LRC", "REPV2", "LRD", "CORE", "DASH", "POKT", "SAGA", "LSD", "JOD", "GUSD", "HKD",
	"JOE",
	"ARSMEP", "LSL", "LSK", "PUMPBTC", "SAFE", "RDOG", "DCR", "WIF", "LTC", "METIS", "ECOIN", "STETH",
	"NXM", "LTL",
	"SAR", "DYDX", "AGIX", "MUBI", "POND", "JPY", "SBD", "DDX", "LINK", "QTUM", "WILD", "POLS", "FIL",
	"POLY", "EBTC",
	"BAL", "BAM", "BAN", "FIS", "BAT", "NZD", "COTI",
}
var wg sync.WaitGroup

func parseAndUpdate(v ExchangeRate, repo repository.CurrencyRepository, egCtx context.Context) error {
	var data []models.Currency
	for currencyName, currencyValue := range v.Rates {
		d := models.Currency{
			Name:  currencyName,
			Value: currencyValue,
			Date:  v.Date,
		}
		data = append(data, d)
	}
	d := models.Currency{
		Name:  v.Base,
		Value: float64(1),
		Date:  v.Date,
	}
	data = append(data, d)
	if err := repo.Update(egCtx, data); err != nil {
		return err
	}
	fmt.Println(time.Now())
	return nil
}

func generateData(startDate, endDate time.Time, repo repository.CurrencyRepository) {
	currentDate := startDate

	genDataChannel := make(chan ExchangeRate)

	go func() {
		for !currentDate.After(endDate) {
			rates := make(map[string]float64)
			for _, currency := range currencies {
				rates[currency] = rand.Float64()*199.5 + 0.5
			}
			r := ExchangeRate{
				Date:  currentDate,
				Base:  "USD",
				Rates: rates,
			}
			genDataChannel <- r
			currentDate = currentDate.AddDate(0, 0, 1)
		}
		close(genDataChannel)
	}()

	eg, egCtx := errgroup.WithContext(context.Background())
	eg.SetLimit(100)
	for v := range genDataChannel {
		v := v
		eg.Go(func() error {
			return parseAndUpdate(v, repo, egCtx)
		})
	}

	if err := eg.Wait(); err != nil {
		fmt.Printf("error during update: %v\n", err)
	}

	wg.Wait()
	fmt.Println("all updated")
}

func main() {
	absPath, _ := filepath.Abs(".")
	log.Println("current working directory", zap.String("path", absPath))

	conf, err := NewConf(absPath)
	if err != nil {
		log.Fatalf("failed to launch currency_app config %v", err)
	}
	logger, err := appLogger.NewLogger(conf.LogLevel)
	if err != nil {
		log.Fatalf("Failed to launch currency_app logger %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to Sync logger %v", err)
		}
	}()

	pgConf := postgres.Config{
		PostgresHost:     conf.PostgresHost,
		PostgresUser:     conf.PostgresUser,
		PostgresPassword: conf.PostgresPassword,
		PostgresDB:       conf.PostgresDB,
		PostgresPort:     conf.PostgresPort,
		MaxOpenConns:     100,
		MaxIdleConns:     10,
		ConnMaxLifetime:  time.Second * 5,
	}

	currencyRepository, err := postgresCurrencyRepo.NewCurrencyRepository(logger.Named("postgresRepo"), pgConf)
	if err != nil {
		logger.Fatal("failed to launch currency repository", zap.Error(err))
	}

	startDate := time.Date(2000, 12, 12, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 12, 0, 0, 0, 0, time.UTC)
	generateData(startDate, endDate, currencyRepository)
}
