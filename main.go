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
    	log.Fatalln("ERROR:Failed to open log file:", err)
	}
	log.SetOutput(file)

	//Получение сетевых интерфейсов
	devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatalln("ERROR:Failed to get devices:", err)
    }
	
}