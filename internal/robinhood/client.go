package robinhood

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	// This is a client id that used for browser.
	clientID    = "c82SH0WZOsabOXGP2sxqcj34FxkvfnWRZBKlBjFS"
	contentType = "application/json"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:75.0) Gecko/%d Firefox/75.0"
)

type Client interface {
	GetDividends(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Dividend], error)
	GetInstrument(ctx context.Context, id string) (*Instrument, error)
	GetMarket(ctx context.Context, id string) (*Market, error)
	GetOrders(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Order], error)
	GetPositions(ctx context.Context, auth *ResponseToken, cursor string) (*ResponseList[*Position], error)
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

func (dc *defaultClient) GetDividends(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Dividend], error) {
	return doList[ResponseList[*Dividend]](ctx, dc.c, auth, EndpointDividends, cursor)
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
	return handleResponse[ResponseToken](resp)
}

func (dc *defaultClient) GetMarket(ctx context.Context, id string) (*Market, error) {
	return doGet[Market](ctx, dc.c, nil, getDetailURL(EndpointMarket, id))
}

func (dc *defaultClient) GetOrders(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Order], error) {
	return doList[ResponseList[*Order]](ctx, dc.c, auth, EndpointOrders, cursor)
}

func (dc *defaultClient) GetPositions(
	ctx context.Context,
	auth *ResponseToken,
	cursor string,
) (*ResponseList[*Position], error) {
	return doList[ResponseList[*Position]](ctx, dc.c, auth, EndpointPositions, cursor)
}

func (dc *defaultClient) GetInstrument(ctx context.Context, id string) (*Instrument, error) {
	return doGet[Instrument](ctx, dc.c, nil, getDetailURL(EndpointInstrument, id))
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

func handleResponse[T any](resp *resty.Response) (*T, error) {
	if resp.IsError() {
		err, ok := resp.Error().(error)
		if ok {
			return nil, err
		}
		return nil, fmt.Errorf("invalid error type: %s", reflect.TypeOf(resp.Error()))
	}
	result, ok := resp.Result().(*T)
	if ok {
		return result, nil
	}
	return nil, fmt.Errorf("invalid result type: %s", reflect.TypeOf(resp.Result()))
}

func doGet[T any](ctx context.Context, client *resty.Client, auth *ResponseToken, url string) (*T, error) {
	r := client.R().SetContext(ctx).SetResult(new(T))
	if auth != nil {
		r = r.SetAuthScheme(auth.TokenType).SetAuthToken(auth.AccessToken)
	}
	resp, err := r.Get(url)
	if err != nil {
		return nil, err
	}
	return handleResponse[T](resp)
}

func doList[T any](
	ctx context.Context,
	client *resty.Client,
	auth *ResponseToken,
	endpoint, cursor string,
) (*T, error) {
	return doGet[T](ctx, client, auth, getListURL(endpoint, cursor))
}
