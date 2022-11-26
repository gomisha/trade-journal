package parse

import "testing"
import "github.com/stretchr/testify/require"

func TestConvert(t *testing.T) {
	// read csv file
	// read trade
	//

	expected := [][]string{
		{"2022-11-25", "TFSA", "", "Trade", "", "", "PR", "", "", "", "Buy", "", "600", "10.588333333", "", "", "", "", "", "", "", "-3"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 9 C", "Sell", "-6", "", "1.971666667", "", "", "", "", "", "", "", "-3.0190707"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 5 P", "Buy", "6", "", "0.053333333", "", "", "", "", "", "", "", "-0.9789"},
	}

	actual := [][]string{
		{"2022-11-25", "TFSA", "", "Trade", "", "", "PR", "", "", "", "Buy", "", "600", "10.588333333", "", "", "", "", "", "", "", "-3"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 9 C", "Sell", "-6", "", "1.971666667", "", "", "", "", "", "", "", "-3.0190707"},
		{"2022-11-25", "TFSA", "", "Trade - Option", "", "", "PR", "", "", "PR 20JAN23 5 P", "Buy", "6", "", "0.053333333", "", "", "", "", "", "", "", "-0.9789"},
	}
	require.Equal(t, expected, actual)
}
