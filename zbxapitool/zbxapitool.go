package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

var groupName, actionName, actionCondition, screenName, graphName, configType string

func init() {
	flag.StringVar(&groupName, "groupName", "", "Host group name")
	flag.StringVar(&actionName, "actionName", "", "Action name (EventSource: Auto registration)")
	flag.StringVar(&actionCondition, "actionCondition", "", "Action condition value (Host name like xxx)")
	flag.StringVar(&screenName, "screenName", "", "Screen name")
	flag.StringVar(&graphName, "graphName", "", `Graph name,e.g: "CPU load,CPU utilization"`)
	flag.StringVar(&configType, "configType", "", "Configure only supports: hostgroup, screen")
}

func zabbixHttpPost(jsonData string) *simplejson.Json {
	resp, err := http.Post("http://zabbix.intra.wepiao.com/api_jsonrpc.php", "application/json-rpc", strings.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	json, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatalf("Unmarshal json: %v", err)
	}
	if errJSON, ok := json.CheckGet("error"); ok {
		log.Fatal(errJSON)
	}
	resultJSON := json.Get("result")
	return resultJSON
}

func userLogin(user, password string) string {
	var userLoginJSON = `{
		"jsonrpc": "2.0",
		"method": "user.login",
		"params": {
			"user": "%s",
			"password": "%s"
		},
		"id": 1
	}`
	userLoginJSON = fmt.Sprintf(userLoginJSON, user, password)
	resultJSON := zabbixHttpPost(userLoginJSON)
	token := resultJSON.MustString()
	return token
}

func hostGroupCreate(token, groupName string) string {
	var gid string
	exists := hostGroupExists(token, groupName)
	fmt.Printf("exists: %v\n", exists)
	if !exists {
		var hostGroupCreateJSON = `{
			"jsonrpc": "2.0",
			"method":  "hostgroup.create",
			"params": {
				"name": "%s"
			},
			"auth": "%s",
			"id": 1
		}`
		hostGroupCreateJSON = fmt.Sprintf(hostGroupCreateJSON, groupName, token)
		resultJSON := zabbixHttpPost(hostGroupCreateJSON)
		gid = resultJSON.Get("groupids").GetIndex(0).MustString()
	}
	return gid
}

func hostGroupExists(token, groupName string) bool {
	var hostGroupExistsJSON = `{
			"jsonrpc": "2.0",
			"method": "hostgroup.get",
			"params": {
				"filter": {
					"name": "%s"
				}
			},
	    	"auth": "%s",
	    	"id": 1
	}`
	hostGroupExistsJSON = fmt.Sprintf(hostGroupExistsJSON, groupName, token)
	resultJSON := zabbixHttpPost(hostGroupExistsJSON)
	for _, v := range resultJSON.MustArray() {
		if res, ok := v.(map[string]interface{}); ok {
			if _, ok := res["groupid"]; ok {
				return true
			}
		} else {
			log.Fatalln("result type is not map!")
		}
	}
	return false
}

func actionAutoRegCreate(token, groupid, actionName, actionCondition string) string {
	var aid string
	exists := actionAutoRegExists(token, actionName)
	fmt.Printf("exists: %v\n", exists)
	if !exists {
		var actionAutoRegCreateJSON = `{
		        "jsonrpc": "2.0",
		        "method": "action.create",
		        "params": {
		             "name": "%s",
					  "status": 0,
		              "esc_period": 0,
		              "eventsource": 2,
		              "filter": {
					       "evaltype": 0,
		                   "conditions": [
		                        {
		                             "conditiontype": 22,
		                             "value": "%s",
		                             "operator": 2
		                        }
		                   ]    
		              },
		              "operations": [
		                   {
		                      "operationtype": 4,
		                      "opgroup": [
		                           {
		                              "groupid": "%s"
		                           }
		                       ]
						   }
		
		              ]
				},
				"auth": "%s",
				"id": 1
		}`
		actionAutoRegCreateJSON = fmt.Sprintf(actionAutoRegCreateJSON, actionName, actionCondition, groupid, token)
		resultJSON := zabbixHttpPost(actionAutoRegCreateJSON)
		res := resultJSON.Get("actionids").GetIndex(0).MustInt()
		fmt.Println("actionid:", res)
		aid = strconv.Itoa(res)
	}
	return aid
}

func actionAutoRegExists(token, actionName string) bool {
	var actionAutoRegExistsJSON = `{
			"jsonrpc": "2.0",
			"method": "action.get",
			"params": {
				"filter": {
					"eventsource": 2,
					"name": "%s"
				}
	    	},
	    	"auth": "%s",
	    	"id": 1
	}`
	actionAutoRegExistsJSON = fmt.Sprintf(actionAutoRegExistsJSON, actionName, token)
	resultJSON := zabbixHttpPost(actionAutoRegExistsJSON)
	for _, v := range resultJSON.MustArray() {
		if res, ok := v.(map[string]interface{}); ok {
			if _, ok := res["actionid"]; ok {
				return true
			}
		} else {
			log.Fatalln("result type is not map!")
		}
	}
	return false
}

func actionTriggerCreate(token, groupid string) string {
	var actionTriggerCreateJSON = `{
	    "jsonrpc": "2.0",
	    "method": "action.create",
	    "params": {
	        "name": "test",
	        "eventsource": 0,
	        "esc_period": 3600,
	        "filter": {
	            "evaltype": 0,
	            "conditions": [
	                {
	                    "conditiontype": 1,
	                    "operator": 0,
	                    "value": "%s"
	                }
	            ]
		    }
		},
	    "auth": "%s",
	    "id": 1
	}`
	actionTriggerCreateJSON = fmt.Sprintf(actionTriggerCreateJSON, groupid, token)
	resultJSON := zabbixHttpPost(actionTriggerCreateJSON)
	aid := resultJSON.Get("actionids").GetIndex(0).MustInt()
	return strconv.Itoa(aid)
}

func hostidsGet(token, groupName string) []string {
	var hostids []string
	var hostgroupGetJSON = `{
                "jsonrpc": "2.0",
                "method":  "hostgroup.get",
                "params": {
                        "output": "extend",
                        "filter": {
                                "name": "%s"
                        },
                        "selectHosts": ["hostid"]
                },
                "auth": "%s",
                "id":      1
    }`
	hostgroupGetJSON = fmt.Sprintf(hostgroupGetJSON, groupName, token)
	resultJSON := zabbixHttpPost(hostgroupGetJSON)
	for _, v := range resultJSON.MustArray() {
		if m, ok := v.(map[string]interface{}); ok {
			var hosts = m["hosts"].([]interface{})
			for _, hostid := range hosts {
				if hm, ok := hostid.(map[string]interface{}); ok {
					hostids = append(hostids, hm["hostid"].(string))
				} else {
					log.Fatalln("hostid type is not map")
				}
			}
		} else {
			log.Fatalln("result type  is not map")
		}
	}
	return hostids
}

func graphGet(token, groupName, graphNames string) map[string][]string {
	var graphGetJSON = `{
	                "jsonrpc": "2.0",
	                "method":  "graph.get",
	                "params": {
	                        "output": ["graphid", "name"],
	                        "filter": {
	                                "hostid": %s,
	                                "name": %s
	                        }
	                },
	                "auth": "%s",
	                "id":      1
	}`
	var graphMap = make(map[string][]string)
	hostids := hostidsGet(token, groupName)
	for _, hostid := range hostids {
		graphJSON := fmt.Sprintf(graphGetJSON, hostid, graphNames, token)
		resultJSON := zabbixHttpPost(graphJSON)
		//graphidList := make([]string, 0)
		var graphidList []string
		for _, v := range resultJSON.MustArray() {
			if m, ok := v.(map[string]interface{}); ok {
				graphidList = append(graphidList, m["graphid"].(string))
			} else {
				log.Fatalln("result type is not map")
			}
			graphMap[hostid] = graphidList
		}
	}
	return graphMap
}

func screenCreate(token, groupName, screenName string, graphNamesCount int) string {
	var screenCreateJSON = `{
		    "jsonrpc": "2.0",
		    "method":  "screen.create",
		    "params": {
		         "name": "%s",
		         "hsize": %d,
		         "vsize": %d
		    },
		    "auth": "%s",
		    "id":      1
	}`
	hostids := hostidsGet(token, groupName)
	screenRows := len(hostids)
	screenColumns := graphNamesCount
	screenCreateJSON = fmt.Sprintf(screenCreateJSON, screenName, screenColumns, screenRows, token)
	resultJSON := zabbixHttpPost(screenCreateJSON)
	screenid := resultJSON.Get("screenids").GetIndex(0).MustString()
	return screenid
}

func screenExists(token, screenName string) bool {
	var screenExistsJSON = `{
			"jsonrpc": "2.0",
			"method": "screen.get",
			"params": {
				"filter": {
					"name": "%s"
				}
			},
	    	"auth": "%s",
	    	"id": 1
	}`
	screenExistsJSON = fmt.Sprintf(screenExistsJSON, screenName, token)
	resultJSON := zabbixHttpPost(screenExistsJSON)
	for _, v := range resultJSON.MustArray() {
		if res, ok := v.(map[string]interface{}); ok {
			if _, ok := res["screenid"]; ok {
				return true
			}
		} else {
			log.Fatalln("result type is not map!")
		}
	}
	return false
}

func screenItemCreate(token, groupName, screenName, graphName string) {
	exists := screenExists(token, screenName)
	if !exists {
		var screenItemCreateJSON = `{
    	            "jsonrpc": "2.0",
    	            "method":  "screenitem.create",
    	            "params": {
    	                    "screenid": %s,
    	                    "resourcetype": 0,
    	                    "resourceid": %s,
    	                    "height": 100,
    	                    "width": 500,
    	                    "x": %d,
    	                    "y": %d
    	            },
    	            "auth": "%s",
    	            "id":      1
    	}`

		graphNameList := strings.Split(graphName, ",")
		for i, v := range graphNameList {
			graphNameList[i] = fmt.Sprintf(`"%s"`, v)
		}
		graphNames := "[" + strings.Join(graphNameList, ",") + "]"
		graphNamesCount := len(graphNameList)

		screenid := screenCreate(token, groupName, screenName, graphNamesCount)
		hostidGraphidMap := graphGet(token, groupName, graphNames)
		y := 0
		var screenitemidList []string
		for _, graphidList := range hostidGraphidMap {
			for x, graphid := range graphidList {
				itemJSON := fmt.Sprintf(screenItemCreateJSON, screenid, graphid, x, y, token)
				resultJSON := zabbixHttpPost(itemJSON)
				screenitemid := resultJSON.Get("screenitemids").GetIndex(0).MustString()
				screenitemidList = append(screenitemidList, screenitemid)
			}
			y++
		}
		//fmt.Printf("screenitemids: %v\n", screenitemidList)
	}
}

func main() {
	flag.Parse()

	user, password := "zabbixapi", "zabbixapi"
	token := userLogin(user, password)
	//fmt.Printf("Password: %s\n", token)
	if flag.NFlag() == 4 {
		if strings.ToLower(configType) == "hostgroup" {
			groupid := hostGroupCreate(token, groupName)
			//fmt.Printf("GroupID: %s\n", groupid)
			actionAutoRegCreate(token, groupid, actionName, actionCondition)
			// 无法将新组添加到报警邮件组中
			//triggerActionid := actionTriggerCreate(token, groupid)
		}
		if strings.ToLower(configType) == "screen" {
			screenItemCreate(token, groupName, screenName, graphName)
		}
	} else {
		fmt.Printf("%s: missing some options\n", os.Args[0])
		fmt.Printf("Try '%s -h' for more information.\n", os.Args[0])
	}
}
