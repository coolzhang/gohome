package main

import (
	"database/sql"
	"flag"
	"log"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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

func orderSum() int {
	db, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var tables = make([]string, 0)
	for i := 0; i < partition; i++ {
		tables = append(tables, table+"_"+strconv.Itoa(i))
	}

	var total int
	ch := make(chan int, 256)
	group := &sync.WaitGroup{}
	for j := 0; j < partition; j++ {
		schema := schema + "_" + strconv.Itoa(j)
		group.Add(1)
		go func() {
			defer group.Done()
			for _, table := range tables {
				var count int
				query := "select count(*) from " + schema + "." + table + " where create_time >=? and create_time <? and status in(6,20)"
				err := db.QueryRow(query, startdate, stopdate).Scan(&count)
				if err != nil {
					log.Fatal(err)
				}
				ch <- count
			}
		}()
	}

	go func() {
		group.Wait()
		close(ch)
	}()

	for count := range ch {
		total = total + count
	}
	return total
}

func orderSumSave() {
	sum := orderSum()

	db, err := sql.Open("mysql", "runaway:run@2018@tcp(10.1.1.62:3306)/")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into runaway.order_summary(order_count,cdate) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(sum, startdate)
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()
}

func main() {
	flag.Parse()

	orderSumSave()
}
