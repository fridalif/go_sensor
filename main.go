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
	"github.com/google/gopacket/layers"
)

//Конфиг
type Config struct {
	ComputerName string `json:"computerName"`
	Snaplen int  `json:"snaplen"`
	Promisc bool `json:"promisc"`
	Timeout time.Duration `json:"timeout"`
}

//Правило
type Rule struct {
	Layer string `json:"layer"`
	Definition map[string]interface{} `json:"definition"`
}


var ( rules []Rule )

func initRules() {
	firstRule := Rule{
		Layer: "IPv4",
		Definition: map[string]interface{}{
			"SrcIp": "127.0.0.1",
			"DstIp": "127.0.0.1",
			"TTL": 64,
		},
	}

	rules = append(rules, firstRule)
}
func checkIPv4(ipLayer gopacket.Layer) bool {
	ipv4, ok := ipLayer.(*layers.IPv4)
	if !ok {
		log.Println("ERROR: Ошибка преобразования к типу IPv4")
		return false
	}
	/*
	Version    uint8
	IHL        uint8
	TOS        uint8
	Length     uint16
	TTL        uint8
	Protocol   IPProtocol
	Checksum   uint16
	SrcIP      net.IP
	DstIP      net.IP
	*/
	for _, rule := range rules {
		if rule.Layer != "IPv4" {
			continue
		}
		thisRule := true
		for key, value := range rule.Definition {
			switch key {
				case "SrcIp":
					if ipv4.SrcIP.String() != value && value != "*" {
						thisRule = false
						break
					}
					continue
				case "DstIp":
					if ipv4.DstIP.String() != value && value != "*" {
						thisRule = false
						break
					}
					continue
				case "Protocol":
					if ipv4.Protocol.String() != value && value != "*" {
						thisRule = false
						break
					}
					continue
				case "IHL": 
					if value != -1 && uint8(value.(int)) != ipv4.IHL {
						thisRule = false
						break
					}
					continue
				case "TOS": 
					if value != -1 && uint8(value.(int)) != ipv4.TOS {
						thisRule = false
						break
					}
					continue
				case "Length": 
					if value != -1 && uint16(value.(int)) != ipv4.Length {
						thisRule = false
						break
					}
					continue
				case "TTL": 
					if value != -1 && ipv4.TTL != uint8(value.(int)) {
						fmt.Println(ipv4.TTL, value)
						thisRule = false
						break
					}
					continue
				case "Checksum": 
					if value != -1 && uint16(value.(int)) != ipv4.Checksum {
						thisRule = false
						break
					}
					continue
				default:
					thisRule = false
					log.Println("ERROR: Неизвестный ключ в правиле:", key)
			}
		}
		if thisRule {
			fmt.Println("Правило прошло проверку")
			return true
		}
	}
	return false
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
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
        if ipLayer != nil {
			checkIPv4(ipLayer)
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
	initRules()
	wg := new(sync.WaitGroup)

	for _, device := range devices {
		wg.Add(1)
		
		go sniffer(device.Name, wg, cfg)
	}
	wg.Wait()

	//Закрытие логирования
	log.Println("Info: Программа завершена")
}