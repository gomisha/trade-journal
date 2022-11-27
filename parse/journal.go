package parse1

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Transaction struct {
	ticker          string
	account         string
	date            string
	commission      string
	stockPrice      string
	optionPrice     string
	optionContracts string
	shares          string
}

type Journal struct {
	// each entry is for a ticker and all the transactions associated with that ticker
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

func (j *Journal) ParseTrades(csvPath string) [][]string {
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
		}
		if err != nil {
			log.Fatal(err)
		}

		// find account alias
		if rec[0] == "Account Information" && rec[1] == "Data" && rec[2] == "Account Alias" {
			accountAlias = rec[3]
		}

		// find trade entries
		if rec[0] == "Trades" && rec[1] == "Data" && rec[2] == "Order" {
			fmt.Printf("%+v\n", rec)
			dateTime := strings.Split(rec[6], " ")

			transaction := Transaction{
				account:    accountAlias,
				date:       dateTime[0],
				commission: rec[11],
			}
			switch rec[3] {
			case "Stocks":
				// stock ticker will be in this column
				transaction.ticker = rec[5]
				transaction.optionPrice = ""
				transaction.stockPrice = rec[8]
				transaction.shares = rec[7]
				transaction.optionContracts = ""

			case "Equity and Index Options":
				optionTicker := strings.Split(rec[5], " ")
				// options ticker will be in first split index: PR 20JAN23 9 C
				transaction.ticker = optionTicker[0]
				transaction.stockPrice = ""
				transaction.optionPrice = rec[8]
				transaction.shares = ""
				transaction.optionContracts = rec[7]

			default:
				log.Fatal("Invalid transaction type: ", rec[3])
			}
			j.addTransaction(transaction)
		}
	}
	return [][]string{}
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
	//j.trades[transaction.ticker] = append(j.trades[transaction.ticker], transaction)
}

func (j *Journal) toCsv([][]string) []string {
	return []string{}
}
