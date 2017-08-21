#!/bin/bash
#
#

source /root/.bash_profile

start=$(date --date="1 day ago" +%Y-%m-%d )
stop=$(date +%Y-%m-%d)
go run /data/mygo/src/gohome/orderinfo.go -u admin -p opencmug -h 10.1.1.85 -P 3306 -d ordercenter -t order_info -n 16 -start ${start} -stop ${stop}
