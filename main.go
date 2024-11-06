package main

import (
	"os"
	"fmt"
	"log"
	"github.com/google/gopacket/pcap"
)

func main() {

	//Инициализация логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
    	log.Fatal("ERROR:Failed to open log file:", err)
	}
	log.SetOutput(file)


	
	devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Println(err)
    }
	for _, device := range devices {
        fmt.Println(device.Name)
        for _, address := range device.Addresses {
            fmt.Printf(" IP: %s\n", address.IP)
			fmt.Printf(" Netmask: %s\n", address.Netmask)
        }
    }
}