package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
)

func main() {
	js := map[string]interface{}{
		"Action":          "DescribeInstances",
		"Nonce":           11886,
		"Region":          "gz",
		"SecretId":        "AKIDz8krbsJ5yKBZQpn74WFkmLcmugPhESA",
		"SignatureMethod": "HmacSHA256",
		"Timestamp":       1465185768,
		"instanceIds.0":   "ins-09dx96cmug",
		"limit":           20,
		"offset":          0,
		"debug":           0,
	}
	s := url.Values{}
	for k, v := range js {
		s.Set(fmt.Sprintf("%v", k), fmt.Sprintf("%v", v))
	}
	fmt.Println(s.Encode())
	m, err := url.ParseQuery(s.Encode())
	if err != nil {
		log.Fatal(err)
	}
	debug, _ := strconv.ParseBool(fmt.Sprint(js["debug"]))
	if debug {
		fmt.Println(m)
	}

	f, err := os.Open("file.txt")
	if err != nil {
		fmt.Println(err)
	}
	b := make([]byte, 128)
	_, err = f.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
}
