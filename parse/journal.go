package parse

import (
	"encoding/csv"
	"io"
	"log"
	"os"
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

	forexUSDBuy  string // USD bought during CAD -> USD forex
	forexUSDCAD  string // exchange rate USD/CAD
	forexCADSell string // CAD sold during CAD -> USD forex

	dividend string // dividend payment
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
		} else

		// dividend transaction
		if rec[0] == "Dividends" && rec[1] == "Data" && rec[2] == "USD" {
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
		} else

		// find trade transactions
		if rec[0] == "Trades" && rec[1] == "Data" && rec[2] == "Order" {
			dateTime := strings.Split(rec[6], ", ")

			transaction := Transaction{
				date:       dateTime[0],
				account:    accountAlias,
				price:      rec[8],
				commission: rec[11],
			}
			switch rec[3] {
			case "Stocks":
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

			case "Equity and Index Options":
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
			case "Forex":
				continue
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
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, tx.dividend)
		row = append(row, tx.commission)
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
		row = append(row, "")
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
