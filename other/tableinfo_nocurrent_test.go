package main

import "testing"

func init() {
	partition = 16
	table = "order_info"
	schema = "ordercenter"
	password = "test"
	user = "cmug"
	port = "3306"
	host = "10.1.1.62"
}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		a()
	}
}
