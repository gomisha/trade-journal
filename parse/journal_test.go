package parse1

import "testing"
import "github.com/stretchr/testify/require"

func TestConvert(t *testing.T) {
	// read original csv file trade data
	filePath := "../testdata/input/1-dmc.csv"
	journal := NewJournal()
	actualTransactions := journal.ParseTrades(filePath)

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

	require.Equal(t, expectedTransactions, actualTransactions)
}
