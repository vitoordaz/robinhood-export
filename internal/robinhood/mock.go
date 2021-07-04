package robinhood

import "context"

type MockClient struct {
	GetInstrumentFunc func(ctx context.Context, id string) (*Instrument, error)
	GetMarketFunc     func(ctx context.Context, id string) (*Market, error)
	GetOrdersFunc     func(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseOrders, error)
	GetPositionsFunc  func(ctx context.Context, auth *ResponseToken, cursor string) (*ResponsePositions, error)
	GetTokenFunc      func(ctx context.Context, username, password, mfa string) (*ResponseToken, error)
}

func (c *MockClient) GetInstrument(ctx context.Context, id string) (*Instrument, error) {
	if c.GetInstrumentFunc != nil {
		return c.GetInstrumentFunc(ctx, id)
	}
	return nil, nil
}

func (c *MockClient) GetMarket(ctx context.Context, id string) (*Market, error) {
	if c.GetMarketFunc != nil {
		return c.GetMarketFunc(ctx, id)
	}
	return nil, nil
}

func (c *MockClient) GetOrders(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseOrders, error) {
	if c.GetOrdersFunc != nil {
		return c.GetOrdersFunc(ctx, auth, cursor)
	}
	return nil, nil
}

func (c *MockClient) GetPositions(ctx context.Context, auth *ResponseToken, cursor string) (*ResponsePositions, error) {
	if c.GetPositionsFunc != nil {
		return c.GetPositionsFunc(ctx, auth, cursor)
	}
	return nil, nil
}

func (c *MockClient) GetToken(ctx context.Context, username, password, mfa string) (*ResponseToken, error) {
	if c.GetTokenFunc != nil {
		return c.GetTokenFunc(ctx, username, password, mfa)
	}
	return nil, nil
}
