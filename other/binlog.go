package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	type EventHeader struct {
		Timestamp   uint32
		TypeCode    byte
		ServerID    uint32
		EventLength uint32
		NextPos     uint32
		Flags       uint16
		//ExtHeader   uint16
	}

	f, err := os.Open("mysql-bin.002236")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data := make([]byte, 500)
	if _, err = f.Read(data); err != nil {
		log.Fatal(err)
	}
	h := new(EventHeader)
	/*
		d1 := data[0:4]
		d2 := data[4:5]
		d3 := data[5:9]
		d4 := data[9:13]
		d5 := data[13:17]
		d6 := data[17:19]
		h.Timestamp = binary.LittleEndian.Uint32(d1)
		h.TypeCode = d2
		h.ServerID = binary.LittleEndian.Uint32(d3)
		h.EventLength = binary.LittleEndian.Uint32(d4)
		h.NextPos = binary.LittleEndian.Uint32(d5)
		h.Flags = binary.LittleEndian.Uint16(d6)
	*/
	pos := 4

	h.Timestamp = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	h.TypeCode = data[pos]
	pos++

	h.ServerID = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	h.EventLength = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	h.NextPos = binary.LittleEndian.Uint32(data[pos:])
	pos += 4

	h.Flags = binary.LittleEndian.Uint16(data[pos:])
	pos += 2
	fmt.Printf("Timestamp: %v\n", time.Unix(int64(h.Timestamp), 0))
	fmt.Printf("TypeCode: %d\n", h.TypeCode)
	fmt.Printf("ServerID: %d\n", h.ServerID)
	fmt.Printf("EventLength: %d\n", h.EventLength)
	fmt.Printf("NextPos: %d\n", h.NextPos)
	fmt.Printf("Flags: %d\n", h.Flags)

	type FormatDescriptionEvent struct {
		BinlogVersion         uint16
		ServerVersion         []byte
		CreatedTimestamp      uint32
		EventHeaderLength     uint8
		EventTypeHeaderLength []byte
	}

	fde := new(FormatDescriptionEvent)

	fde.BinlogVersion = binary.LittleEndian.Uint16(data[pos:])
	pos += 2
	fde.ServerVersion = make([]byte, 50)
	copy(fde.ServerVersion, data[pos:])
	pos += 50
	fde.CreatedTimestamp = binary.LittleEndian.Uint32(data[pos:])
	pos += 4
	fde.EventHeaderLength = data[pos]
	pos++
	fde.EventTypeHeaderLength = data[pos:]
	fmt.Printf("Binlog Version: %d\n", fde.BinlogVersion)
	fmt.Printf("Server Version: %s\n", string(fde.ServerVersion))
	fmt.Printf("Created Timestamp: %v\n", fde.CreatedTimestamp)

}
