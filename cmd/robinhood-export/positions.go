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

func doPositions(args arguments) {
	if !args.verbose {
		logVerbose.SetOutput(io.Discard) // disable verbose logging
	}

	ctx := context.Background()
	client := robinhood.New()

	token, err := getAuthToken(ctx, client, args.username)
	if err != nil {
		logError.Fatalln(err)
	}

	logVerbose.Println("loading positions")
	positions, err := loadPositions(ctx, client, token)
	if err != nil {
		logError.Fatalln(err)
	}
	if !args.all {
		positions = getOpenPositions(positions)
	}
	logVerbose.Printf("loaded %d positions\n", len(positions))

	logVerbose.Println("loading instruments")
	instruments, err := loadInstruments(ctx, client, getPositionsInstrumentIds(positions))
	if err != nil {
		logError.Fatalln(err)
	}
	logVerbose.Printf("loaded %d instruments\n", len(instruments))

	logVerbose.Println("loading markets")
	markets, err := loadMarkets(ctx, client, getInstrumentsMarketIds(instruments))
	if err != nil {
		logError.Fatalln(err)
	}
	logVerbose.Println("loading markets")

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
	if err := outputPositions(f, positions, instruments, markets); err != nil {
		logError.Fatalln(err)
	}
}

func getPositionsInstrumentIds(positions []*robinhood.Position) []string {
	return utils.GetIDs(positions, func(position interface{}) string {
		return position.(*robinhood.Position).Instrument
	})
}

func loadPositions(
	ctx context.Context,
	client robinhood.Client,
	token *robinhood.ResponseToken,
) ([]*robinhood.Position, error) {
	positions := make([]*robinhood.Position, 0)
	err := utils.LoadList(ctx, &positions, func(c context.Context, cursor string) (interface{}, string, error) {
		result, er := client.GetPositions(c, token, cursor)
		if er != nil {
			return nil, "", er
		}
		return result.Results, result.Next, nil
	})
	if err != nil {
		return nil, err
	}
	return positions, nil
}

// getOpenPositions filters out closed (quantity == 0) positions from a give slice of positions.
func getOpenPositions(positions []*robinhood.Position) []*robinhood.Position {
	result := make([]*robinhood.Position, 0)
	for _, position := range positions {
		if isZero, _ := utils.IsZero(position.Quantity); !isZero {
			result = append(result, position)
		}
	}
	return result
}

func outputPositions(
	w io.Writer,
	positions []*robinhood.Position,
	instruments []*robinhood.Instrument,
	markets []*robinhood.Market,
) error {
	instrumentByURL := getInstrumentByURL(instruments)
	marketByURL := getMarketByURL(markets)
	header := []string{"market", "symbol", "quantity", "avg price"}
	writer := csv.NewWriter(w)
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, position := range positions {
		instrument, ok := instrumentByURL[position.Instrument]
		if !ok {
			return fmt.Errorf("missing instrument: %s", position.Instrument)
		}
		market, ok := marketByURL[instrument.Market]
		if !ok {
			return fmt.Errorf("instrument %s missing market: %s", instrument.Symbol, instrument.Market)
		}
		record := []string{market.Acronym, instrument.Symbol, position.Quantity, position.AverageBuyPrice}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
