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

	token, err := getAuthToken(ctx, client, args.username)
	if err != nil {
		logError.Fatalln(err)
	}

	logVerbose.Println("loading dividends")
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

	var f *os.File
	if args.output == "" {
		f = os.Stdout
	} else {
		if f, err = os.Create(args.output); err != nil {
			logError.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				logError.Println(err) // still exit with code 0, because this error is not critical
			}
		}()
	}
	if err := outputDividends(f, dividends, accounts, instruments); err != nil {
		logError.Fatalln(err)
	}
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

func outputDividends(
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
