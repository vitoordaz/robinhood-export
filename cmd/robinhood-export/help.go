package main

import (
	"fmt"
	"os"
)

func doHelp(cmd string) {
	switch cmd {
	case "dividends":
		dividendsCmd.Usage()
	case "orders":
		ordersCmd.Usage()
	case "positions":
		positionsCmd.Usage()
	}
}

func printUsage() {
	fmt.Printf("%s is a tool for exporting data from Robinhood (https://robinhood.com/).\n\n", os.Args[0])
	fmt.Printf("Usage:\n\n")
	fmt.Printf("	%s <command> [arguments]\n\n", os.Args[0])
	fmt.Println(`The commands are:`)
	fmt.Println()
	fmt.Println(`	dividends	exports all dividends`)
	fmt.Println(`	orders		exports all orders`)
	fmt.Println(`	positions	exports all positions`)
	fmt.Println()
	fmt.Printf("Use '%s help <command>' for more information about a command.\n\n", os.Args[0])
}

func printHelpUsage() {
	fmt.Printf("Usage: %s help <command>\n", os.Args[0])
}

func printDividendsUsage() {
	fmt.Printf("Usage: %s dividends [arguments]\n", os.Args[0])
	ordersCmd.PrintDefaults()
}

func printOrdersUsage() {
	fmt.Printf("Usage: %s orders [arguments]\n", os.Args[0])
	ordersCmd.PrintDefaults()
}

func printPositionsUsage() {
	fmt.Printf("Usage: %s positions [arguments]\n", os.Args[0])
	positionsCmd.PrintDefaults()
}
