package robinhood

import (
	"context"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	clientID    = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"
	contentType = "application/json"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/%d Firefox/75.0"
)

type Client interface {
	GetInstrument(ctx context.Context, id string) (*Instrument, error)
	GetMarket(ctx context.Context, id string) (*Market, error)
	GetOrders(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseOrders, error)
	GetPositions(ctx context.Context, auth *ResponseToken, cursor string) (*ResponsePositions, error)
	GetToken(ctx context.Context, username, password, mfa string) (*ResponseToken, error)
}

func New() Client {
	return &defaultClient{
		c: resty.New().
			SetContentLength(true).
			SetHeader("Content-Type", contentType).
			SetHeader("User-Agent", userAgent).
			SetHeader("Accept-Language", "en-US,en").
			SetError(&ResponseError{}),
	}
}

type defaultClient struct {
	c *resty.Client
}

func (dc *defaultClient) GetToken(ctx context.Context, username, password, mfa string) (*ResponseToken, error) {
	resp, err := dc.c.R().
		SetContext(ctx).
		SetBody(&RequestToken{
			GrantType:                    "password",
			Scope:                        "internal",
			ClientID:                     clientID,
			ExpiresIn:                    86400,
			DeviceToken:                  uuid.New().String(),
			Username:                     username,
			Password:                     password,
			MFACode:                      mfa,
			LongSession:                  true,
			CreateReadOnlySecondaryToken: true,
		}).
		SetResult(&ResponseToken{}).
		Post(EndpointToken)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(error)
	}
	return resp.Result().(*ResponseToken), nil
}

func (dc *defaultClient) GetMarket(ctx context.Context, id string) (*Market, error) {
	result, err := dc.get(ctx, nil, getDetailURL(EndpointMarket, id), &Market{})
	if err != nil {
		return nil, err
	}
	return result.(*Market), nil
}

func (dc *defaultClient) GetOrders(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseOrders, error) {
	result, err := dc.list(ctx, auth, EndpointOrders, cursor, &ResponseOrders{})
	if err != nil {
		return nil, err
	}
	return result.(*ResponseOrders), nil
}

func (dc *defaultClient) GetPositions(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponsePositions, error) {
	result, err := dc.list(ctx, auth, EndpointPositions, cursor, &ResponsePositions{})
	if err != nil {
		return nil, err
	}
	return result.(*ResponsePositions), nil
}

func (dc *defaultClient) GetInstrument(ctx context.Context, id string) (*Instrument, error) {
	result, err := dc.get(ctx, nil, getDetailURL(EndpointInstrument, id), &Instrument{})
	if err != nil {
		return nil, err
	}
	return result.(*Instrument), nil
}

func (dc *defaultClient) get(
	ctx context.Context,
	auth *ResponseToken,
	url string,
	result interface{},
) (interface{}, error) {
	r := dc.c.R().SetContext(ctx).SetResult(result)
	if auth != nil {
		r = r.SetAuthScheme(auth.TokenType).SetAuthToken(auth.AccessToken)
	}
	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, resp.Error().(error)
	}
	return resp.Result(), nil
}

func (dc *defaultClient) list(
	ctx context.Context,
	auth *ResponseToken,
	endpoint, cursor string,
	result interface{},
) (interface{}, error) {
	return dc.get(ctx, auth, getListURL(endpoint, cursor), result)
}

func getDetailURL(prefix, id string) string {
	if isURL(id) {
		return id
	}
	return prefix + id
}

func getListURL(endpoint, cursor string) string {
	if cursor == "" {
		return endpoint
	}
	return cursor
}

// isURL returns true if a given v is url.
func isURL(v string) bool {
	return strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://")
}
