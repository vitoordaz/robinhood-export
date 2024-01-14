package robinhood

const (
	EndpointAccounts   = "https://api.robinhood.com/accounts/"
	EndpointDividends  = "https://api.robinhood.com/dividends/"
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

type Account struct {
	URL                                  string `json:"url"`
	PortfolioCash                        string `json:"portfolio_cash"`
	CanDowngradeToCash                   string `json:"can_downgrade_to_cash"`
	User                                 string `json:"user"`
	AccountNumber                        string `json:"account_number"`
	Type                                 string `json:"type"`
	BrokerageAccountType                 string `json:"brokerage_account_type"`
	CreatedAt                            string `json:"created_at"`
	UpdatedAt                            string `json:"updated_at"`
	Deactivated                          bool   `json:"deactivated"`
	DepositHalted                        bool   `json:"deposit_halted"`
	WithdrawalHalted                     bool   `json:"withdrawal_halted"`
	OnlyPositionClosingTrades            bool   `json:"only_position_closing_trades"`
	BuyingPower                          string `json:"buying_power"`
	ONBP                                 string `json:"onbp"`
	CashAvailableForWithdrawal           string `json:"cash_available_for_withdrawal"`
	Cash                                 string `json:"cash"`
	AmountEligibleForDepositCancellation string `json:"amount_eligible_for_deposit_cancellation"`
	CashHeldForOrders                    string `json:"cash_held_for_orders"`
	UnclearedDeposits                    string `json:"uncleared_deposits"`
	SMA                                  string `json:"sma"`
	SMAHeldForOrders                     string `json:"sma_held_for_orders"`
	UnsettledFunds                       string `json:"unsettled_funds"`
	UnsettledDebit                       string `json:"unsettled_debit"`
	CryptoBuyingPower                    string `json:"crypto_buying_power"`
	MaxACHEarlyAccessAmount              string `json:"max_ach_early_access_amount"`
	CashBalances                         string `json:"cash_balances"`
	SweepEnabled                         bool   `json:"sweep_enabled"`
	SweepEnrolled                        bool   `json:"sweep_enrolled"`
	OptionLevel                          string `json:"option_level"`
	IsPinnacleAccount                    bool   `json:"is_pinnacle_account"`
	RHSAccountNumber                     int64  `json:"rhs_account_number"`
	State                                string `json:"state"`
	ActiveSubscriptionID                 string `json:"active_subscription_id"`
	Locked                               bool   `json:"locked"`
	PermanentlyDeactivated               bool   `json:"permanently_deactivated"`
	IPOAccessRestricted                  bool   `json:"ipo_access_restricted"`
	IPOAccessRestrictedReason            string `json:"ipo_access_restricted_reason"`
	ReceivedACHDebitLocked               bool   `json:"received_ach_debit_locked"`
	DripEnabled                          bool   `json:"drip_enabled"`
	EligibleForFractionals               bool   `json:"eligible_for_fractionals"`
	EligibleForDrip                      bool   `json:"eligible_for_drip"`
	EligibleForCashManagement            bool   `json:"eligible_for_cash_management"`
	CashManagementEnabled                bool   `json:"cash_management_enabled"`
	OptionTradingOnExpirationEnabled     bool   `json:"option_trading_on_expiration_enabled"`
	CashHeldForOptionsCollateral         string `json:"cash_held_for_options_collateral"`
	FractionalPositionClosingOnly        bool   `json:"fractional_position_closing_only"`
	UserID                               string `json:"user_id"`
	EquityTradingLock                    string `json:"equity_trading_lock"`
	OptionTradingLock                    string `json:"option_trading_lock"`
	DisableADT                           bool   `json:"disable_adt"`
	ManagementType                       string `json:"management_type"`
	DynamicInstantLimit                  string `json:"dynamic_instant_limit"`
	Affiliate                            string `json:"affiliate"`
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

type Dividend struct {
	ID               string `json:"id,omitempty"`
	URL              string `json:"url,omitempty"`
	Account          string `json:"account,omitempty"`
	Instrument       string `json:"instrument,omitempty"`
	Amount           string `json:"amount,omitempty"`
	Rate             string `json:"rate,omitempty"`
	Position         string `json:"position,omitempty"`
	Withholding      string `json:"withholding,omitempty"`
	RecordDate       string `json:"record_date,omitempty"`
	PayableDate      string `json:"payable_date,omitempty"`
	PaidAt           string `json:"paid_at,omitempty"`
	State            string `json:"state,omitempty"`
	CashDividendID   string `json:"cash_dividend_id,omitempty"`
	DripEnabled      bool   `json:"drip_enabled,omitempty"`
	NRAWithholding   string `json:"nra_withholding,omitempty"`
	IsPrimaryAccount bool   `json:"is_primary_account,omitempty"`
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
