package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"github.com/google/gopacket/pcap"
	"github.com/gookit/config/v2"
)

type Config struct {
	ComputerName string `json:"computerName"`
	Snaplen int  `json:"snaplen"`
	Promisc bool `json:"promisc"`
	Timeout int `json:"timeout"`
}

func sniffer(iface pcap.Interface, wg *sync.WaitGroup, cfg *Config) {
	defer wg.Done()
}

func main() {

	//Инициализация логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
    	log.Fatalln("ERROR:Ошибка при открытии файла:", err)
	}
	log.SetOutput(file)

	//Получение сетевых интерфейсов
	devices, err := pcap.FindAllDevs()
    if err != nil {
        log.Fatalln("ERROR:Ошибка получения сетевых интерфейсов:", err)
    }
	
	//Загрузка конфига
	err = config.LoadFiles("config/sniffing_config.json")
	if err != nil {
		log.Fatalln("ERROR: Ошибка загрузки конфига:", err)
	}
	
	snaplen := config.Int("snaplen")
    promisc := config.Bool("promisc")
	computerName := config.String("computerName")
	timeout := config.Int("timeout")

	cfg := &Config{
		ComputerName: computerName,
		Snaplen: snaplen,
		Promisc: promisc,
		Timeout: timeout,
	}
	
	if timeout == 0 {
		cfg.Timeout = 1000
	}

	if snaplen == 0 {
		cfg.Snaplen = 1600
	}

	if computerName == "" {
		cfg.ComputerName = "myComputer"
	}
	
	//Запуск прослушивания интерфейсов
	wg := new(sync.WaitGroup)

	for _, device := range devices {
		wg.Add(1)
		go sniffer(device, wg, cfg)
	}
	wg.Wait()

	//Закрытие логирования
	log.Println("Info: Программа завершена")
}