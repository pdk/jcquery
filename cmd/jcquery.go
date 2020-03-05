package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pdk/jcquery/typer"
)

func main() {

	tableDef := flag.Bool("tabledef", false, "dump sql create table statement")
	csvNoHeaders := flag.Bool("noheaders", false, "csv files have no headers")
	dumpTypes := flag.Bool("types", false, "dump type info")
	queryFilename := flag.String("query", "", "filename of query")

	flag.Parse()

	if *tableDef {
		for _, inputFilename := range flag.Args() {

			csvFile, err := os.Open(inputFilename)
			if err != nil {
				log.Fatalf("%v", err)
			}

			csvDumpCreateTable(csvFile, *csvNoHeaders)
		}

		return
	}

	if *dumpTypes {
		for _, inputFilename := range flag.Args() {
			fmt.Printf("column/type counts for %s:\n", inputFilename)

			csvFile, err := os.Open(inputFilename)
			if err != nil {
				log.Fatalf("%v", err)
			}

			csvDumpTypes(csvFile, *csvNoHeaders)
		}

		return
	}

	if *queryFilename == "" {
		log.Fatalf("-query is required")
	}

	// query := os.Args[1]
	csvFilename := os.Args[2]

	csvFile, err := os.Open(csvFilename)
	if err != nil {
		log.Fatalf("%v", err)
	}

	db, _ := sql.Open("sqlite3", "foobar.db")
	db.Exec("create table x (id,email)")

	insertCSV(db, "blah", csvFile)
}

func csvDumpCreateTable(csvFile *os.File, csvNoHeaders bool) {

	r := csv.NewReader(bufio.NewReader(csvFile))

	columnNames, stats := csvGatherStats(r, csvNoHeaders)

	sb := &strings.Builder{}

	fmt.Fprintf(sb, "create table % (\n")

	for i, colName := range columnNames {

		dbType := pickType(stats[colName])
		fmt.Fprintf(sb, "    %-30s %s", colName, dbType)
		if i < len(columnNames)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	fmt.Fprintf(sb, ");\n")

	fmt.Printf("%s", sb.String())
}

func pickType(stat map[string]int) string {

	if len(stat) == 1 {
		for t := range stat {
			return t
		}
	}

	return "text"
}

func csvDumpTypes(csvFile *os.File, csvNoHeaders bool) {

	r := csv.NewReader(bufio.NewReader(csvFile))

	columnNames, stats := csvGatherStats(r, csvNoHeaders)

	for _, colName := range columnNames {
		for typeName, count := range stats[colName] {
			fmt.Printf("%-20s %-15s %5d\n", colName, typeName, count)
			colName = ""
		}
	}

}

func csvGatherStats(r *csv.Reader, csvNoHeaders bool) (columnNames []string, stats map[string]map[string]int) {

	stats = map[string]map[string]int{}

	lineNo := int64(0)
	for {
		lineNo++

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if lineNo == 1 {
			if !csvNoHeaders {
				columnNames = record
				continue
			}

			for i := range record {
				columnNames = append(columnNames, fmt.Sprintf("col%03d", i+1))
			}
		}

		for i, col := range columnNames {
			var gt typer.GuessedValue
			if i > len(record) {
				gt = typer.NullValue
			} else {
				gt = typer.GuessType(record[i])
			}

			if stats[col] == nil {
				stats[col] = map[string]int{}
			}
			typeLabel := gt.GuessedType.String()
			if typeLabel == "Timestamp" {
				typeLabel = gt.TimestampFormatName
			}
			stats[col][typeLabel]++
		}
	}

	return columnNames, stats
}

func createTableStatement(tableName string, columnNames []string, stats map[string]map[string]int) string {

	sb := &strings.Builder{}
	fmt.Fprintf(sb, "create table %s (\n", tableName)

	for i, col := range columnNames {
		typeString := "text"
		fmt.Fprintf(sb, "    %-30s %s", col, typeString)
		if i < len(columnNames) {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	fmt.Fprintf(sb, ")\n")

	return sb.String()
}

func insertCSV(db *sql.DB, tableName string, csvFile *os.File) {

	r := csv.NewReader(bufio.NewReader(csvFile))

	lineNo := int64(0)
	var columnNames []string
	for {
		lineNo++

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if lineNo == 1 {
			columnNames = record
			createTable(db, tableName, columnNames)
		} else {
			insertRow(db, tableName, columnNames, record)
		}
	}
}

func createTable(db *sql.DB, tableName string, columnNames []string) {

	log.Printf("creating table %s: %v", tableName, columnNames)

	_, err := db.Exec("create table " + tableName + "(" + strings.Join(columnNames, ",") + ")")
	if err != nil {
		log.Fatal(err)
	}
}

func insertRow(db *sql.DB, tableName string, columnNames, values []string) {

	log.Printf("inserting row: %v", values)

	bindValues := make([]interface{}, len(values), len(values))
	for i := 0; i < len(values); i++ {
		bindValues[i] = values[i]
	}

	_, err := db.Exec("insert into "+tableName+"("+strings.Join(columnNames, ",")+
		") values ("+bindMarks(len(columnNames))+")", bindValues...)
	if err != nil {
		log.Fatal(err)
	}
}

func bindMarks(n int) string {
	return "?" + strings.Repeat(",?", n-1)
}
