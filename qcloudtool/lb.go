package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/QcloudApi/qcloud_sign_golang"
)

func main() {
	// 公共参数
	secretId := ""
	secretKey := ""
	// 连接配置
	config := map[string]interface{}{"secretId": secretId, "secretKey": secretKey}
	// 请求参数
	params := map[string]interface{}{"Region": "gz", "Action": "DescribeLoadBalancers", "loadBalancerType": 3, "loadBalancerVips.n": "10.3.105.146"}

	// 发送请求
	retData, err := QcloudApi.SendRequest("lb", params, config)
	if err != nil {
		fmt.Print("Error: ", err)
		return
	}

	var jsonObj interface{}
	err = json.Unmarshal([]byte(retData), &jsonObj)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonOut, _ := json.MarshalIndent(jsonObj, "", "  ")
	b2 := append(jsonOut, '\n')
	os.Stdout.Write(b2)
	return

}
