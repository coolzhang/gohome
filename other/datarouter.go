package main

import (
	"bufio"
	"flag"
	"fmt"
	"hash/fnv"
	"log"
	"os"
)

var fileName string
var tableCount int

func init() {
	flag.StringVar(&fileName, "f", "", "a file includes open_ids")
	flag.IntVar(&tableCount, "c", 0, "partitioning table count")
}

func hashCode(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {
	flag.Parse()

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	tableidxOpenids := make(map[uint32][]string)
	for scanner.Scan() {
		hash := hashCode(scanner.Text())
		tableIndex := hash % uint32(tableCount)
		tableidxOpenids[tableIndex] = append(tableidxOpenids[tableIndex], scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < tableCount; i++ {
		fmt.Printf("TableIndex: %d Count: %d\n", i, len(tableidxOpenids[uint32(i)]))
	}
}
