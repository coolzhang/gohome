package main

import (
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var groupName, actionName, actionCondition, idType, filterName, methodName, screenName, graphName, graphNameList string
var screenColumns, graphNameLen int

func init() {
	flag.StringVar(&groupName, "g", "", "Host group name")
	flag.StringVar(&actionName, "a", "", "Action name (EventSource: Auto registration)")
	flag.StringVar(&actionCondition, "c", "", "Action condition value (Host name like xxx)")
	flag.StringVar(&methodName, "m", "", "Method name (create, get)")
	flag.StringVar(&idType, "t", "", "ID type name")
	flag.StringVar(&filterName, "n", "", "Source name")
	flag.StringVar(&screenName, "s", "", "Screen name")
	flag.StringVar(&graphName, "gp", "", `Graph name,e.g: "CPU load,CPU utilization"`)
}

func graphNameInfo() (string, int) {
	graphNameList := strings.Split(graphName, ",")
	for i, v := range graphNameList {
		graphNameList[i] = fmt.Sprintf(`"%s"`, v)
	}
	graphNameStr := "[" + strings.Join(graphNameList, ",") + "]"
	graphNameLen := len(graphNameList)
	return graphNameStr, graphNameLen
}

func zabbixHttpPost(jsonData string) *simplejson.Json {
	resp, err := http.Post("http://zabbix.intra.wepiao.com/api_jsonrpc.php", "application/json-rpc", strings.NewReader(jsonData))
	if err != nil {
		log.Fatalf("Post: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	json, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatalf("Unmarshal json: %v", err)
	}
	if errJson, ok := json.CheckGet("error"); ok {
		log.Fatal(errJson)
	}
	resultJson := json.Get("result")
	return resultJson
}

func userLogin(user, password string) string {
	var userLoginJson = `{
		"jsonrpc": "2.0",
		"method": "user.login",
		"params": {
			"user": "%s",
			"password": "%s"
		},
		"id": 1
	}`
	userLoginJson = fmt.Sprintf(userLoginJson, user, password)
	resultJson := zabbixHttpPost(userLoginJson)
	token := resultJson.MustString()
	return token
}

func hostGroupCreate(token, groupName string) string {
	var hostGroupCreateJson = `{
		"jsonrpc": "2.0",
		"method":  "hostgroup.create",
		"params": {
			"name": "%s"
		},
		"auth": "%s",
		"id":      1
	}`
	hostGroupCreateJson = fmt.Sprintf(hostGroupCreateJson, groupName, token)
	resultJson := zabbixHttpPost(hostGroupCreateJson)
	gid := resultJson.Get("groupids").GetIndex(0).MustString()
	return gid
}

func hostGroupExists(token, groupid string) bool {
	var hostGroupExistsJson = `{
			"jsonrpc": "2.0",
			"method": "hostgroup.exists",
			"params": {
				"groupid": "%s"
	    		},
	    		"auth": "%s",
	    		"id": 1
	}`
	if groupid == "" {
		return false
	}
	hostGroupExistsJson = fmt.Sprintf(hostGroupExistsJson, groupid, token)
	resultJson := zabbixHttpPost(hostGroupExistsJson)
	ok := resultJson.MustBool()
	return ok

}

func actionAutoRegCreate(token, groupid, actionName, actionCondition string) string {
	var actionAutoRegCreateJson = `{
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
	actionAutoRegCreateJson = fmt.Sprintf(actionAutoRegCreateJson, actionName, actionCondition, groupid, token)
	resultJson := zabbixHttpPost(actionAutoRegCreateJson)
	aid := resultJson.Get("actionids").GetIndex(0).MustInt()
	return strconv.Itoa(aid)
}

func actionAutoRegExists(token, actionid string) bool {
	var actionAutoRegExistsJson = `{
			"jsonrpc": "2.0",
			"method": "action.exists",
			"params": {
				"actionid": "%s"
	    		},
	    		"auth": "%s",
	    		"id": 1
	}`
	if actionid == "" {
		return false
	}
	actionAutoRegExistsJson = fmt.Sprintf(actionAutoRegExistsJson, actionid, token)
	resultJson := zabbixHttpPost(actionAutoRegExistsJson)
	ok := resultJson.MustBool()
	return ok

}

func actionTriggerCreate(token, groupid string) string {
	var actionTriggerCreateJson = `{
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
	actionTriggerCreateJson = fmt.Sprintf(actionTriggerCreateJson, groupid, token)
	resultJson := zabbixHttpPost(actionTriggerCreateJson)
	aid := resultJson.Get("actionids").GetIndex(0).MustInt()
	return strconv.Itoa(aid)
}

func createForMonitor(token string) {
	groupid := hostGroupCreate(token, groupName)
	autoRegActionid := actionAutoRegCreate(token, groupid, actionName, actionCondition)
	// 无法将新组添加到报警邮件组中
	//triggerActionid := actionTriggerCreate(token, groupid)

	if ok := hostGroupExists(token, groupid); ok {
		fmt.Printf("Host group: %s created successfully\n", groupName)
	} else {
		fmt.Printf("Host group: %s created unsuccessfully\n", groupName)
	}

	if ok := actionAutoRegExists(token, autoRegActionid); ok {
		fmt.Printf("Action-AutoRegistration: %s created successfully\n", actionName)
	} else {
		fmt.Printf("Action-AutoRegistration: %s created unsuccessfully\n", actionName)
	}

}

func idGet(token string) string {
	var getJson = `{
    		"jsonrpc": "2.0",
    		"method": "%s.get",
    		"params": {
        		"output": "extend",
        	"filter": {
                	"name": ["%s"],
			"selectHosts": ["hostid"]
        		}
    		},
    		"auth": "%s",
    		"id": 1
	}`
	idName := idType + "id"
	getJson = fmt.Sprintf(getJson, idType, filterName, token)
	resultJson := zabbixHttpPost(getJson)
	var id string
	for i, v := range resultJson.MustArray() {
		if m, ok := v.(map[string]interface{}); ok {
			fmt.Printf("#%d: %s: %s name: %s\n", i, idName, m[idName], m["name"])
			//fmt.Printf("%T\n",m[idName])
			id = m[idName].(string)
		} else {
			log.Fatalln("result type is not map")
		}
	}
	return id
}

func hostidsGet(token string) []string {
	hostgroupGetJson := `{
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

	hostgroupGetJson = fmt.Sprintf(hostgroupGetJson, groupName, token)
	resultJson := zabbixHttpPost(hostgroupGetJson)
	for _, v := range resultJson.MustArray() {
		if m, ok := v.(map[string]interface{}); ok {
			var hosts = m["hosts"].([]interface{})
			var hostids = make([]string, 0)
			for _, hostid := range hosts {
				if hm, ok := hostid.(map[string]interface{}); ok {
					hostids = append(hostids, hm["hostid"].(string))
				} else {
					log.Fatalln("hostid type is not map")
				}
			}
			return hostids
		} else {
			log.Fatalln("result type  is not map")
		}
	}
	return nil
}

func graphGet(token string) map[string][]string {
	var graphGetJson = `{
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
	hostids := hostidsGet(token)
	for _, hostid := range hostids {
		graphJson := fmt.Sprintf(graphGetJson, hostid, graphNameList, token)
		resultJson := zabbixHttpPost(graphJson)
		graphidList := make([]string, 0)
		for _, v := range resultJson.MustArray() {
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

func screenCreate(token string) string {
	screenCreateJson := `{
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

	hostids := hostidsGet(token)
	screenRows := len(hostids)
	screenColumns := graphNameLen
	screenCreateJson = fmt.Sprintf(screenCreateJson, screenName, screenColumns, screenRows, token)
	resultJson := zabbixHttpPost(screenCreateJson)
	screenid := resultJson.Get("screenids").GetIndex(0).MustString()
	return screenid
}

func screenItemCreate(token string) {
	var screenItemCreateJson = `{
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

	screenid := screenCreate(token)
	hostidGraphidMap := graphGet(token)
	y := 0
	screenitemidList := make([]string, 0)
	for _, graphidList := range hostidGraphidMap {
		for x, graphid := range graphidList {
			itemJson := fmt.Sprintf(screenItemCreateJson, screenid, graphid, x, y, token)
			resultJson := zabbixHttpPost(itemJson)
			screenitemid := resultJson.Get("screenitemids").GetIndex(0).MustString()
			screenitemidList = append(screenitemidList, screenitemid)
		}
		y++
	}
	fmt.Printf("screenitemids: %v\n", screenitemidList)
}

func main() {
	flag.Parse()
	graphNameList, graphNameLen = graphNameInfo()

	user, password := "zabbixapi", "zabbixapi"
	token := userLogin(user, password)

	if strings.ToLower(methodName) == "create" && strings.ToLower(idType) == "hostgroup" {
		createForMonitor(token)
	} else if strings.ToLower(methodName) == "create" && strings.ToLower(idType) == "screen" {
		screenItemCreate(token)
	} else if strings.ToLower(methodName) == "get" {
		idGet(token)
	} else {
		fmt.Println("wrong option parameter")
	}
}
