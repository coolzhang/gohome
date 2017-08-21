package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

var user, password, host, port, schema, table, startdate, stopdate string
var partition, shard int

func init() {
	flag.StringVar(&user, "u", "", "db username")
	flag.StringVar(&password, "p", "", "db password")
	flag.StringVar(&host, "h", "", "db IP")
	flag.StringVar(&port, "P", "", "db port")
	flag.StringVar(&schema, "d", "", "db name")
	flag.StringVar(&table, "t", "", "table name")
	flag.IntVar(&partition, "n", 0, "table partition number")
	flag.IntVar(&shard, "m", 0, "db sharding number")
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
	var wg sync.WaitGroup
	ch := make(chan int, 256)
	for j := 0; j < shard; j++ {
		schema := schema + "_" + strconv.Itoa(j)
		wg.Add(1)
		go func() {
			defer wg.Done()
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
	//go func() {
	wg.Wait()
	close(ch)
	//}()

	for c := range ch {
		total += c
	}
	fmt.Println(total)
	return int(total)
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

func testOnce() {
	var once sync.Once
	/*
		onceBody := func() {
			fmt.Println("Only once")
		}
	*/
	done := make(chan int)
	for i := 0; i < 10; i++ {
		go func(i int) {
			once.Do(func() { fmt.Println("Only once") })
			done <- i
		}(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(<-done)
	}
}

func main() {
	/*
		starttime := time.Now()
		fmt.Printf("Start...on %s\n", starttime)
		flag.Parse()
		orderSum()
		fmt.Printf("Cost time: %s\n", time.Now().Sub(starttime).String())
	*/
	testOnce()

	//orderSumSave()
}
