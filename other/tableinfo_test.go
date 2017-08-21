package main

import (
	"database/sql"
	"testing"
)

func init() {
	partition = 16
	table = "order_info"
	schema = "ordercenter"
	password = "opencmug"
	user = "admin"
	port = "3306"
	host = "10.1.1.85"
	startdate = "2016-09-16"
	stopdate = "2016-09-17"
	var err error
	db, err = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/")
	if err != nil {
		panic(err)
	}
}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a()
	}
}
