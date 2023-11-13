package saxotrader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"io"

	"net/http"
)

type SaxoAPI struct {
	Endpoint   string
	Params     map[string]string
	ClientKey  string
	AccountKey string
	LoginToken string
	Body       []byte
	BodyObject interface{}
}

type RESTCall struct {
	Method string
	Path   string
}

var SaxoEndpoints = map[string]RESTCall{
	"user":        {"GET", "port/v1/users/me"},
	"balance":     {"GET", "port/v1/balances"},
	"client":      {"GET", "port/v1/clients/me"},
	"account":     {"GET", "port/v1/accounts/me"},
	"instruments": {"GET", "ref/v1/instruments"},
	"prices":      {"GET", "trade/v1/infoprices/list"},
	"trade":       {"POST", "trade/v2/orders"},
	"orderlist":   {"GET", "port/v1/orders/me"},
	"positions":   {"Get", "port/v1/positions/me"},
}

type SaxoQuote struct {
	AskPrice         float64
	AskSize          float64
	BidPrice         float64
	BidSize          float64
	Amount           float64
	DelayedByMinutes float64
	ErrorCode        string
	MarketState      string
	Mid              float64
	PriceSource      string
	PriceSourceType  string
	PriceTypeAsk     string
	PriceTypeBid     string
}

type SaxoPrice struct {
	AssetType        string
	LastUpdated      string
	PriceSource      string
	Quote            SaxoQuote
	DisplayAndFormat SaxoFormat
	Uic              int
}

type SaxoFormat struct {
	Format        string
	Currency      string
	Decimals      int
	Description   string
	OrderDecimals int
	Symbol        string
}

type SaxoUser struct {
	Active                            bool
	ClientKey                         string
	Culture                           string
	Language                          string
	LastLoginTime                     string
	LastLoginStatus                   string
	LegalAssetTypes                   []string
	MarketDataViaOpenApiTermsAccepted bool
	Name                              string
	TimezoneId                        int
	UserId                            string
	UserKey                           string
}

type SaxoClient struct {
	AccountValueProtectionLimit         float64
	AllowedNettingProfiles              []string
	AllowedTradingSessions              string
	ClientId                            string
	ClientKey                           string
	ClientType                          string
	CurrencyDecimals                    int
	DefaultAccountKey                   string
	DefaultAccountId                    string
	DefaultCurrency                     string
	ForceOpenDefaultValue               bool
	IsMarginTradingAllowed              bool
	IsVariationMarginEligible           bool
	LegalAssetTypes                     []string
	LegalAssetTypesAreIndicative        bool
	MarginCalculationMethod             string
	MarginMonitoringMode                string
	Name                                string
	PartnerPlatformId                   string
	PositionNettingMethod               string
	PositionNettingMode                 string
	PositionNettingProfile              string
	ReduceExposureOnly                  bool
	SupportsAccountValueProtectionLimit bool
}

type SaxoAccount struct {
	AccountId                             string
	AccountKey                            string
	AccountGroupKey                       string
	AccountName                           string
	AccountType                           string
	AccountSubType                        string
	AccountValueProtectionLimit           float64
	AccountValueProtectionLimitCurrency   string
	Active                                bool
	CanUseCashPositionsAsMarginCollateral bool
	CfdBorrowingCostsActive               bool
	ClientId                              string
	ClientKey                             string
	CreationDate                          string
	Currency                              string
	CurrencyDecimals                      int
	DirectMarketAccess                    bool
	FractionalOrderEnabled                bool
	FractionalOrderEnabledAssetTypes      []string
	IndividualMargining                   bool
	IsCurrencyConversionAtSettlementTime  bool
	IsMarginTradingAllowed                bool
	IsShareable                           bool
	IsTrialAccount                        bool
	LegalAssetTypes                       []string
	ManagementType                        string
	MarginCalculationMethod               string
	MarginLendingEnabled                  bool
	PortfolioBasedMarginEnabled           bool
	Sharing                               []string
	SupportsAccountValueProtectionLimit   bool
	UseCashPositionsAsMarginCollateral    bool
}

type SaxoBalance struct {
	CalculationReliability  string
	CashAvailableForTrading float64
	CashBalance             float64
	CashBlocked             float64
	ChangesScheduled        bool
	ClosedPositionsCount    int
	CollateralAvailable     float64
	CollateralCreditValue   struct {
		Line           float64
		UtilizationPct float64
	}
	CorporateActionUnrealizedAmounts float64
	CostToClosePositions             float64
	Currency                         string
	CurrencyDecimals                 int
	InititialMargin                  struct {
		CollateralAvailable   float64
		CollateralCreditValue struct {
			Line           float64
			UtilizationPct float64
		}
		MarginAvailable              float64
		MarginCollateralNotAvailable float64
		MarginUsedByCurrentPositions float64
		MarginUtilizationPct         float64
		NetEquityForMargin           float64
		OtherCollateralDeduction     float64
	}
	IsPortfolioMarginModelSimple      bool
	MarginAndCollateralUtilizationPct float64
	MarginAvailableForTrading         float64
	MarginCollateralNotAvailable      float64
	MarginExposureCoveragePct         float64
	MarginNetExposure                 float64
	MarginUsedByCurrentPositions      float64
	MarginUtilizationPct              float64
	NetEquityForMargin                float64
	NetPositionsCount                 int
	NonMarginPositionsValue           float64
	OpenIpoOrdersCount                int
	OpenPositionsCount                int
	OptionPremiumsMarketValue         float64
	OrdersCount                       int
	OtherCollateral                   float64
	SettlementValue                   float64
	SpendingPowerDetail               map[string]string
	TotalValue                        float64
	TransactionsNotBooked             float64
	TriggerOrdersCount                int
	UnrealizedMarginClosedProfitLoss  float64
	UnrealizedMarginOpenProfitLoss    float64
	UnrealizedMarginProfitLoss        float64
	UnrealizedPositionsValue          float64
}

type SaxoAsset struct {
	AssetType    string
	CurrencyCode string
	ExchangeId   string
	Description  string
	GroupId      string
	Idenfitier   string
	SummaryType  string
	Symbol       string
	TradableAs   []string
}

type SaxoOrderInstruction struct {
	Uic           int
	BuySell       string
	AssetType     string
	Amount        float64
	OrderPrice    float64
	OrderType     string
	OrderRelation string
	OrderDuration struct {
		DurationType string
	}
	ManualOrder bool
	AccountKey  string
}

type SaxoOrder struct {
	AccountId  string
	AccountKey string
	AdviceNote string
	Amount     float64
	Ask        float64
	AssetType  string
	Bid        float64
	BuySell    string

	CalculationReliability   string
	ClientId                 string
	ClientKey                string
	ClientName               string
	ClientNote               string
	CorrelationKey           string
	CurrentPrice             float64
	CurrentPriceDelayMinutes float64
	CurrentPriceType         string
	DisplayAndFormat         SaxoFormat
	DistanceToMarket         float64
	Duration                 struct {
		DurationType string
	}
	Exchange struct {
		ExchangeId  string
		Description string
		IsOpen      bool
		TimezoneId  string
	}
	IpoSubscriptionFee     float64
	IsExtendedHoursEnabled bool
	IsForceOpen            bool
	IsMarketOpen           bool
	MarketPrice            float64
	MarketState            string
	MarketValue            float64
	NonTradableReason      string
	OpenOrderType          string
	OrderAmountType        string
	OrderId                string
	OrderRelation          string
	OrderTime              string
	Price                  float64
	RelatedOpenOrders      []string
	Status                 string
	TradingStatus          string
	Uic                    int
}

type SaxoPosition struct {
	DisplayAndFormat SaxoFormat
	NetPositionId    string
	PositionBase     struct {
		Amount                     float64
		AccountId                  string
		AccountKey                 string
		AssetType                  string
		CanBeClosed                bool
		ClientId                   string
		CloseConversionRateSettled bool
		CorrelationKey             string
		ExecutionOpenTime          string
		IsForceOpen                bool
		IsMarketOpen               bool
		LockedByBackOffice         bool
		OpenPrice                  float64
		OpenPriceIncludingCosts    float64
		RelatedOpenOrders          []string
		SourceOrderId              string
		SpotDate                   string
		Status                     string
		Uic                        int
		ValueDate                  string
	}
	PositionId   string
	PositionView struct {
		Ask                             float64
		Bid                             float64
		CalculationReliability          string
		ConversionRateCurrent           float64
		ConversionRateOpen              float64
		CurrentPrice                    float64
		CurrentPriceDelayMinutes        float64
		CurrentPriceTypeId              string
		Exposure                        float64
		ExposureCurrency                string
		ExposureInBaseCurrency          float64
		InstrumentPriceDayPercentChange float64
		MarketState                     string
		MarketValue                     float64
		MarketValueInBaseCurrency       float64
		ProfitLossOnTrade               float64
		ProfitLossOnTradeInBaseCurrency float64
		TradeCostsTotal                 float64
		TradeCostsTotalInBaseCurrency   float64
	}
}

type SaxoInstrument struct {
	Data []SaxoAsset
}

func (api *SaxoAPI) Call(call string) ([]byte, error) {
	uri := api.Endpoint + SaxoEndpoints[call].Path
	var rdr io.Reader
	if api.BodyObject != nil {
		ba, err := json.Marshal(api.BodyObject)
		if err != nil {
			return nil, err
		}
		api.Body = ba
		rdr = bytes.NewReader(ba)
	} else {
		rdr = bytes.NewReader(api.Body)
	}
	req, err := http.NewRequest(SaxoEndpoints[call].Method, uri, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.LoginToken))
	// add content type if we have a body
	if api.BodyObject != nil || len(api.Body) > 0 {
		req.Header.Add("Content-Type", "application/json")
	}
	// for GET requests we need to add the params to the URL
	if SaxoEndpoints[call].Method == "GET" {
		q := req.URL.Query()
		for k, v := range api.Params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	res := http.Request(*req)
	if res.Response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Error: %s", res.Response.Status))
	}
	ba, err := io.ReadAll(res.Response.Body)
	if err != nil {
		return nil, err
	}
	return ba, nil
}

func (api *SaxoAPI) User() (*SaxoUser, error) {
	data, err := api.Call("user")
	if err != nil {
		return nil, err
	}
	var user SaxoUser
	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (api *SaxoAPI) Client() (*SaxoClient, error) {
	data, err := api.Call("client")
	if err != nil {
		return nil, err
	}
	var client SaxoClient
	err = json.Unmarshal(data, &client)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (api *SaxoAPI) Account() (*SaxoAccount, error) {
	data, err := api.Call("account")
	if err != nil {
		return nil, err
	}
	var account SaxoAccount
	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (api *SaxoAPI) Balance() (*SaxoBalance, error) {
	data, err := api.Call("balance")
	if err != nil {
		return nil, err
	}
	var balance SaxoBalance
	err = json.Unmarshal(data, &balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (api *SaxoAPI) Instruments() (*SaxoInstrument, error) {
	data, err := api.Call("instruments")
	if err != nil {
		return nil, err
	}
	var instruments SaxoInstrument
	err = json.Unmarshal(data, &instruments)
	if err != nil {
		return nil, err
	}
	return &instruments, nil
}

func NewSaxoAPICall(loginToken string) *SaxoAPI {
	return &SaxoAPI{Endpoint: "https://gateway.saxobank.com/sim/openapi/", LoginToken: loginToken, Params: make(map[string]string)}
}
