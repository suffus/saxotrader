package saxotrader

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	"user":               {"GET", "port/v1/users/me"},
	"balance":            {"GET", "port/v1/balances"},
	"client":             {"GET", "port/v1/clients/me"},
	"account":            {"GET", "port/v1/accounts/me"},
	"instruments":        {"GET", "ref/v1/instruments"},
	"instrument_details": {"GET", "ref/v1/instruments/details"},
	"prices":             {"GET", "trade/v1/infoprices/list"},
	"make_order":         {"POST", "trade/v2/orders"},
	"order_list":         {"GET", "port/v1/orders/me"},
	"order_details":      {"GET", "port/v1/orders/{ClientKey}/{OrderId}/"},
	"positions":          {"GET", "port/v1/positions/me"},
	"net_positions":      {"GET", "port/v1/netpositions/me"},
	"order":              {"GET", "trade/v2/orders/"},
	"cancel_order":       {"DELETE", "trade/v2/orders/"},
	"replace_order":      {"PUT", "trade/v2/orders/"},
	"quotes":             {"GET", "trade/v1/infoprices/snapshot"},
	"chart":              {"GET", "chart/v1/charts"},
	"chart_data":         {"GET", "chart/v1/charts/"},
	"chart_list":         {"GET", "chart/v1/charts/me"},
	"chart_config":       {"GET", "chart/v1/configurations"},
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
	MarginLendingEnabled                  string
	PortfolioBasedMarginEnabled           bool
	Sharing                               []string
	SupportsAccountValueProtectionLimit   bool
	UseCashPositionsAsMarginCollateral    bool
}

type SaxoAccounts SaxoData[SaxoAccount]

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
	AssetType      string
	CurrencyCode   string
	ExchangeId     string
	Description    string
	GroupId        int
	Identifier     int
	PrimaryListing int
	SummaryType    string
	IssuerCountry  string
	Symbol         string
	TradableAs     []string
}

type SaxoAssetDetails struct {
	AssetType           string
	AmountDecimals      int
	CurrencyCode        string
	DefaultAmount       float64
	DefaultSlippage     float64
	DefaultSlippageType string
	Description         string
	Exchange            struct {
		ExchangeId  string
		Name        string
		CountryCode string
	}
	Format struct {
		Decimals      int
		OrderDecimals int
		Format        string
	}
	FxForwardMaxForwardDate string
	FxForwardMinForwardDate string
	GroupId                 int
	IncrementSize           float64
	IsRedemptionByAmounts   bool
	IsTradable              bool
	NonTradableReason       string
	OrderDistances          struct {
		EntryDefaultDistance          float64
		EntryDefaultDistanceType      string
		LimitDefaultDistance          float64
		LimitDefaultDistanceType      string
		StopLimitDefaultDistance      float64
		StopLimitDefaultDistanceType  string
		StopLossDefaultDistance       float64
		StopLossDefaultDistanceType   string
		StopLossDefaultOrderType      string
		TakeProfitDefaultDistance     float64
		TakeProfitDefaultDistanceType string
		TakeProfitDefaultOrderType    string
	}
	StandardAmounts     []float64
	SupportedOrderTypes []string
	Symbol              string
	TickSize            float64
	TradableAs          []string
	TradableOn          []string
	TradingSignals      string
	TradingStatus       string
	Uic                 int
}

type SaxoInstruction struct {
	ExchangeId                    string
	Keywords                      string
	AssetTypes                    []string
	CanParticipateInMultiLegOrder bool
	TradingStatus                 string
	FieldGroups                   []string
	Class                         []string
	IncludeNonTradable            bool
	Uics                          []int
	Uic                           int
	UnderlyingUic                 int
	ExpiryDates                   string
	OptionSpaceSegment            string
	Tags                          []string
	Amount                        float64
	AmountType                    string
	ForwardDate                   string
	ForwardDateFarLeg             string
	ForwardDateNearLeg            string
	LowerBarrier                  float64
	UpperBarrier                  float64
	OrderBidPrice                 float64
	OrderAskPrice                 float64
	PutCall                       string
	StrikePrice                   float64
	QuoteCurrency                 string
	ToOpenClose                   string
}

type SaxoOrderInstruction struct {
	Uic           int
	BuySell       string
	AssetType     string
	Amount        float64
	OrderPrice    float64
	OrderType     string
	OrderDuration struct {
		DurationType string
	}
	ManualOrder bool
	AccountKey  string
}

type SaxoOrderPost struct {
	Orders          []SaxoOrderInstruction
	ManualOrder     bool
	WaitForApproval bool
	WithAdvice      bool
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

type SaxoNetPosition struct {
	DisplayAndFormat SaxoFormat
	NetPositionId    string
	NetPositionBase  struct {
		Amount                     float64
		AccountId                  string
		AccountKey                 string
		AssetType                  string
		CanBeClosed                bool
		ClientId                   string
		CloseConversionRateSettled bool
		CorrelationKey             string
		HasForceOpenPositions      bool
		IsMarketOpen               bool
		NonTradableReason          string
		NumberOfRelatedOrders      int
		OpeningDirection           string
		OpenIpoOrdersCount         int
		OpenOrdersCount            int
		OpenTriggerOrdersCount     int
		PositionsAccount           string
		SinglePositionStatus       string
		Uic                        int
		ValueDate                  string
	}
	NetPositionView struct {
		AverageOpenPrice                float64
		AverageOpenPriceIncludingCosts  float64
		CalculationReliability          string
		CurrentPrice                    float64
		CurrentPriceDelayMinutes        float64
		CurrentPriceType                string
		Exposure                        float64
		ExposureInBaseCurrency          float64
		InstrumentPriceDayPercentChange float64
		PositionCount                   int
		PositionsNotClosedCount         int
		ProfitLossOnTrade               float64
		Status                          string
		TradeCostsTotal                 float64
		TradeCostsTotalInBaseCurrency   float64
	}
}

type SaxoData[T any] struct {
	Data []T
}

type SaxoAssetSet SaxoData[SaxoAsset]

type SaxoAssetDetailsSet SaxoData[SaxoAssetDetails]

type SaxoError struct {
	ErrorCode  string
	Message    string
	ModelState map[string][]string
}

func (api *SaxoAPI) MakeOrder(amount, price float64, uic int, asset, buysell, duration, orderType string) (SaxoOrderInstruction, error) {
	if duration == "" {
		duration = "DayOrder"
	}
	if orderType == "" {
		orderType = "Limit"
	}
	if api.AccountKey == "" {
		return SaxoOrderInstruction{}, errors.New("No account key set")
	}
	return SaxoOrderInstruction{
		Uic:         uic,
		BuySell:     buysell,
		Amount:      amount,
		AssetType:   asset,
		OrderPrice:  price,
		OrderType:   orderType,
		AccountKey:  api.AccountKey,
		ManualOrder: true,
		OrderDuration: struct {
			DurationType string
		}{DurationType: duration},
	}, nil
}

func (api *SaxoAPI) Call(call string) ([]byte, error) {
	client := http.Client{}
	uri := api.Endpoint + SaxoEndpoints[call].Path
	var rdr io.Reader
	if api.BodyObject != nil {
		ba, err := json.Marshal(api.BodyObject)
		if err != nil {
			return nil, err
		}
		api.Body = ba
		fmt.Println("Body is ", string(ba))
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

	// for GET and DELETE requests we need to add the params to the URL
	if SaxoEndpoints[call].Method == "GET" || SaxoEndpoints[call].Method == "DELETE" {
		q := req.URL.Query()
		for k, v := range api.Params {
			q.Add(k, v)
		}
		if len(api.ClientKey) > 0 {
			q.Add("ClientKey", api.ClientKey)
		}
		if len(api.AccountKey) > 0 {
			q.Add("AccountKey", api.AccountKey)
		}

		paramrx := regexp.MustCompile(`\{([a-zA-Z0-9]+)\}`)
		for _, p := range paramrx.FindAllString(uri, -1) {
			uri = strings.Replace(uri, p, api.Params[strings.Trim(p, "{}")], 1)
		}

		req.URL.RawQuery = q.Encode()
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	ba, err := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		//fmt.Println("Error, result body is ", string(ba))
		return nil, errors.New(fmt.Sprintf("Error: %s %s", res.Status, string(ba)))
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
	api.ClientKey = client.ClientKey
	return &client, nil
}

func (api *SaxoAPI) Accounts() (*SaxoAccounts, error) {
	data, err := api.Call("account")
	if err != nil {
		return nil, err
	}
	var account SaxoAccounts
	err = json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (instr *SaxoInstruction) MakeParams(cmd string) map[string]string {
	params := make(map[string]string)
	if len(instr.AssetTypes) > 0 {
		params["AssetTypes"] = strings.Join(instr.AssetTypes, ",")
	}
	if len(instr.Class) > 0 {
		params["Class"] = strings.Join(instr.Class, ",")
	}
	if len(instr.Tags) > 0 {
		params["Tags"] = strings.Join(instr.Tags, ",")
	}
	if len(instr.FieldGroups) > 0 {
		params["FieldGroups"] = strings.Join(instr.FieldGroups, ",")
	}
	if len(instr.ExchangeId) > 0 {
		params["ExchangeId"] = instr.ExchangeId
	}
	if len(instr.Keywords) > 0 {
		params["Keywords"] = instr.Keywords
	}
	if len(instr.OptionSpaceSegment) > 0 {
		params["OptionSpaceSegment"] = instr.OptionSpaceSegment
	}
	if len(instr.TradingStatus) > 0 {
		params["TradingStatus"] = instr.TradingStatus
	}
	if instr.UnderlyingUic > 0 {
		params["UnderlyingUic"] = strconv.Itoa(instr.UnderlyingUic)
	}
	if instr.Uic > 0 {
		params["Uics"] = strconv.Itoa(instr.Uic)
	}
	if len(instr.Uics) > 0 {
		var strs []string
		for _, i := range instr.Uics {
			strs = append(strs, strconv.Itoa(i))
		}
		params["Uics"] = strings.Join(strs, ",")
	}
	if instr.Amount != 0 {
		params["Amount"] = fmt.Sprintf("%f", instr.Amount)
	}
	if len(instr.AmountType) > 0 {
		params["AmountType"] = instr.AmountType
	}
	if len(instr.ForwardDate) > 0 {
		params["ForwardDate"] = instr.ForwardDate
	}
	if len(instr.ForwardDateFarLeg) > 0 {
		params["ForwardDateFarLeg"] = instr.ForwardDateFarLeg
	}
	if len(instr.ForwardDateNearLeg) > 0 {
		params["ForwardDateNearLeg"] = instr.ForwardDateNearLeg
	}
	if len(instr.PutCall) > 0 {
		params["PutCall"] = instr.PutCall
	}
	if len(instr.QuoteCurrency) > 0 {
		params["QuoteCurrency"] = instr.QuoteCurrency
	}
	if len(instr.ToOpenClose) > 0 {
		params["ToOpenClose"] = instr.ToOpenClose
	}
	return params

}

func (api *SaxoAPI) Balance() (*SaxoBalance, error) {
	if api.ClientKey == "" {
		return nil, errors.New("No client key set")
	}
	data, err := api.Call("balance")
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	var balance SaxoBalance
	err = json.Unmarshal(data, &balance)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (api *SaxoAPI) Instruments(instr SaxoInstruction) ([]SaxoAsset, error) {
	api.Params = instr.MakeParams("instruments")
	api.Params["$top"] = "1000"
	data, err := api.Call("instruments")
	if err != nil {
		return nil, err
	}
	var instruments SaxoAssetSet
	//fmt.Println(string(data))
	err = json.Unmarshal(data, &instruments)
	if err != nil {
		fmt.Println("UME:", err)
		return nil, err
	}
	return instruments.Data, nil
}

func (api *SaxoAPI) InstrumentDetails(instr SaxoInstruction) ([]SaxoAssetDetails, error) {
	if instr.Uic > 0 {
		instr.Uics = append(instr.Uics, instr.Uic)
	}
	api.Params = instr.MakeParams("instrument_details")
	data, err := api.Call("instrument_details")
	if err != nil {
		return nil, err
	}
	var details SaxoAssetDetailsSet
	err = json.Unmarshal(data, &details)
	if err != nil {
		return nil, err
	}
	return details.Data, nil
}

func (api *SaxoAPI) Prices(instr SaxoInstruction) (*SaxoPrice, error) {
	api.Params = instr.MakeParams("prices")
	data, err := api.Call("prices")
	if err != nil {
		return nil, err
	}
	var prices SaxoPrice
	err = json.Unmarshal(data, &prices)
	if err != nil {
		return nil, err
	}
	return &prices, nil
}

func (api *SaxoAPI) NetPositions(instr SaxoInstruction) ([]SaxoNetPosition, error) {
	if api.ClientKey == "" {
		return nil, errors.New("No client key set")
	}
	api.Params["ClientKey"] = api.ClientKey
	data, err := api.Call("net_positions")
	if err != nil {
		return nil, err
	}
	var netpositions SaxoData[SaxoNetPosition]
	err = json.Unmarshal(data, &netpositions)
	if err != nil {
		return nil, err
	}
	return netpositions.Data, nil
}

func (api *SaxoAPI) PlaceOrder(instr SaxoOrderInstruction) ([]SaxoOrder, error) {
	if api.ClientKey == "" {
		return nil, errors.New("No client key set")
	}
	api.BodyObject = instr
	data, err := api.Call("make_order")
	if err != nil {
		fmt.Println(data)
		return nil, err
	}
	var orders SaxoData[SaxoOrder]
	err = json.Unmarshal(data, &orders)
	if err != nil {
		return nil, err
	}
	return orders.Data, nil
}

func (api *SaxoAPI) OrderList() ([]SaxoOrder, error) {
	if api.ClientKey == "" {
		return nil, errors.New("No client key set")
	}
	api.Params["ClientKey"] = api.ClientKey
	data, err := api.Call("order_list")
	if err != nil {
		return nil, err
	}
	var orders SaxoData[SaxoOrder]
	err = json.Unmarshal(data, &orders)
	if err != nil {
		return nil, err
	}
	return orders.Data, nil
}

func NewSaxoAPICall(loginToken string) *SaxoAPI {
	return &SaxoAPI{Endpoint: "https://gateway.saxobank.com/sim/openapi/", LoginToken: loginToken, Params: make(map[string]string)}
}
