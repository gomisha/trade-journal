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
			stockPrice: "10.588333333",
			shares:     "600",
			action:     "Buy",
		},
		{
			ticker:          "PR",
			account:         "TFSA",
			date:            "2022-11-25",
			commission:      "-3.0190707",
			action:          "Sell",
			optionContracts: "-6",
			optionContract:  "PR 20JAN23 9 C",
			optionPrice:     "1.971666667",
		},
		{
			ticker:          "PR",
			account:         "TFSA",
			date:            "2022-11-25",
			commission:      "-0.9789",
			action:          "Buy",
			optionContracts: "6",
			optionContract:  "PR 20JAN23 5 P",
			optionPrice:     "0.053333333",
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
