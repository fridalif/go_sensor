package main

import (
	"os"
	//	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	//	"encoding/json"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gookit/config/v2"
	"github.com/gorilla/websocket"
)

// Конфиг
type Config struct {
	ComputerName string        `json:"computerName"`
	Snaplen      int           `json:"snaplen"`
	Promisc      bool          `json:"promisc"`
	Timeout      time.Duration `json:"timeout"`
	ServerAddr   string        `json:"server_addr"`
}

// Правило
type Rule struct {
	ID         uint                   `json:"id"`
	Layer      string                 `json:"layer"`
	Definition map[string]interface{} `json:"definition"`
}

var (
	rules []Rule
)

var mu sync.Mutex

func sendAlert(ruleId uint, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	err := conn.WriteJSON(map[string]interface{}{
		"rule_id":   ruleId,
		"timestamp": time.Now(),
	})
	if err != nil {
		log.Fatalln("ERROR: Не получилось отпраивть сработку серверу:", err)
		return
	}
}

func initRules(cfg *Config, conn *websocket.Conn) {

	// Устанавливаем соединение с WebSocket-сервером

	err := conn.WriteJSON(map[string]interface{}{
		"name": cfg.ComputerName,
	})
	if err != nil {
		log.Fatalln("ERROR: Не получилось отпраивть идентификационные данные серверу:", err)
		return
	}

	for {
		var serverMessage = map[string]interface{}{}
		if err := conn.ReadJSON(&serverMessage); err != nil {
			log.Println("ERROR:Ошибка при чтении сообщения:", err)
			return
		}
		tableName, exists := serverMessage["table_name"]
		var rule Rule
		if !exists {
			log.Println("ERROR: Не удалось получить имя таблицы")
			continue
		}
		if tableName == "rules" || tableName == "new_rule" {
			var ruleJSON map[string]interface{}
			if ruleJSON, exists = serverMessage["data"].(map[string]interface{}); !exists {
				log.Println("ERROR: Не удалось преобразовать правило в JSON")
				continue
			}
			layer, exists := ruleJSON["Netlayer"].(map[string]interface{})["Name"].(string)
			if !exists || layer == "" || (layer != "TCP" && layer != "IPv4" && layer != "IPv6") {
				log.Println("ERROR: Не удалось получить имя слоя")
				continue
			}
			rule.Layer = layer
			if _, exists = ruleJSON["ID"]; !exists {
				log.Println("ERROR: Не удалось преобразовать правило в JSON")
				continue
			}
			rule.ID = uint(ruleJSON["ID"].(float64))
			definition := map[string]interface{}{}
			for key, value := range ruleJSON {
				if key == "SrcIp" {
					definition["SrcIp"] = value
				}
				if key == "DstIp" {
					definition["DstIp"] = value
				}
				if key == "TTL" {
					definition["TTL"] = int64(value.(float64))
				}
				if key == "Checksum" {
					definition["Checksum"] = int64(value.(float64))
				}
				if key == "SrcPort" {
					definition["SrcPort"] = value
				}
				if key == "DstPort" {
					definition["DstPort"] = value
				}
				if key == "PayloadContains" {
					definition["PayloadContains"] = value
				}
			}
			rule.Definition = definition
			rules = append(rules, rule)
			log.Printf("INFO: Правило %d было добавлено", rule.ID)
		}
		if tableName == "delete_rule" {
			var id int
			if id, exists = serverMessage["Id"].(int); !exists {
				log.Println("ERROR: Не удалось преобразовать правило в JSON")
				continue
			}

			for i, rule := range rules {
				if rule.ID == uint(id) {
					rules = append(rules[:i], rules[i+1:]...)
					log.Printf("INFO: Правило %d было удалено", id)
					break
				}
			}
		}
	}
}

// Проверка парвил TCP
func checkTCP(tcpLayer gopacket.Layer, conn *websocket.Conn) bool {
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
				if value != -1 && uint32(value.(int64)) != tcp.Seq {
					thisRule = false
					break
				}
			case "Ack":
				if value != -1 && uint32(value.(int64)) != tcp.Ack {
					thisRule = false
					break
				}
			case "DataOffset":
				if value != -1 && uint8(value.(int64)) != tcp.DataOffset {
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
			}
		}
		if thisRule {
			log.Println("INFO: Правило TCP прошло проверку")
			sendAlert(rule.ID, conn)
			return true
		}
	}
	return false
}

// Проверка парвил IPv4
func checkIPv4(ipLayer gopacket.Layer, conn *websocket.Conn) bool {
	ipv4, ok := ipLayer.(*layers.IPv4)
	if !ok {
		log.Println("ERROR: Ошибка преобразования к типу IP")
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
			case "TTL":
				if value.(int64) != -1 && ipv4.TTL != uint8(value.(int64)) {
					thisRule = false
					break
				}
			case "Checksum":
				if value.(int64) != -1 && uint16(value.(int64)) != ipv4.Checksum {
					thisRule = false
					break
				}
			}
		}
		if thisRule {
			//log.Println("INFO:Правило прошло проверку")
			sendAlert(rule.ID, conn)
			return true
		}
	}
	return false
}

func sniffer(iface string, wg *sync.WaitGroup, cfg *Config, conn *websocket.Conn) {
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
			ipLayer := packet.Layer(layers.LayerTypeIPv4)

			if ipLayer != nil {

				checkIPv4(ipLayer, conn)
			}
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer != nil {
				checkTCP(tcpLayer, conn)
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
		Snaplen:      snaplen,
		Promisc:      promisc,
		Timeout:      pcap.BlockForever,
		ServerAddr:   serverAddr,
	}

	if snaplen == 0 {
		cfg.Snaplen = 1600
	}

	if computerName == "" {
		cfg.ComputerName = "myComputer"
	}

	if serverAddr == "" {
		cfg.ServerAddr = "127.0.0.1:9000"
	}

	conn, _, err := websocket.DefaultDialer.Dial(cfg.ServerAddr, nil)
	if err != nil {
		log.Fatal("ERROR:Ошибка подключения к вебу:", err)
		os.Exit(1)
	}

	defer conn.Close()
	log.Println("Info: Программа запущена")
	//Запуск прослушивания интерфейсов
	go initRules(cfg, conn)
	wg := new(sync.WaitGroup)

	for _, device := range devices {
		wg.Add(1)

		go sniffer(device.Name, wg, cfg, conn)
	}
	wg.Wait()

	//Закрытие логирования
	log.Println("Info: Программа завершена")
}
