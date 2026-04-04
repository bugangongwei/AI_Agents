package tools

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

type IPRecord struct {
	StartIP uint32
	EndIP   uint32
	City    string
}

var ipRecords []IPRecord

// InitLocationTool loads the IP2Location CSV data into memory.
func InitLocationTool(csvPath string) error {
	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open IP database: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV record: %v", err)
			continue
		}

		if len(record) < 6 {
			continue
		}

		startIP, _ := strconv.ParseUint(record[0], 10, 32)
		endIP, _ := strconv.ParseUint(record[1], 10, 32)
		city := record[5]

		ipRecords = append(ipRecords, IPRecord{
			StartIP: uint32(startIP),
			EndIP:   uint32(endIP),
			City:    city,
		})
	}

	log.Printf("Loaded %d IP records", len(ipRecords))
	return nil
}

// GetCityFromIP resolves an IP string to a city name.
func GetCityFromIP(ipStr string) string {
	ipNum := ipToUint32(ipStr)
	if ipNum == 0 {
		return ""
	}

	// Binary search for the IP range
	low, high := 0, len(ipRecords)-1
	for low <= high {
		mid := low + (high-low)/2
		if ipNum >= ipRecords[mid].StartIP && ipNum <= ipRecords[mid].EndIP {
			return ipRecords[mid].City
		}
		if ipNum < ipRecords[mid].StartIP {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return ""
}

func ipToUint32(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
