package parse1

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Journal struct {
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

func (j Journal) ParseTrades(csvPath string) [][]string {
	ScrubFile(csvPath)

	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	// expect variable number of columns so parser won't crash
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)

	return [][]string{}
}

func (j Journal) toCsv([][]string) []string {
	return []string{}
}
