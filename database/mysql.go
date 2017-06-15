package database

import (
	"database/sql"
	"log"

	"github.com/tealeg/xlsx"
)

// GenerateExcel generate excel
func GenerateExcel(host string, port string, username string, password string, schema string, table string) {
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/information_schema")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	var rows *sql.Rows

	if len(table) != 0 {
		stmtOut, err := db.Prepare("select TABLE_SCHEMA as '库名', TABLE_NAME as '表名', COLUMN_NAME as '字段名字', DATA_TYPE as '字段类型', COLUMN_COMMENT as '字段注释' from COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()
		rows, err = stmtOut.Query(schema, table)
	} else {
		stmtOut, err := db.Prepare("select TABLE_SCHEMA as '库名', TABLE_NAME as '表名', COLUMN_NAME as '字段名字', DATA_TYPE as '字段类型', COLUMN_COMMENT as '字段注释' from COLUMNS where TABLE_SCHEMA = ?")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()
		rows, err = stmtOut.Query(schema)
	}
	defer rows.Close()

	mappers := make(map[string]*xlsx.Sheet)

	file := xlsx.NewFile()

	if !rows.Next() {
		log.Fatalln("该库中没有字典信息，请补充后重新执行该命令")
	}

	for rows.Next() {
		var tableSchema string
		var tableName string
		var columnName string
		var dataType string
		var columnComment string
		err = rows.Scan(&tableSchema, &tableName, &columnName, &dataType, &columnComment)
		if err != nil {
			panic(err.Error())
		}
		var sheet *xlsx.Sheet
		if s, ok := mappers[tableName]; ok {
			sheet = s
		} else {
			sheet, err = file.AddSheet(tableName)
			if err != nil {
				panic(err.Error())
			}
			sRow := sheet.AddRow()
			cellOne := sRow.AddCell()
			cellOne.Value = "库名"
			cellTwo := sRow.AddCell()
			cellTwo.Value = "表名"
			cellThree := sRow.AddCell()
			cellThree.Value = "字段名字"
			cellFour := sRow.AddCell()
			cellFour.Value = "字段类型"
			cellFive := sRow.AddCell()
			cellFive.Value = "字段注释"
			mappers[tableName] = sheet
		}
		sheetRow := sheet.AddRow()
		cell1 := sheetRow.AddCell()
		cell1.Value = tableSchema
		cell2 := sheetRow.AddCell()
		cell2.Value = tableName
		cell3 := sheetRow.AddCell()
		cell3.Value = columnName
		cell4 := sheetRow.AddCell()
		cell4.Value = dataType
		cell5 := sheetRow.AddCell()
		cell5.Value = columnComment
	}
	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}
	err = file.Save(schema + ".xlsx")
	if err != nil {
		panic(err.Error())
	}
}
