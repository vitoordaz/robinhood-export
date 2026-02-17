package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func doDividends(args arguments) {
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

	logVerbose.Println("loading dividends")
	switch args.format {
	case "json":
		dividends, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointDividends, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d dividends\n", len(dividends))
		outputFunc = func(f *os.File) error { return outputItemsJSON(f, dividends) }
	case "csv":
		dividends, err := loadDividends(ctx, client, token)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d dividends\n", len(dividends))

		logVerbose.Println("loading accounts")
		accounts, err := loadAccounts(ctx, client, token, getDividendsAccountIds(dividends))
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d accounts\n", len(accounts))

		logVerbose.Println("loading instruments")
		instruments, err := loadInstruments(ctx, client, getDividendsInstrumentIds(dividends))
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d instruments\n", len(instruments))
		outputFunc = func(f *os.File) error { return outputDividendsCSV(f, dividends, accounts, instruments) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}

func loadDividends(
	ctx context.Context,
	client robinhood.Client,
	token *robinhood.ResponseToken,
) ([]*robinhood.Dividend, error) {
	return utils.LoadList(ctx, func(ctx context.Context, cursor string) ([]*robinhood.Dividend, string, error) {
		result, err := client.GetDividends(ctx, token, cursor)
		if err != nil {
			return nil, "", err
		}
		return result.Results, result.Next, nil
	})
}

func getDividendsInstrumentIds(dividends []*robinhood.Dividend) []string {
	return utils.GetIDs(dividends, func(dividend *robinhood.Dividend) string {
		return dividend.Instrument
	})
}

func getDividendsAccountIds(dividends []*robinhood.Dividend) []string {
	return utils.GetIDs(dividends, func(dividend *robinhood.Dividend) string {
		return dividend.Account
	})
}

func outputDividendsCSV(
	w io.Writer,
	dividends []*robinhood.Dividend,
	accounts []*robinhood.Account,
	instruments []*robinhood.Instrument,
) error {
	accountByURL := getAccountByURL(accounts)
	instrumentByURL := getInstrumentByURL(instruments)
	header := []string{
		"account",
		"symbol",
		"state",
		"amount",
		"rate",
		"withholding",
		"nra_withholding",
		"record_date",
		"payable_date",
		"paid_at",
	}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, dividend := range dividends {
		account, ok := accountByURL[dividend.Account]
		if !ok {
			return fmt.Errorf("missing account: %s", dividend.Account)
		}
		instrument, ok := instrumentByURL[dividend.Instrument]
		if !ok {
			return fmt.Errorf("missing instrument: %s", dividend.Instrument)
		}
		record := []string{
			account.AccountNumber,
			instrument.Symbol,
			dividend.State,
			dividend.Amount,
			dividend.Rate,
			dividend.Withholding,
			dividend.NRAWithholding,
			dividend.RecordDate,
			dividend.PayableDate,
			dividend.PaidAt,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
