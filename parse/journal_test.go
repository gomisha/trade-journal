package parse1

import "testing"
import "github.com/stretchr/testify/require"

func TestConvert(t *testing.T) {
	// read original csv file trade data
	filePath := "../testdata/input"
	file := "1-dmc.csv"
	journal := NewJournal()
	actualParsedTrades := journal.ParseTrades(filePath, file)

	expectedTrades := [][]string{
		{"2022-11-25", "TFSA", "", "Trade", "", "", "PR", "", "", "", "Buy", "", "600", "10.588333333", "", "", "", "", "", "", "", "-3"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 9 C", "Sell", "-6", "", "1.971666667", "", "", "", "", "", "", "", "-3.0190707"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 5 P", "Buy", "6", "", "0.053333333", "", "", "", "", "", "", "", "-0.9789"},
	}

	require.Equal(t, expectedTrades, actualParsedTrades)

	// convert to CSV
	actualCsv := journal.toCsv(actualParsedTrades)

	expectedCsv := []string{
		"2022-11-25,TFSA,,Trade,,,PR,,,,Buy,,600,10.588333333,,,,,,,,-3",
		"2022-11-25,TFSA,,Trade - Option,,,PR,,,PR 20JAN23 9 C,Sell,-6,,1.971666667,,,,,,,,-3.0190707",
		"2022-11-25,TFSA,,Trade - Option,,,PR,,,PR 20JAN23 5 P,Buy,6,,0.053333333,,,,,,,,-0.9789",
	}
	require.Equal(t, expectedCsv, actualCsv)
}
