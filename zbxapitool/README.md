## zbxapitool.go
`zbxapitool`: a tool that could modify the configuaration of Zabbix via Zabbix API.

### Usage
**Help**:  

    # zbxapitool -h  

**Create hostgroup and action**:  

    # zbxapitool -groupName testAPI -actionName testAPI -actionCondition testAPI -configType hostgroup

**Create screen**:  

    # zbxapitool -screenName testScreen -graphName "CPU load,CPU utilization" -groupName testAPI -configType screen
