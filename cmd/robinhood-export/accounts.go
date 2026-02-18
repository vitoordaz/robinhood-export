package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strconv"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func doAccounts(args arguments) {
	if !args.verbose {
		logVerbose.SetOutput(io.Discard) // disable verbose logging
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := robinhood.New()

	var (
		token *robinhood.ResponseToken
		err   error
	)
	if args.pathToTokenFile == "" {
		token, err = getAuthToken(ctx, client, args.username)
	} else {
		token, err = getAuthTokenFromFile(args.pathToTokenFile)
	}
	if err != nil {
		logError.Fatalln(err)
	}

	var outputFunc func(f *os.File) error

	logVerbose.Println("loading accounts")
	switch args.format {
	case "json":
		accounts, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointAccounts, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d accounts\n", len(accounts))
		outputFunc = func(f *os.File) error { return outputItemsJSON(f, accounts) }
	case "csv":
		accounts, err := listAccounts(ctx, client, token)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d accounts\n", len(accounts))
		outputFunc = func(f *os.File) error { return outputAccountsCSV(f, accounts) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}

func listAccounts(
	ctx context.Context,
	client robinhood.Client,
	token *robinhood.ResponseToken,
) ([]*robinhood.Account, error) {
	return utils.LoadList(
		ctx,
		func(ctx context.Context, cursor string) ([]*robinhood.Account, string, error) {
			result, err := client.ListAccounts(ctx, token, cursor)
			if err != nil {
				return nil, "", err
			}
			return result.Results, result.Next, nil
		},
	)
}

func outputAccountsCSV(w io.Writer, accounts []*robinhood.Account) error {
	header := []string{"account number", "user", "type", "brokerage account type", "state", "cash", "locked"}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, account := range accounts {
		record := []string{
			account.AccountNumber,
			account.User,
			account.Type,
			account.BrokerageAccountType,
			account.State,
			account.Cash,
			strconv.FormatBool(account.Locked),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
