## zbxapitool.go
`zbxapitool`: a tool that could modify the configuaration of Zabbix via Zabbix API.

### Usage
**Help**:  

    # zbxapitool -h  
    Usage of ./zbxapitool:
      -actionCondition string
            Action condition value (Host name like xxx)
      -actionName string
            Action name (EventSource: Auto registration)
      -configType string
            Configure only supports: hostgroup, screen
      -graphName string
            Graph name,e.g: "CPU load,CPU utilization"
      -groupName string
            Host group name
      -screenName string
            Screen name

**Create hostgroup and action**:  

    # zbxapitool -groupName testAPI -actionName testAPI -actionCondition testAPI -configType hostgroup

**Create screen**:  

    # zbxapitool -screenName testScreen -graphName "CPU load,CPU utilization" -groupName testAPI -configType screen
