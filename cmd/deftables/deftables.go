package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pdk/jcquery/jsonkeys"
)

func main() {

	keys, err := jsonkeys.GetKeys(os.Stdin)
	if err != nil {
		log.Fatalf("%v", err)
	}

	defs, err := computeTableDefs(os.Args[1], keys)
	if err != nil {
		log.Fatalf("%v", err)
	}

	for i, tab := range defs {
		printTableDef(tab)
		if i < len(defs)-1 {
			fmt.Printf("\n")
		}
	}
}

func printTableDef(def tableDef) {
	fmt.Printf("create table if not exists %s (\n    %s\n);\n", def.tableName, strings.Join(def.columnNames, ",\n    "))
}

type tableDef struct {
	tableName   string
	columnNames []string
}

func computeTableDefs(rootTableName string, jsonPaths []string) ([]tableDef, error) {

	// ordered list of table name
	tableNames := []string{rootTableName}

	rootTable := tableDef{
		tableName:   rootTableName,
		columnNames: []string{"id"},
	}

	// map of table name to tableDef
	defs := map[string]tableDef{
		rootTableName: rootTable,
	}

	for _, p := range jsonPaths {

		elems := strings.Split(rootTableName+p, "/")
		tableName := strings.Join(elems[0:len(elems)-1], "_")
		tableName = strings.TrimRight(tableName, "_")
		tableName = strings.ReplaceAll(tableName, "__", "_")
		columnName := elems[len(elems)-1]
		if columnName == "" {
			columnName = "value"
		}

		tab, ok := defs[tableName]
		if !ok {
			tableNames = append(tableNames, tableName)
			tab = tableDef{
				tableName:   tableName,
				columnNames: []string{"id", "fk"},
			}
		}

		tab.columnNames = appendNew(tab.columnNames, columnName)
		defs[tableName] = tab
	}

	result := []tableDef{}

	for _, tn := range tableNames {
		result = append(result, defs[tn])
	}

	return result, nil
}

func appendNew(base []string, next string) []string {

	for _, s := range base {
		if s == next {
			return base
		}
	}

	return append(base, next)
}
