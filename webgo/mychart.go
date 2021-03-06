package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
)

var runawayDb, ordercenterDb *sql.DB
var ordercenterSchema, ordercenterTable, runawaySchema, runawayTable string

type tomlConfig struct {
	Ordercenter dsn
	Runaway     dsn
}

type dsn struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
	Table    string
}

func init() {
	var conf tomlConfig
	var err error
	if _, err = toml.DecodeFile("dsn.conf", &conf); err != nil {
		fmt.Println(err)
	}
	ordercenterSchema = conf.Ordercenter.Db
	ordercenterTable = conf.Ordercenter.Table
	runawaySchema = conf.Runaway.Db
	runawayTable = conf.Runaway.Table
	runawayDsn := conf.Runaway.User + ":" + conf.Runaway.Password + "@tcp(" + conf.Runaway.Host + ":" + conf.Runaway.Port + ")/"
	ordercenterDsn := conf.Ordercenter.User + ":" + conf.Ordercenter.Password + "@tcp(" + conf.Ordercenter.Host + ":" + conf.Ordercenter.Port + ")/"

	runawayDb, err = sql.Open("mysql", runawayDsn)
	if err != nil {
		log.Fatal(err)
	}
	ordercenterDb, err = sql.Open("mysql", ordercenterDsn)
	if err != nil {
		log.Fatal(err)
	}
}

func orderTrend(start string, stop string) map[string]interface{} {
	if stop == "" {
		stop = time.Now().Local().Format("2006-01-02")
	}
	if start == "" {
		start = time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	}
	query := "select * from " + runawaySchema + "." + runawayTable + " where cdate >=? and cdate<=?"
	rows, err := runawayDb.Query(query, start, stop)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var count int
	var date string
	var orderCount []int
	var orderDate []string
	for rows.Next() {
		err := rows.Scan(&count, &date)
		if err != nil {
			log.Fatal(err)
		}
		orderCount = append(orderCount, count)
		orderDate = append(orderDate, date)
	}
	trend := map[string]interface{}{
		"orderCount": orderCount,
		"orderDate":  orderDate,
	}
	return trend
}

func orderLive(start string, stop string) int {
	var tables = make([]string, 0)
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			tables = append(tables, ordercenterSchema+"_"+strconv.Itoa(i)+"."+ordercenterTable+"_"+strconv.Itoa(j))
		}
	}

	var count, total int
	for _, table := range tables {
		query := "select count(*) from " + table + " where create_time >=from_unixtime(?) and create_time <from_unixtime(?) and status in(6,20)"
		err := ordercenterDb.QueryRow(query, start, stop).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}

		total = total + count

	}
	return total
}

func orderTotalConcurrent(start string, stop string) int {
	ch := make(chan int, 256)
	group := &sync.WaitGroup{}
	for i := 0; i < 16; i++ {
		i := i
		group.Add(1)
		go func() {
			defer group.Done()
			var tables = make([]string, 0)
			for j := 0; j < 16; j++ {
				tables = append(tables, ordercenterSchema+"_"+strconv.Itoa(i)+"."+ordercenterTable+"_"+strconv.Itoa(j))
			}
			var count int
			for _, table := range tables {
				query := "select count(*) from " + table + " where create_time >=? and create_time <=? and status in(6,20)"
				err := ordercenterDb.QueryRow(query, start, stop).Scan(&count)
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

	var total int
	for count := range ch {
		total = total + count
	}
	return total
}

func orderTotal(start string, stop string) int {
	var tables []string
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			tables = append(tables, ordercenterSchema+"_"+strconv.Itoa(i)+"."+ordercenterTable+"_"+strconv.Itoa(j))
		}
	}

	var total, count int
	for _, table := range tables {
		query := "select count(*) from " + table + " where create_time >=? and create_time <=? and status in(6,20)"
		err := ordercenterDb.QueryRow(query, start, stop).Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
		total = total + count
	}
	return total
}

func liveInit() map[string]interface{} {
	var tData, yData []map[string]interface{}
	var hData []interface{}
	var i int64
	today := time.Now().Unix()
	yeday := today - 86400
	for i = -59; i <= 0; i++ {
		tStopTime := today + i*60
		tStartTime := today + i*60 - 60
		yStopTime := yeday + i*60
		yStartTime := yeday + i*60 - 60
		tOrderCount := orderLive(fmt.Sprint(tStartTime), fmt.Sprint(tStopTime))
		yOrderCount := orderLive(fmt.Sprint(yStartTime), fmt.Sprint(yStopTime))
		tt := map[string]interface{}{
			"x": tStopTime * 1000,
			"y": tOrderCount,
		}
		yy := map[string]interface{}{
			"x": tStopTime * 1000,
			"y": yOrderCount,
		}
		tData = append(tData, tt)
		yData = append(yData, yy)
	}

	date := time.Now().Local().Format("2006-01-02")
	hour, _ := strconv.Atoi(time.Now().Local().Format("15"))
	for i := 0; i < hour; i++ {
		start := date + " " + strconv.Itoa(i) + ":00:00"
		stop := date + " " + strconv.Itoa(i) + ":59:59"
		count := orderTotal(start, stop)
		htime := strconv.Itoa(i) + ":00"
		hc := []interface{}{htime, count}
		hData = append(hData, hc)
	}
	orderCount := map[string]interface{}{
		"todayOrderByMin":  tData,
		"yedayOrderByMin":  yData,
		"totalOrderByHour": hData,
	}
	//fmt.Println(orderCount["totalOrderByHour"])
	return orderCount
}

func trendHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("trend.html")
	if err != nil {
		log.Fatal(err)
	}
	start := r.FormValue("from")
	stop := r.FormValue("to")
	t.Execute(w, orderTrend(start, stop))
}

func totalGetDataHandler(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("starttime")
	stop := r.FormValue("stoptime")
	//count := orderTotalConcurrent(start, stop)
	count := orderTotal(start, stop)
	fmt.Fprint(w, count)
}

func liveGetDataHandler(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("starttime")
	stop := r.FormValue("stoptime")
	count := orderLive(start, stop)
	fmt.Fprint(w, count)
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	/*
		body, err := ioutil.ReadFile("live.html")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(body))
	*/
	t, err := template.ParseFiles("live.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, liveInit())
}

func main() {
	defer runawayDb.Close()
	defer ordercenterDb.Close()
	http.HandleFunc("/", liveHandler)
	http.HandleFunc("/live", liveGetDataHandler)
	http.HandleFunc("/total", totalGetDataHandler)
	http.HandleFunc("/trend", trendHandler)
	http.ListenAndServe(":8080", nil)
}
