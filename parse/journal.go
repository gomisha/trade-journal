package parse1

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Journal struct {
}

func NewJournal() Journal {
	return Journal{}
}

// ScrubFile removed lines that will break the CSV parser.
// Specifically lines with double quotes in the middle of the column that are not escaped.
func ScrubFile(csvDir string, csvFile string) {
	csvFilePath := filepath.Join(csvDir, csvFile)

	input, err := os.ReadFile(csvFilePath)
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
	err = os.WriteFile(csvFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (j Journal) ParseTrades(csvDir string, csvFile string) [][]string {
	ScrubFile(csvDir, csvFile)

	//filepath.Join(rawJsonFilePath, testData.RawJSONTestRunFile)
	csvFilePath := filepath.Join(csvDir, csvFile)
	file, err := os.Open(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// expect variable number of columns so parser won't crash
	reader.FieldsPerRecord = -1

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", rec)
	}

	//records, err := reader.ReadAll()

	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(records)

	return [][]string{}
}

func (j Journal) toCsv([][]string) []string {
	return []string{}
}
