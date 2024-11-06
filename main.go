package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gookit/config/v2"
)

type Config struct {
	ComputerName string `json:"computerName"`
	Snaplen int  `json:"snaplen"`
	Promisc bool `json:"promisc"`
	Timeout time.Duration `json:"timeout"`
}

func sniffer(iface string, wg *sync.WaitGroup, cfg *Config) {
	defer wg.Done()
	if iface == "dbus-system" || iface == "dbus-session" {
		return
	}
	handle, err := pcap.OpenLive(iface, int32(cfg.Snaplen), cfg.Promisc, cfg.Timeout)
    if err != nil {
        log.Println("ERROR: Ошибка открытия интерфейса:", err)
		return
    }
    defer handle.Close()
	//Если нужны фильтры
	//if err := handle.SetBPFFilter(filter); err != nil {
	//	log.Println("ERROR: Ошибка установки фильтра:",err)
	//	return
	//}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		fmt.Println(iface)
		fmt.Println(cfg.ComputerName)
		for _, layer := range packet.Layers() {
			fmt.Println(layer.LayerType())
		}
	}
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

	cfg := &Config{
		ComputerName: computerName,
		Snaplen: snaplen,
		Promisc: promisc,
		Timeout: pcap.BlockForever,
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
		
		go sniffer(device.Name, wg, cfg)
	}
	wg.Wait()

	//Закрытие логирования
	log.Println("Info: Программа завершена")
}