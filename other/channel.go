package main

import (
	"fmt"
	"time"
)

func main() {
	message := make(chan int)
	count := 3

	go func() {
		for i := 1; i <= count; i++ {
			fmt.Println("send message #", i)
			message <- i
		}
		//close(message)
	}()
	close(message)
	time.Sleep(1 * time.Second)
	for c := range message {
		fmt.Println(c)
	}

	/*
		time.Sleep(1 * time.Second)
			for i := 1; i <= count; i++ {
				fmt.Println(<-message)
			}
			/*
				done := make(chan bool)
				go func() {
					fmt.Println("goroutine message")
					done <- true
				}()
				<-done
				fmt.Println("main func messgae")
				//<-done
	*/
}
