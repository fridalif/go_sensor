package main

import (
	"os"
	"fmt"
	"log"
	"sync"
	"time"
	"strings"
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
	
}


//Проверка парвил TCP
func checkTCP(tcpLayer gopacket.Layer) bool {
	tcp, ok := tcpLayer.(*layers.TCP)
	
	if !ok {
		log.Println("ERROR: Ошибка преобразования к типу TCP")
		return false
	}
	payload := string(tcp.Payload)
	for _, rule := range rules {
		if rule.Layer != "TCP" {
			continue
		}
		thisRule := true
		for key, value := range rule.Definition {
			switch key {
				case "SrcPort":
					if tcp.SrcPort.String() != value && value != "*" {
						thisRule = false
						break
					}
				case "DstPort":
					if tcp.DstPort.String() != value && value != "*" {
						thisRule = false
						break
					}
				case "Seq":
					if value != -1 && uint32(value.(int)) != tcp.Seq {
						thisRule = false
						break
					}
				case "Ack":
					if value != -1 && uint32(value.(int)) != tcp.Ack {
						thisRule = false
						break
					}
				case "DataOffset":
					if value != -1 && uint8(value.(int)) != tcp.DataOffset {
						thisRule = false
						break
					}
				case "FIN":
					if value != -1 && tcp.FIN != value.(bool) {
						thisRule = false
						break
					}
				case "SYN":
					if value != -1 && tcp.SYN != value.(bool) {
						thisRule = false
						break
					}
				case "RST":
					if value != -1 && tcp.RST != value.(bool) {
						thisRule = false
						break
					}
				case "PSH":
					if value != -1 && tcp.PSH != value.(bool) {
						thisRule = false
						break
					}
				case "ACK":
					if value != -1 && tcp.ACK != value.(bool) {
						thisRule = false
						break
					}
				case "URG":
					if value != -1 && tcp.URG != value.(bool) {
						thisRule = false
						break
					}
				case "ECE":
					if value != -1 && tcp.ECE != value.(bool) {
						thisRule = false
						break
					}
				case "CWR":
					if value != -1 && tcp.CWR != value.(bool) {
						thisRule = false
						break
					}
				case "NS":
					if value != -1 && tcp.NS != value.(bool) {
						thisRule = false
						break
					}
				case "PayloadContains":
					if value.(string) != "*" && !strings.Contains(strings.ToLower(payload), strings.ToLower(value.(string))) {
						thisRule = false
						break	
					}
				default:
					thisRule = false
					log.Println("ERROR: Неизвестный ключ в правиле TCP:", key)
			}
		}
		if thisRule {
			fmt.Println("INFO: Правило TCP прошло проверку")
			return true
		}
	}
	return false
}
//Проверка парвил IPv4
func checkIPv4(ipLayer gopacket.Layer) bool {
	ipv4, ok := ipLayer.(*layers.IPv4)
	if !ok {
		log.Println("ERROR: Ошибка преобразования к типу IPv4")
		return false
	}
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
				case "DstIp":
					if ipv4.DstIP.String() != value && value != "*" {
						thisRule = false
						break
					}
				case "Protocol":
					if ipv4.Protocol.String() != value && value != "*" {
						thisRule = false
						break
					}
				case "IHL": 
					if value != -1 && uint8(value.(int)) != ipv4.IHL {
						thisRule = false
						break
					}
				case "TOS": 
					if value != -1 && uint8(value.(int)) != ipv4.TOS {
						thisRule = false
						break
					}
				case "Length": 
					if value != -1 && uint16(value.(int)) != ipv4.Length {
						thisRule = false
						break
					}
				case "TTL": 
					if value != -1 && ipv4.TTL != uint8(value.(int)) {
						fmt.Println(ipv4.TTL, value)
						thisRule = false
						break
					}
				case "Checksum": 
					if value != -1 && uint16(value.(int)) != ipv4.Checksum {
						thisRule = false
						break
					}
				default:
					thisRule = false
					log.Println("ERROR: Неизвестный ключ в правиле IPv4:", key)
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

	//Запуск прослушивания
	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		wg.Add(1)
		go func(packet gopacket.Packet) {
			defer wg.Done()
			/*ipLayer := packet.Layer(layers.LayerTypeIPv4)
        	if ipLayer != nil {
				checkIPv4(ipLayer)
        	}*/
			/*ipv6Layer := packet.Layer(layers.LayerTypeIPv6)
			if ipv6Layer != nil {
				checkIPv6(ipv6Layer)
			}*/
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				checkTCP(tcpLayer)
			}
		}(packet)
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
	serverAddr := config.String("serverAddr")
	cfg := &Config{
		ComputerName: computerName,
		Snaplen: snaplen,
		Promisc: promisc,
		Timeout: pcap.BlockForever,
		ServerAddr: serverAddr,
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