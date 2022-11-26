package parse1

type Journal struct {
}

func NewJournal() Journal {
	return Journal{}
}

func (j Journal) ParseTrades(csvFilePath string) [][]string {
	return [][]string{}
}

func (j Journal) toCsv([][]string) []string {
	return []string{}
}
