package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"math"
	"os"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]

	pcapFile := argsWithoutProg[0]

	var prevTenMinuteT int64 = 0
	var cumSize int64 = 0


	if handle, err := pcap.OpenOffline(pcapFile); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			packetTime := packet.Metadata().Timestamp
			size1 := packet.Metadata().CaptureLength
			size2 := packet.Metadata().Length

			size := math.Max(float64(size1), float64(size2))

			tenMinuteT := packetTime.Unix() / 600

			cumSize += int64(size)

			if tenMinuteT != prevTenMinuteT {
				currentTime := time.Unix(tenMinuteT * 600, 0)
				kbSize := cumSize / 1000
				kbps := (kbSize * 8) / 600
				fmt.Println("Timestamp: ", currentTime, " -> ", cumSize / 1000, " KB", "\t ", kbps, " kbps")
				cumSize = 0
			}

			prevTenMinuteT = tenMinuteT
		}
	}
}
