package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func getInstrumentsMarketIds(instruments []*robinhood.Instrument) []string {
	return utils.GetIDs(instruments, func(instrument *robinhood.Instrument) string {
		return instrument.Market
	})
}

func loadAccounts(
	ctx context.Context,
	client robinhood.Client,
	auth *robinhood.ResponseToken,
	ids []string,
) ([]*robinhood.Account, error) {
	return utils.LoadDetails(ctx, ids, func(ctx context.Context, id string) (*robinhood.Account, error) {
		return client.GetAccount(ctx, auth, id)
	})
}

func loadInstruments(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Instrument, error) {
	return utils.LoadDetails(ctx, ids, client.GetInstrument)
}

func loadMarkets(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Market, error) {
	return utils.LoadDetails(ctx, ids, client.GetMarket)
}

func getAccountByURL(accounts []*robinhood.Account) map[string]*robinhood.Account {
	accountByURL := make(map[string]*robinhood.Account, len(accounts))
	for _, account := range accounts {
		accountByURL[account.URL] = account
	}
	return accountByURL
}

func getInstrumentByURL(instruments []*robinhood.Instrument) map[string]*robinhood.Instrument {
	instrumentByURL := make(map[string]*robinhood.Instrument, len(instruments))
	for _, instrument := range instruments {
		instrumentByURL[instrument.URL] = instrument
	}
	return instrumentByURL
}

func getMarketByURL(markets []*robinhood.Market) map[string]*robinhood.Market {
	marketByURL := make(map[string]*robinhood.Market, len(markets))
	for _, market := range markets {
		marketByURL[market.URL] = market
	}
	return marketByURL
}

func readLine(reader *bufio.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", errors.New("line is too long")
	}
	return string(line), nil
}

func getAuthToken(
	ctx context.Context,
	client robinhood.Client,
	username string,
) (*robinhood.ResponseToken, error) {
	password := ""
	reader := bufio.NewReader(os.Stdin)
	for username == "" || password == "" {
		if username == "" {
			fmt.Print("Enter username (email): ")
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			username = strings.TrimSpace(line)
			if username == "" {
				logError.Println("username (email) is required")
				continue
			}
		}
		if password == "" {
			fmt.Print("Enter password: ")
			line, err := term.ReadPassword(0)
			if err != nil {
				return nil, fmt.Errorf("ERROR: %w", err)
			}
			password = strings.TrimSpace(string(line))
			fmt.Println() // NOTE: term.ReadPassword doesn't add new line after enter
			if password == "" {
				logError.Println("password is required")
				continue
			}
		}
	}
	var (
		err  error
		resp *robinhood.ResponseToken
		mfa  = ""
	)
	for resp == nil || resp.AccessToken == "" {
		msg := "Trying to log in using username, password"
		if mfa != "" {
			msg += " and OTP code"
		}
		logVerbose.Println(msg)
		resp, err = client.GetToken(ctx, username, password, mfa)
		if err != nil {
			return nil, err
		}
		if resp.MFARequired {
			fmt.Print("Enter OTP code: ")
			line, err := readLine(reader)
			if err != nil {
				return nil, err
			}
			mfa = strings.TrimSpace(line)
			if mfa == "" {
				logError.Println("OTP code is required")
				continue
			}
		}
	}
	logVerbose.Println("Successfully logged in")
	return resp, nil
}

func getAuthTokenFromFile(pathToTokenFile string) (*robinhood.ResponseToken, error) {
	data, err := os.ReadFile(pathToTokenFile)
	if err != nil {
		return nil, err
	}
	resp := &robinhood.ResponseToken{}
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func outputItemsJSON(w io.Writer, items []map[string]any) error {
	encoder := json.NewEncoder(w)
	for _, item := range items {
		if err := encoder.Encode(item); err != nil {
			return err
		}
	}
	return nil
}

func loadItems[T any](
	ctx context.Context,
	getItemsFunc func(ctx context.Context, cursor string) (*robinhood.ResponseList[T], error),
) ([]T, error) {
	items, err := utils.LoadList(ctx, func(c context.Context, cursor string) ([]T, string, error) {
		result, err := getItemsFunc(c, cursor)
		if err != nil {
			return nil, "", err
		}
		return result.Results, result.Next, nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}
