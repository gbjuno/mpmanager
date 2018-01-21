package main

import (
	"testing"
)

func Test_ParseExcel(t *testing.T) {
	InitializeDB("root", "gf37888676", "127.0.0.1", "3306")
	var errorLines string
	var errorStr string
	var err error
	if errorLines, errorStr, err = parseExcel("/tmp/company.xlsx"); err != nil {
		t.Fatalf("ParseExcelNewCompany run failed, err %s\n", err)
	}
	t.Fatalf("ParseExcelNewCompany run successfully, error lines: %s, errorStr %s\n", errorLines, errorStr)
}

func Test_getData(t *testing.T) {
	InitializeDB("root", "gf37888676", "127.0.0.1", "3306")
	var err error
	if err = getData(); err != nil {
		t.Fatalf("getData run failed, err %s\n", err)
	}
	t.Fatalf("getData run successfully\n")
}
