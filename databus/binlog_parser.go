package main

import (
	"os"

	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"golang.org/x/net/context"
)

func main() {
	//binlogFile := "mysql-bin.000004"
	//binlogPos := 4

	cfg := replication.BinlogSyncerConfig{
		ServerID: 999,
		Flavor:   "mysql",
		Host:     "10.1.1.81",
		Port:     3306,
		User:     "admin",
		Password: "opencmug",
	}
	syncer := replication.NewBinlogSyncer(&cfg)
	streamer, _ := syncer.StartSync(mysql.Position{Name: "mysql-bin.000024", Pos: 4})
	for {
		ev, _ := streamer.GetEvent(context.Background())
		ev.Dump(os.Stdout)
	}
}
