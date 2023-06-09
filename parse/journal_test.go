package parse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestData struct {
	expectedTransactions []Transaction
	filePath             string
}

func TestReadTransactions(t *testing.T) {
	expectedTransactions1 := []Transaction{
		{
			date:       "2022-11-25",
			account:    "TFSA",
			action:     "Trade",
			ticker:     "PR",
			buySell:    "Buy",
			shares:     "600",
			price:      "10.588333333",
			commission: "-3",
		},
		{
			date:            "2022-11-25",
			account:         "TFSA",
			action:          "Trade - Option",
			ticker:          "PR",
			optionContract:  "20JAN23 9 C",
			buySell:         "Sell",
			optionContracts: "-6",
			price:           "1.971666667",
			commission:      "-3.0190707",
		},
		{
			date:            "2022-11-25",
			account:         "TFSA",
			action:          "Trade - Option",
			ticker:          "PR",
			optionContract:  "20JAN23 5 P",
			buySell:         "Buy",
			optionContracts: "6",
			price:           "0.053333333",
			commission:      "-0.9789",
		},
	}

	expectedTransactions2 := []Transaction{
		{
			date:     "2023-06-08",
			account:  "RRSP",
			action:   "Dividend",
			ticker:   "MSFT",
			dividend: "136",
			notes:    "MSFT(US5949181045) Cash Dividend USD 0.68 per Share (Ordinary Dividend)",
		},
	}

	testDataMap := map[string]TestData{
		"stock, short call, long put": {
			expectedTransactions: expectedTransactions1,
			filePath:             "../testdata/input/1-dmc.csv",
		},
		"dividend": {
			expectedTransactions: expectedTransactions2,
			filePath:             "../testdata/input/2-dividend.csv",
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
