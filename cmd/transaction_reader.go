package main

import (
	"flag"
	"fmt"
	parse1 "github.com/gomisha/trade-journal/parse"
	"os"
)

// usage: go run cmd/transaction_reader.go --data "./testdata/input/1-dmc.csv"
func main() {
	dataFlag := flag.String("data", "", "Path to CSV data.")

	flag.Parse()

	if *dataFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Println("Hello world")

	journal := parse1.NewJournal()
	transactions := journal.ReadTransactions(*dataFlag)
	//csvTransactions := journal.ToCsv(transactions)

	for i, transaction := range transactions {
		fmt.Println("transaction: ", i, " ", transaction)
	}
}
