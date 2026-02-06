package main

import (
	"flag"
	"log"
	"os"
)

type arguments struct {
	username        string // -u, robinhood account username or email
	pathToTokenFile string // -t, path to token file
	output          string // -o, path to output file
	verbose         bool   // -v, enable verbose messages
	all             bool   // -a, return everything
	format          string // -f, output format: json or csv
}

const (
	exitCodeError = 2
	exitCodeOk    = 0
)

var (
	logVerbose = log.New(os.Stdout, "D: ", 0)
	logError   = log.New(os.Stderr, "ERROR: ", 0)

	helpCmd = flag.NewFlagSet("help", flag.ExitOnError)

	dividendsCmd              = flag.NewFlagSet("dividends", flag.ExitOnError)
	dividendsCmdPathTokenFile = dividendsCmd.String("t", "", "path to token file")                   // optional
	dividendsCmdUsername      = dividendsCmd.String("u", "", "Robinhood account username or email.") // optional
	dividendsCmdOutput        = dividendsCmd.String("o", "", "path to output file.")                 // optional
	dividendsCmdVerbose       = dividendsCmd.Bool("v", false, "enable verbose messages.")            // optional
	dividendsCmdFormat        = dividendsCmd.String("f", "json", "output format, json or csv.")      // optional

	ordersCmd              = flag.NewFlagSet("orders", flag.ExitOnError)
	ordersCmdUsername      = ordersCmd.String("u", "", "Robinhood account username or email.") // optional
	ordersCmdPathTokenFile = ordersCmd.String("t", "", "path to token file")                   // optional
	ordersCmdOutput        = ordersCmd.String("o", "", "path to output file.")                 // optional
	ordersCmdVerbose       = ordersCmd.Bool("v", false, "enable verbose messages.")            // optional
	ordersCmdFormat        = ordersCmd.String("f", "json", "output format, json or csv.")      // optional

	positionsCmd         = flag.NewFlagSet("positions", flag.ExitOnError)
	positionsCmdUsername = positionsCmd.String("u", "", "Robinhood account username or email.") // optional
	positionsCmdOutput   = positionsCmd.String("o", "", "path to output file.")                 // optional
	positionsCmdAll      = positionsCmd.Bool("a", false, "return all positions (even closed).") // optional
	positionsCmdVerbose  = positionsCmd.Bool("v", false, "enable verbose messages.")            // optional

	optionsOrdersCmd              = flag.NewFlagSet("options orders", flag.ExitOnError)
	optionsOrdersCmdUsername      = optionsOrdersCmd.String("u", "", "Robinhood account username or email.") // optional
	optionsOrdersCmdPathTokenFile = optionsOrdersCmd.String("t", "", "path to token file")                   // optional
	optionsOrdersCmdOutput        = optionsOrdersCmd.String("o", "", "path to output file.")                 // optional
	optionsOrdersCmdVerbose       = optionsOrdersCmd.Bool("v", false, "enable verbose messages.")            // optional
	optionsOrdersCmdFormat        = optionsOrdersCmd.String("f", "json", "output format, json or csv.")      // optional
)

func main() {
	flag.Usage = printUsage
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(exitCodeError)
	}

	dividendsCmd.Usage = printDividendsUsage
	helpCmd.Usage = printHelpUsage
	ordersCmd.Usage = printOrdersUsage
	positionsCmd.Usage = printPositionsUsage
	optionsOrdersCmd.Usage = printOptionOrdersUsage

	switch flag.Arg(0) {
	case "dividends":
		if err := dividendsCmd.Parse(os.Args[2:]); err != nil {
			logError.Println(err)
			dividendsCmd.Usage()
			os.Exit(exitCodeError)
		}
		doDividends(arguments{
			username:        *dividendsCmdUsername,
			pathToTokenFile: *dividendsCmdPathTokenFile,
			verbose:         *dividendsCmdVerbose,
			output:          *dividendsCmdOutput,
			format:          *dividendsCmdFormat,
		})
		os.Exit(exitCodeError)
	case "help":
		if err := helpCmd.Parse(os.Args[2:]); err != nil || len(helpCmd.Args()) < 1 {
			if err != nil {
				logError.Println(err)
			}
			helpCmd.Usage()
			os.Exit(exitCodeError)
		}
		doHelp(helpCmd.Args())
		os.Exit(exitCodeOk)
	case "orders":
		if err := ordersCmd.Parse(os.Args[2:]); err != nil {
			logError.Println(err)
			ordersCmd.Usage()
			os.Exit(exitCodeError)
		}
		doOrders(arguments{
			username:        *ordersCmdUsername,
			pathToTokenFile: *ordersCmdPathTokenFile,
			verbose:         *ordersCmdVerbose,
			output:          *ordersCmdOutput,
			format:          *ordersCmdFormat,
		})
		os.Exit(exitCodeOk)
	case "positions":
		if err := positionsCmd.Parse(os.Args[2:]); err != nil {
			logError.Println(err)
			positionsCmd.Usage()
			os.Exit(exitCodeError)
		}
		doPositions(arguments{
			username: *positionsCmdUsername,
			verbose:  *positionsCmdVerbose,
			output:   *positionsCmdOutput,
			all:      *positionsCmdAll,
		})
		os.Exit(exitCodeOk)
	case "options":
		switch flag.Arg(1) {
		case "orders":
			if err := optionsOrdersCmd.Parse(os.Args[3:]); err != nil {
				logError.Println(err)
				optionsOrdersCmd.Usage()
				os.Exit(exitCodeError)
			}
			doOptionsOrders(arguments{
				username:        *optionsOrdersCmdUsername,
				pathToTokenFile: *optionsOrdersCmdPathTokenFile,
				verbose:         *optionsOrdersCmdVerbose,
				output:          *optionsOrdersCmdOutput,
				format:          *optionsOrdersCmdFormat,
			})
			os.Exit(exitCodeOk)
		default:
			flag.Usage()
			os.Exit(exitCodeError)
		}
	default:
		flag.Usage()
		os.Exit(exitCodeError)
	}
}
