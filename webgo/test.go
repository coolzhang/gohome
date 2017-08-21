package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Now().Local().Format("15"))
	fmt.Println(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
}
