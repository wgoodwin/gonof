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
		nextDay := 1

		if len(allRecords) > 0 {
			nextDay, err = strconv.Atoi(allRecords[len(allRecords)-1][0])
			if err != nil {
				panic(err)
			}
		}

		var dayRecord = len(allRecords)
		updateRecord := createBaseRecordRow(nextDay)
		if *day > 0 {
			if *day <= nextDay {
				dayVal := fmt.Sprintf("%d", *day)
				for i, row := range allRecords {
					if row[0] == dayVal {
						dayRecord = i
						updateRecord = row
						break
					}
				}
			}
		}
		runUpdateLoop(keys, updateRecord)

		if dayRecord == len(allRecords) {
			allRecords = append(allRecords, updateRecord)
		} else {
			allRecords[dayRecord] = updateRecord
		}

		flushRecords(*file, keys, allRecords...)
	}

	// TODO Update this and add graphing logic
	/*
		ruleMap := map[string]rules.Rule{
			"Meditation": rules.NewChainRule(
				rules.NewEqualRule(0, -10),
				rules.NewLTRule(16, 3),
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

func flushRecords(file string, headers []string, rows ...[]string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err) // might want to change this later
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write(headers)
	for _, r := range rows {
		writer.Write(r)
	}
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

func printOptions(keys, row []string) {
	for i, k := range keys[1:] {
		fmt.Printf("%d) %s: %s\n", i+1, k, row[i+1])
	}
}

func runUpdateLoop(keys, row []string) []string {
	var arg string
	for {
		fmt.Println()
		printOptions(keys, row)
		fmt.Print("Select an option (or \"done\" to quit): ")
		fmt.Scanln(&arg)
		if arg == "done" {
			break
		}
		keyIndex, err := strconv.Atoi(arg)
		if err != nil || (keyIndex > len(keys)-1 || keyIndex < 1) {
			fmt.Printf("Invalid key [%s], please try again.", arg)
			continue
		}

		fmt.Printf("Enter value for %s or \"done\" to stop: ", keys[keyIndex])
		fmt.Scanln(&arg)
		if arg == "done" {
			break
		}
		row = validateAndUpdate(arg, keyIndex, keys, row)
	}
	return row
}

func validateAndUpdate(newVal string, index int, keys, row []string) []string {
	switch keys[index] {
	case "Relapse":
		if validateYesNo(keys[index], newVal) {
			row[index] = strings.ToLower(newVal)
		}
	case "Mediation (minutes)":
		if f, ok := validateFloat(keys[index], newVal); ok {
			row[index] = fmt.Sprintf("%.02f", f)
		}
	case "Sleep (hours)":
		if f, ok := validateFloat(keys[index], newVal); ok {
			row[index] = fmt.Sprintf("%.02f", f)
		}
	case "Connection":
		if validateYesNo(keys[index], newVal) {
			row[index] = strings.ToLower(newVal)
		}
	case "Purpose Building":
		if validateYesNo(keys[index], newVal) {
			row[index] = strings.ToLower(newVal)
		}
	case "Workout":
		if validateYesNo(keys[index], newVal) {
			row[index] = strings.ToLower(newVal)
		}
	case "Social Media (minutes)":
		if f, ok := validateFloat(keys[index], newVal); ok {
			row[index] = fmt.Sprintf("%.02f", f)
		}
	case "Number of Urges (max 5)":
		if i, ok := validateIntBounds(keys[index], newVal, 0, 5); ok {
			row[index] = fmt.Sprintf("%d", i)
		}
	case "Urge Intensity (1-5)":
		if i, ok := validateIntBounds(keys[index], newVal, 0, 5); ok {
			row[index] = fmt.Sprintf("%d", i)
		}
	}
	return row
}

func validateIntBounds(key, val string, lower, upper int) (int, bool) {
	i, err := strconv.Atoi(val)
	if err != nil {
		fmt.Printf("Value %s must be an integer.\n", key)
		return 0, false
	}

	if i < lower || i > upper {
		fmt.Printf("Value %s out of bounds, must be within [%d,%d].\n", key, lower, upper)
		return 0, false
	}
	return i, true
}

func validateYesNo(key, val string) bool {
	val = strings.ToLower(val)
	if val != "yes" && val != "no" {
		fmt.Printf("Value %s must be either \"yes\" or \"no\".\n", key)
		return false
	}
	return true
}

func validateFloat(key, val string) (float64, bool) {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		fmt.Printf("Value %s must be a decimal number.\n", key)
		return 0, false
	}
	return f, true
}
