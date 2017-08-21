package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	_ "time"
)

var user, password, host, port, schema, table, startdate, stopdate string
var partition int

func init() {
	flag.StringVar(&user, "u", "", "db username")
	flag.StringVar(&password, "p", "", "db password")
	flag.StringVar(&host, "h", "", "db IP")
	flag.StringVar(&port, "P", "", "db port")
	flag.StringVar(&schema, "d", "", "db name")
	flag.StringVar(&table, "t", "", "table name")
	flag.IntVar(&partition, "n", 0, "partition number")
	flag.StringVar(&startdate, "start", "", "start date")
	flag.StringVar(&stopdate, "stop", "", "stop date")
}

var db *sql.DB

func main() {
	//starttime := time.Now()
	flag.Parse()
	a()
}

func a() {
	var tables = make([]string, 0)
	for i := 0; i < partition; i++ {
		for j := 0; j < partition; j++ {
			tables = append(tables, schema+"_"+strconv.Itoa(i)+"."+table+"_"+strconv.Itoa(j))
		}
	}
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var count, total int
	for _, table := range tables {
		query := "select count(*) from " + table + " where create_time >=? and create_time <? and status in(6,20)"
		err = db.QueryRow(query, startdate, stopdate).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		total = total + count
	}
	//fmt.Println(time.Now().Sub(starttime).String())
	fmt.Printf("order tatal count: %d\n", total)
}
