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
			date:                 "2022-11-25",
			account:              "TFSA",
			action:               "Trade",
			ticker:               "PR",
			buySell:              "Buy",
			shares:               "600",
			price:                "10.588333333",
			proceeds:             "-6353.00",
			costBasisBuyOrOption: "-6355.999999799999",
			costBasisTotal:       "-6355.999999799999",
			commission:           "-3",
		},
		{
			date:                 "2022-11-25",
			account:              "TFSA",
			action:               "Trade - Option",
			ticker:               "PR",
			optionContract:       "20JAN23 9 C",
			buySell:              "Sell",
			optionContracts:      "-6",
			price:                "1.971666667",
			proceeds:             "1183.00",
			costBasisShare:       "0",
			costBasisBuyOrOption: "1179.9809295",
			commission:           "-3.0190707",
		},
		{
			date:                 "2022-11-25",
			account:              "TFSA",
			action:               "Trade - Option",
			ticker:               "PR",
			optionContract:       "20JAN23 5 P",
			buySell:              "Buy",
			optionContracts:      "6",
			price:                "0.053333333",
			proceeds:             "-32.00",
			costBasisShare:       "0",
			costBasisBuyOrOption: "-32.9788998",
			commission:           "-0.9789",
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

	expectedTransactions3 := []Transaction{
		{
			date:         "2023-06-05",
			account:      "Margin",
			action:       "Forex",
			commission:   "-2",
			forexUSDBuy:  "4,838.82",
			forexUSDCAD:  "1.3433",
			forexCADSell: "-6499.986906",
			notes:        "converted all CAD to USD",
		},
		{
			date:                 "2023-06-05",
			account:              "Margin",
			action:               "Trade",
			ticker:               "TECK",
			buySell:              "Buy",
			shares:               "100",
			price:                "42.09",
			proceeds:             "-4209.00",
			costBasisBuyOrOption: "-4209.370257250001",
			costBasisTotal:       "-4209.370257250001",
			commission:           "-0.37025725",
		},
		{
			date:                 "2023-06-05",
			account:              "Margin",
			action:               "Trade - Option",
			ticker:               "TECK",
			optionContract:       "21JUL23 38 C",
			buySell:              "Sell",
			optionContracts:      "-1",
			price:                "5.07",
			proceeds:             "507.00",
			costBasisShare:       "0",
			costBasisBuyOrOption: "505.944454",
			commission:           "-1.055546",
		},
	}

	expectedTransactions4 := []Transaction{
		{
			date:     "2023-06-09",
			account:  "Margin",
			action:   "Dividend",
			ticker:   "SMG",
			dividend: "66",
			fee:      "-9.9",
			notes:    "SMG(US8101861065) Payment in Lieu of Dividend (Ordinary Dividend)\n15% tax withdrawn",
		},
	}

	expectedTransactions5 := []Transaction{
		{
			date:           "2023-06-08",
			account:        "RRSP",
			action:         "Trade - Option - Assignment",
			ticker:         "FDX",
			optionContract: "16JUN23 155 C",
			buySell:        "Sell",
			shares:         "-100",
			price:          "155",
			proceeds:       "15500.00",
			costBasisShare: "-172.67370257",
			costBasisTotal: "17267.370257", // import IBKR value and multiply by -1
			realizedPL:     "3873.744617",  // imports IBKR value
			commission:     "-0.1385",
		},
	}

	expectedTransactions6 := []Transaction{
		{
			date:                 "2023-06-08",
			account:              "TFSA",
			action:               "Trade - Close",
			ticker:               "BBWI",
			buySell:              "Sell",
			shares:               "-100",
			price:                "41.44",
			proceeds:             "4144",
			costBasisShare:       "-38.17",
			costBasisBuyOrOption: "",
			costBasisTotal:       "-3,817",     // imports IBKR value
			realizedPL:           "326.482091", // imports IBKR value
			commission:           "-0.51790925",
		},
		{
			date:                 "2023-06-08",
			account:              "TFSA",
			action:               "Trade - Option",
			ticker:               "BBWI",
			optionContract:       "16JUN23 35 C",
			buySell:              "Buy",
			optionContracts:      "1",
			price:                "6.53",
			proceeds:             "-653",
			costBasisShare:       "0",
			costBasisBuyOrOption: "-654.05155",
			commission:           "-1.05155",
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
		"forex": {
			expectedTransactions: expectedTransactions3,
			filePath:             "../testdata/input/3-forex.csv",
		},
		"dividend - withholding tax": {
			expectedTransactions: expectedTransactions4,
			filePath:             "../testdata/input/4-dividend-withholding-tax.csv",
		},
		"call assignment": {
			expectedTransactions: expectedTransactions5,
			filePath:             "../testdata/input/5-call-assignment.csv",
		},
		"hit target": {
			expectedTransactions: expectedTransactions6,
			filePath:             "../testdata/input/6-hit-target.csv",
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
