package parse

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Transaction struct {
	ticker     string
	account    string
	date       string
	commission string
	price      string // stock / option price

	optionContracts string // # of contracts
	optionContract  string // contract name e.g. PR 20JAN23 9 C
	shares          string
	buySell         string // buy / sell / transfer
	action          string // trade / trade-option / dividend

	proceeds             string // will be calculated, not imported
	costBasisShare       string // will be calculated, not imported
	costBasisBuyOrOption string // will be calculated, not imported
	costBasisTotal       string // will be calculated, not imported
	realizedPL           string // will be calculated, not imported

	forexUSDBuy  string // USD bought during CAD -> USD forex
	forexUSDCAD  string // exchange rate USD/CAD
	forexCADSell string // CAD sold during CAD -> USD forex

	dividend string // dividend payment
	fee      string // e.g. dividend withholding, monthly live data subscription
	notes    string // automated notes (e.g. dividend payment)
}

type Journal struct {
	// each map entry is for a ticker and all the transactions associated with that ticker
	trades map[string][]Transaction
}

func NewJournal() Journal {
	return Journal{}
}

// ScrubFile removed lines that will break the CSV parser.
// Specifically lines with double quotes in the middle of the column that are not escaped.
func ScrubFile(csvPath string) {
	input, err := os.ReadFile(csvPath)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "Quantities preceded by a \"-\" sign ") {
			lines[i] = "REMOVED LINE"
		}
	}
	output := strings.Join(lines, "\n")
	err = os.WriteFile(csvPath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// ReadTransactions reads the raw CSV transactions that are autogenerated by IBKR "Activity Statement" and
// converts them to a list of Transaction structs.
func (j *Journal) ReadTransactions(csvPath string) []Transaction {
	ScrubFile(csvPath)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// expect variable number of columns so parser won't crash
	reader.FieldsPerRecord = -1
	accountAlias := ""

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else

		// find account alias
		if rec[0] == "Account Information" && rec[1] == "Data" && rec[2] == "Account Alias" {
			accountAlias = rec[3]
		} else if rec[0] == "Dividends" && rec[1] == "Data" && rec[2] == "USD" {
			// dividend transaction
			ticker := strings.Split(rec[4], "(") // e.g. MSFT(US5949181045) Cash Dividend USD 0.68 per Share (Ordinary Dividend)

			transaction := Transaction{
				date:     rec[3],
				account:  accountAlias,
				action:   "Dividend",
				ticker:   ticker[0],
				dividend: rec[5],
				notes:    rec[4],
			}
			j.addTransaction(transaction)
		} else if rec[0] == "Withholding Tax" && rec[1] == "Data" && rec[2] == "USD" {
			// some dividend payments will have 15% withholding tax

			// e.g. SMG(US8101861065) Payment in Lieu of Dividend - US Tax
			ticker := strings.Split(rec[4], "(")[0]

			// look up transactions by ticker and ensure there's a single dividend transaction
			transaction := j.findSingleTransaction(ticker, "Dividend")

			transaction.fee = rec[5]
			transaction.notes += "\n15% tax withdrawn"

			j.updateSingleTransaction(ticker, transaction)
		} else if rec[0] == "Trades" && rec[1] == "Data" && rec[2] == "Order" {
			// find trade transactions

			dateTime := strings.Split(rec[6], ", ")

			transaction := Transaction{
				date:       dateTime[0],
				account:    accountAlias,
				commission: rec[11],
			}
			switch rec[3] {
			case "Stocks":
				transaction.price = rec[8]
				// stock ticker will be in this column
				transaction.ticker = rec[5]
				transaction.shares = rec[7]
				transaction.optionContracts = ""
				transaction.action = "Trade"
				if strings.HasPrefix(transaction.shares, "-") {
					transaction.buySell = "Sell"
				} else {
					transaction.buySell = "Buy"
				}

				shares, err := strconv.ParseFloat(transaction.shares, 64)
				if err != nil {
					panic(err)
				}

				price, err := strconv.ParseFloat(transaction.price, 64)
				if err != nil {
					panic(err)
				}

				// proceeds calculation
				proceeds := -1 * shares * price
				transaction.proceeds = fmt.Sprintf("%.2f", proceeds)

				// cost basis buy or option calculation
				commission, err := strconv.ParseFloat(transaction.commission, 64)
				if err != nil {
					panic(err)
				}

				costBasisBuyOrOption := proceeds + commission
				transaction.costBasisBuyOrOption = fmt.Sprint(costBasisBuyOrOption)
				transaction.costBasisTotal = transaction.costBasisBuyOrOption

			case "Equity and Index Options":
				transaction.price = rec[8]
				optionTicker := strings.Split(rec[5], " ")
				// options ticker will be in first split index: PR 20JAN23 9 C
				optionContract := strings.Split(rec[5], optionTicker[0]+" ")
				transaction.ticker = optionTicker[0]
				transaction.optionContracts = rec[7]
				transaction.optionContract = optionContract[1] //extract "20JAN23 9 C" from "PR 20JAN23 9 C"
				transaction.action = "Trade - Option"
				if strings.HasPrefix(transaction.optionContracts, "-") {
					transaction.buySell = "Sell"
				} else {
					transaction.buySell = "Buy"
				}

				// proceeds calculation
				contracts, err := strconv.ParseFloat(transaction.optionContracts, 64)
				if err != nil {
					panic(err)
				}

				price, err := strconv.ParseFloat(transaction.price, 64)
				if err != nil {
					panic(err)
				}

				proceeds := -100 * contracts * price

				transaction.proceeds = fmt.Sprintf("%.2f", proceeds)

				// cost basis per share calculation
				transaction.costBasisShare = "0"

				// cost basis buy or option calculation
				commission, err := strconv.ParseFloat(transaction.commission, 64)
				if err != nil {
					panic(err)
				}

				costBasisBuyOrOption := proceeds + commission
				transaction.costBasisBuyOrOption = fmt.Sprint(costBasisBuyOrOption)

			case "Forex":
				transaction.action = "Forex"
				// Trades,Data,Order,Forex,CAD,USD.CAD,"2023-06-05, 11:17:59","4,838.82",1.3433,,-6499.986906,-2,,,4.259739,
				transaction.forexUSDBuy = rec[7]
				transaction.forexUSDCAD = rec[8]

				if transaction.commission == "0" {
					transaction.commission = ""
				}

				usdBuyNoComma := strings.ReplaceAll(transaction.forexUSDBuy, ",", "") // remove comma from string
				usdBuy, err := strconv.ParseFloat(usdBuyNoComma, 64)
				if err != nil {
					panic(err)
				}

				usdcad, err := strconv.ParseFloat(transaction.forexUSDCAD, 64)
				if err != nil {
					panic(err)
				}

				cadSell := usdBuy * usdcad * -1

				// use 6 decimal places to correspond with IBKR report
				transaction.forexCADSell = fmt.Sprintf("%.6f", cadSell)

				if usdBuy < 5 {
					transaction.notes = "remaining CAD auto converted"
				} else {
					transaction.notes = "converted all CAD to USD"
				}

			default:
				log.Fatal("Invalid transaction type: ", rec[3])
			}
			j.addTransaction(transaction)
		}
	}

	var transactions []Transaction
	for _, tickerTransactions := range j.trades {
		for _, transaction := range tickerTransactions {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

func (j *Journal) addTransaction(transaction Transaction) {
	if j.trades == nil {
		j.trades = make(map[string][]Transaction)
	}
	// get list of transactions for that ticker
	transactions := j.trades[transaction.ticker]
	if transactions == nil {
		transactions = []Transaction{transaction}
	} else {
		transactions = append(transactions, transaction)
	}
	j.trades[transaction.ticker] = transactions
}

func (j *Journal) findSingleTransaction(ticker string, action string) Transaction {
	if j.trades == nil {
		panic("No transactions in journal")
	}
	// get list of transactions for that ticker
	transactions := j.trades[ticker]
	if transactions == nil {
		panic(fmt.Sprintf("no transactions for ticker %s", ticker))
	}
	if len(transactions) > 1 {
		panic(fmt.Sprintf("expected only 1 transaction for ticker %s but have %d", ticker, len(transactions)))
	}
	if (transactions)[0].action != action {
		panic(fmt.Sprintf("expected transaction for ticker %s has unexpected action %s", ticker, transactions[0].action))
	}
	return (transactions)[0]
}

func (j *Journal) updateSingleTransaction(ticker string, transaction Transaction) {
	if j.trades == nil {
		panic("No transactions in journal")
	}
	// get list of transactions for that ticker
	transactions := j.trades[ticker]
	if transactions == nil {
		panic(fmt.Sprintf("no transactions for ticker %s", ticker))
	}
	if len(transactions) > 1 {
		panic(fmt.Sprintf("expected only 1 transaction for ticker %s but have %d", ticker, len(transactions)))
	}
	j.trades[ticker] = []Transaction{transaction}
}

func (j *Journal) ToCsv(txs []Transaction) {
	// convert [] Transaction to [] string, so they can be written to CSV
	var txsStr [][]string
	for _, tx := range txs {
		var row []string

		//txsStr[i] = make([]string, 10)
		row = append(row, tx.date)
		row = append(row, tx.account)
		row = append(row, "")
		row = append(row, tx.action)
		row = append(row, "")
		row = append(row, "")
		row = append(row, tx.ticker)
		row = append(row, "")
		row = append(row, "")
		row = append(row, tx.optionContract)
		row = append(row, tx.buySell)
		row = append(row, tx.optionContracts)
		row = append(row, tx.shares)
		row = append(row, tx.price)
		row = append(row, tx.proceeds)
		row = append(row, "")
		row = append(row, tx.costBasisShare)
		row = append(row, tx.costBasisBuyOrOption)
		row = append(row, tx.costBasisTotal)
		row = append(row, "")
		row = append(row, tx.dividend)
		row = append(row, tx.commission)
		row = append(row, "")
		row = append(row, "")
		row = append(row, tx.fee)
		row = append(row, "")
		row = append(row, tx.forexUSDBuy)
		row = append(row, tx.forexUSDCAD)
		row = append(row, tx.forexCADSell)
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, tx.notes)

		txsStr = append(txsStr, row)
	}

	// create the file
	f, err := os.Create("./transactions.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	writer.WriteAll(txsStr)
}
