package robinhood

import "context"

type MockClient struct {
	GetAccountFunc    func(ctx context.Context, auth *ResponseToken, id string) (*Account, error)
	GetDividendsFunc  func(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Dividend], error)
	GetInstrumentFunc func(ctx context.Context, id string) (*Instrument, error)
	GetMarketFunc     func(ctx context.Context, id string) (*Market, error)
	GetOrdersFunc     func(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Order], error)
	GetPositionsFunc  func(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Position], error)
	GetTokenFunc      func(ctx context.Context, username, password, mfa string) (*ResponseToken, error)
}

func (c *MockClient) GetAccount(ctx context.Context, auth *ResponseToken, id string) (*Account, error) {
	if c.GetAccountFunc != nil {
		return c.GetAccountFunc(ctx, auth, id)
	}
	return &Account{}, nil
}

func (c *MockClient) GetDividends(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Dividend], error) {
	if c.GetInstrumentFunc != nil {
		return c.GetDividendsFunc(ctx, auth, cursor)
	}
	return &ResponseList[*Dividend]{}, nil
}

func (c *MockClient) GetInstrument(ctx context.Context, id string) (*Instrument, error) {
	if c.GetInstrumentFunc != nil {
		return c.GetInstrumentFunc(ctx, id)
	}
	return &Instrument{}, nil
}

func (c *MockClient) GetMarket(ctx context.Context, id string) (*Market, error) {
	if c.GetMarketFunc != nil {
		return c.GetMarketFunc(ctx, id)
	}
	return &Market{}, nil
}

func (c *MockClient) GetOrders(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Order], error) {
	if c.GetOrdersFunc != nil {
		return c.GetOrdersFunc(ctx, auth, cursor)
	}
	return &ResponseList[*Order]{}, nil
}

func (c *MockClient) GetPositions(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Position], error) {
	if c.GetPositionsFunc != nil {
		return c.GetPositionsFunc(ctx, auth, cursor)
	}
	return &ResponseList[*Position]{}, nil
}

func (c *MockClient) GetToken(ctx context.Context, username, password, mfa string) (*ResponseToken, error) {
	if c.GetTokenFunc != nil {
		return c.GetTokenFunc(ctx, username, password, mfa)
	}
	return &ResponseToken{}, nil
}
