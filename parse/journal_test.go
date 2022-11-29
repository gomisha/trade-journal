package parse

import (
	"testing"
)
import "github.com/stretchr/testify/require"

type TestData struct {
	expectedTransactions []Transaction
	filePath             string
}

func TestReadTransactions(t *testing.T) {
	expectedTransactions := []Transaction{
		{
			ticker:     "PR",
			account:    "TFSA",
			date:       "2022-11-25",
			commission: "-3",
			price:      "10.588333333",
			shares:     "600",
			buySell:    "Buy",
			action:     "Trade",
		},
		{
			ticker:          "PR",
			account:         "TFSA",
			date:            "2022-11-25",
			commission:      "-3.0190707",
			buySell:         "Sell",
			optionContracts: "-6",
			optionContract:  "PR 20JAN23 9 C",
			price:           "1.971666667",
			action:          "Trade - Option",
		},
		{
			ticker:          "PR",
			account:         "TFSA",
			date:            "2022-11-25",
			commission:      "-0.9789",
			buySell:         "Buy",
			optionContracts: "6",
			optionContract:  "PR 20JAN23 5 P",
			price:           "0.053333333",
			action:          "Trade - Option",
		},
	}

	testDataMap := map[string]TestData{
		"stock, short call, long put": {
			expectedTransactions: expectedTransactions,
			filePath:             "../testdata/input/1-dmc.csv",
		},
	}

	for k, testData := range testDataMap {
		t.Run(k, func(t *testing.T) {
			// read original csv file trade data
			journal := NewJournal()
			actualTransactions := journal.ReadTransactions(testData.filePath)

			require.ElementsMatch(t, testData.expectedTransactions, actualTransactions)
		})
	}
}
