## zbxapitool.go
`zbxapitool`: a tool that could modify the configuaration of Zabbix via Zabbix API.

### Usage
**Help**:  

    # zbxapitool -h  

**Create hostgroup and action**:  

    # zbxapitool -m create -t hostgroup -g ordercenter-mysql -a ordercenter-mysql -c ordercenter

**Create screen**:  

    # zbxapitool -m create -t screen -g ordercenter-mysql -s ordercenter-mysql -gp "MySQL Connections,MySQL Queries executed"
