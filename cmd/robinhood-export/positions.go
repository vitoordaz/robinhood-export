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

	logVerbose.Println("loading positions")
	switch args.format {
	case "json":
		positions, err := loadItems(
			ctx,
			func(ctx context.Context, cursor string) (*robinhood.ResponseList[map[string]any], error) {
				return client.GetItems(ctx, robinhood.EndpointPositions, token, cursor)
			},
		)
		if err != nil {
			logError.Fatalln(err)
		}
		logVerbose.Printf("loaded %d positions\n", len(positions))
		outputFunc = func(f *os.File) error { return outputItemsJSON(f, positions) }
	case "csv":
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
		outputFunc = func(f *os.File) error { return outputPositionsCSV(f, positions, instruments, markets) }
	default:
		logError.Fatalln("unsupported output format " + args.format)
	}

	outputData(args.output, outputFunc)
}

func getPositionsInstrumentIds(positions []*robinhood.Position) []string {
	return utils.GetIDs(positions, func(position *robinhood.Position) string {
		return position.Instrument
	})
}

func loadPositions(
	ctx context.Context,
	client robinhood.Client,
	token *robinhood.ResponseToken,
) ([]*robinhood.Position, error) {
	positions, err := utils.LoadList(ctx, func(c context.Context, cursor string) ([]*robinhood.Position, string, error) {
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

func outputPositionsCSV(
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
