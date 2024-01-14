package robinhood

const (
	EndpointInstrument = "https://api.robinhood.com/instruments/"
	EndpointMarket     = "https://api.robinhood.com/markets/"
	EndpointToken      = "https://api.robinhood.com/oauth2/token/"
	EndpointOrders     = "https://api.robinhood.com/orders/"
	EndpointPositions  = "https://api.robinhood.com/positions/"
)

type RequestToken struct {
	GrantType                    string `json:"grant_type,omitempty"`
	Scope                        string `json:"scope,omitempty"`
	ClientID                     string `json:"client_id,omitempty"`
	ExpiresIn                    int64  `json:"expires_in,omitempty"`
	DeviceToken                  string `json:"device_token,omitempty"`
	Username                     string `json:"username,omitempty"`
	Password                     string `json:"password,omitempty"`
	MFACode                      string `json:"mfa_code,omitempty"`
	LongSession                  bool   `json:"long_session,omitempty"`
	CreateReadOnlySecondaryToken bool   `json:"create_read_only_secondary_token,omitempty"`
}

type ResponseList[T any] struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	Results  []T    `json:"results,omitempty"`
}

type ResponseToken struct {
	MFARequired                  bool   `json:"mfa_required,omitempty"`
	MFAType                      string `json:"mfa_type,omitempty"`
	AccessToken                  string `json:"access_token,omitempty"`
	ExpiresIn                    int64  `json:"expires_in,omitempty"`
	TokenType                    string `json:"token_type,omitempty"`
	Scope                        string `json:"scope,omitempty"`
	RefreshToken                 string `json:"refresh_token,omitempty"`
	MFACode                      string `json:"mfa_code,omitempty"`
	BackupCode                   string `json:"backup_code,omitempty"`
	ReadOnlySecondaryAccessToken string `json:"read_only_secondary_access_token,omitempty"`
}

type Instrument struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Market     string `json:"market"`
	SimpleName string `json:"simple_name"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Country    string `json:"country"`
	Type       string `json:"type"`
}

const OrderStateFilled = "filled"

type Order struct {
	ID                 string       `json:"id,omitempty"`
	URL                string       `json:"url,omitempty"`
	Instrument         string       `json:"instrument,omitempty"`
	CumulativeQuantity string       `json:"cumulative_quantity,omitempty"`
	AveragePrice       string       `json:"average_price,omitempty"`
	Fees               string       `json:"fees,omitempty"`
	State              string       `json:"state,omitempty"`
	Side               string       `json:"side,omitempty"`
	Price              string       `json:"price"`
	Quantity           string       `json:"quantity"`
	CreatedAt          string       `json:"created_at"`
	UpdatedAt          string       `json:"updated_at"`
	LastTransactionAt  string       `json:"last_transaction_at"`
	Executions         []*Execution `json:"executions,omitempty"`
	TotalNotional      *Notional    `json:"total_notional,omitempty"`
	ExecutedNotional   *Notional    `json:"executed_notional,omitempty"`
}

// IsFilled returns true if order state is filled.
func (o *Order) IsFilled() bool {
	return o.State == OrderStateFilled
}

type Position struct {
	URL             string `json:"url,omitempty"`
	Instrument      string `json:"instrument,omitempty"`
	AverageBuyPrice string `json:"average_buy_price,omitempty"`
	Quantity        string `json:"quantity,omitempty"`
}

type Market struct {
	URL          string `json:"url,omitempty"`
	MIC          string `json:"mic,omitempty"`
	OperatingMIC string `json:"operating_mic,omitempty"`
	Acronym      string `json:"acronym,omitempty"`
	Name         string `json:"name,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	Timezone     string `json:"timezone,omitempty"`
	Website      string `json:"website,omitempty"`
}

type Execution struct {
	ID             string `json:"id,omitempty"`
	Price          string `json:"price,omitempty"`
	Quantity       string `json:"quantity,omitempty"`
	SettlementDate string `json:"settlement_date,omitempty"`
	Timestamp      string `json:"timestamp,omitempty"`
}

type Notional struct {
	Amount       string `json:"amount,omitempty"`
	CurrencyCode string `json:"currency_code,omitempty"`
	CurrencyID   string `json:"currency_id,omitempty"`
}

type ResponseError struct {
	Detail string `json:"detail,omitempty"`
	Err    string `json:"error,omitempty"` // this can be set for /oauth2/token endpoint
}

func (re *ResponseError) Error() string {
	if re.Detail != "" {
		return re.Detail
	}
	if re.Err != "" {
		return re.Err
	}
	return "unexpected error"
}

// GetNotional returns notional amount of a given order or empty string if it's not available.
func GetNotional(order *Order) string {
	if order.TotalNotional != nil {
		return order.TotalNotional.Amount
	}
	if order.ExecutedNotional != nil {
		return order.ExecutedNotional.Amount
	}
	return ""
}

// GetCurrencyCode returns currency code of a given order or empty string if it's not available.
func GetCurrencyCode(order *Order) string {
	if order.TotalNotional != nil {
		return order.TotalNotional.CurrencyCode
	}
	if order.ExecutedNotional != nil {
		return order.ExecutedNotional.CurrencyCode
	}
	return ""
}
