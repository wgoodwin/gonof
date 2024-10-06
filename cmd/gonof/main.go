package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func main() {
	file := flag.String("file", os.Getenv("GONOF_FILE"), "file storing gonof data in csv format")
	day := flag.Int("day", 0, "day to reference when updating, defaults to next unwritten day")
	update := flag.Bool("update", false, "whether or not to update data, will show the last 10 days regardless")
	flag.Parse()

	if *file == "" {
		fmt.Println("ERROR: No file provided from the -file flag or GONOF_FILE environment variable set")
		flag.PrintDefaults()
		os.Exit(2)
	}

	keys := []string{
		"Day",
		"Relapse",
		"Mediation (minutes)",
		"Sleep (hours)",
		"Connection",
		"Purpose Building",
		"Workout",
		"Social Media (minutes)",
		"Number of Urges (max 5)",
		"Urge Intensity (1-5)",
	}

	if _, err := os.Stat(*file); errors.Is(err, os.ErrNotExist) {
		createBaseFile(*file, keys)
	}

	allRecords, err := readAndValidate(*file, keys)
	if err != nil {
		panic(err)
	}

	recordsToShow := allRecords
	if len(recordsToShow) > 10 {
		recordsToShow = recordsToShow[len(recordsToShow)-9:]
	}
	writeRecords(recordsToShow, keys, *update)

	if *update {
		nextDay, err := strconv.Atoi(allRecords[len(allRecords)-1][0])
		if err != nil {
			panic(err)
		}

		updateRecord := createBaseRecordRow(nextDay)
		if *day > 0 {
			if *day <= nextDay {
				dayVal := fmt.Sprintf("%d", *day)
				for _, row := range allRecords {
					if row[0] == dayVal {
						updateRecord = row
						break
					}
				}
			}
		}

		fmt.Println(updateRecord)
		// TODO start update loop
		// TODO Write back the data presented and exit

	}

	/*
		ruleMap := map[string]rules.Rule{
			"Meditation": rules.NewChainRule(
				rules.NewEqualRule(0, -10),
				rules.NewLTRule(15, 3),
				rules.NewGTRule(15, 10),
			),
			"Sleep": rules.NewChainRule(
				rules.NewLTRule(6, -10),
				rules.NewLTRule(7, 5),
				rules.NewLTRule(8, 6.5),
				rules.NewLTRule(9, 10),
				rules.NewGTRule(9, 6),
				rules.NewElseRule(10),
			),
			"Connection":       rules.NewYesNoRule(5, -5),
			"Purpose Building": rules.NewYesNoRule(6, -5),
			"Workout":          rules.NewYesNoRule(8, -5),
			"Social Media": rules.NewChainRule(
				rules.NewLTRule(10, 10),
				rules.NewLTRule(30, 4),
				rules.NewLTRule(60, -7),
				rules.NewLTRule(90, -8),
				rules.NewElseRule(-10),
			),
		}
	*/
}

func createBaseFile(file string, headers []string) {
	fmt.Println("creating gonof file:", file)
	f, err := os.Create(file)
	if err != nil {
		panic(err) // might want to change this later
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write(headers)
}

func readAndValidate(file string, headers []string) ([][]string, error) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	if len(headers) != len(records[0]) {
		return [][]string{}, fmt.Errorf("file [%s] is not in gonof CSV format", file)
	}

	for index, col := range records[0] {
		if col != headers[index] {
			return [][]string{}, fmt.Errorf("file [%s] is not in gonof CSV format", file)
		}
	}
	return records[1:], nil
}

func writeRecords(records [][]string, headers []string, willUpdate bool) {
	if len(records) == 0 && !willUpdate {
		fmt.Println("No Records to Show; use the -update flag to add the first row")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Options.DrawBorder = false
	t.SetRowPainter(table.RowPainter(mapRelapse))

	t.AppendHeader(toRow(headers))
	for _, row := range records {
		t.AppendRow(toRow(row))
	}
	t.Render()
}

func toRow(records []string) table.Row {
	var result table.Row
	for _, val := range records {
		result = append(result, val)
	}
	return result
}

func mapRelapse(row table.Row) text.Colors {
	if rVal, ok := row[1].(string); ok {
		switch strings.ToLower(rVal) {
		case "yes":
			return text.Colors{text.FgRed}
		case "no":
			return text.Colors{text.FgGreen}
		}
	}
	return text.Colors{}
}

func createBaseRecordRow(newDay int) []string {
	return []string{fmt.Sprintf("%d", newDay), "no", "0", "0", "no", "no", "no", "0", "0", "0"}
}
