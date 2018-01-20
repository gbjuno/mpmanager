package main

import (
	"testing"
)

func Test_ParseExcelNewCompany(t *testing.T) {
	InitializeDB("root", "gf37888676", "127.0.0.1", "3306")
	var errorLines string
	var errorStr string
	var err error
	if errorLines, errorStr, err = ParseExcelNewCompany("/tmp/company.xlsx"); err != nil {
		t.Fatalf("ParseExcelNewCompany run failed\n")
	}
	t.Fatalf("ParseExcelNewCompany run successfully, error lines: %s, errorStr %s\n", errorLines, errorStr)
}
