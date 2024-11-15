package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gookit/config/v2"
	"github.com/gorilla/websocket"
)

var hashRules []map[string]interface{}

func initRules(conn *websocket.Conn, computerName string) {
	err := conn.WriteJSON(map[string]interface{}{
		"name": computerName,
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
		if !exists {
			log.Println("ERROR: Не удалось получить имя таблицы")
			continue
		}
		if tableName == "rules" || tableName == "new_rule_computer" {
			var ruleJSON map[string]interface{}
			if ruleJSON, exists = serverMessage["data"].(map[string]interface{}); !exists {
				log.Println("ERROR: Не удалось преобразовать правило в JSON")
				continue
			}
			log.Println("INFO: Получено правило:", uint(ruleJSON["ID"].(float64)), ruleJSON["hash_sum"])
			hashRules = append(hashRules, map[string]interface{}{"hash_sum": ruleJSON["hash_sum"].(string), "id": uint(ruleJSON["ID"].(float64))})
		}
		if tableName == "delete_rule" {
			var id int
			if id, exists = serverMessage["Id"].(int); !exists {
				log.Println("ERROR: Не удалось преобразовать правило в JSON")
				continue
			}

			for i, rule := range hashRules {
				if rule["id"] == uint(id) {
					hashRules = append(hashRules[:i], hashRules[i+1:]...)
					log.Printf("INFO: Правило %d было удалено", id)
					break
				}
			}
		}
	}
}

func checkFile(filePath string, conn *websocket.Conn) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatal(err)
	}

	hashSum := hasher.Sum(nil)
	for _, rule := range hashRules {
		if fmt.Sprintf("%x", hashSum) == rule["hash_sum"] {
			log.Printf("INFO: File %s matches rule %s\n", filePath, rule)
			conn.WriteJSON(map[string]interface{}{"rule_id": rule["id"]})
		}
	}
}

func checkDir(checkingDir string, interval int, wg *sync.WaitGroup, conn *websocket.Conn) {
	defer wg.Done()
	for {
		err := filepath.Walk(checkingDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				go checkFile(path, conn)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func main() {

	//Инициализация логирования
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("ERROR:Ошибка при открытии файла:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	err = config.LoadFiles("config.json")
	if err != nil {
		log.Fatalln("ERROR: Ошибка загрузки конфига:", err)
	}

	directories := config.Strings("checking_directories")
	interval := config.Int("checking_interval")
	serverUrl := config.String("serverAddr")
	computerName := config.String("computerName")
	if computerName == "" {
		computerName = "Uzel"
	}
	if serverUrl == "" {
		serverUrl = "ws://127.0.0.1:9000/computerconn"
	}
	if interval == 0 {
		interval = 60
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl, nil)
	if err != nil {
		fmt.Println("ERROR: Не удалось подключиться к серверу:", serverUrl)
		log.Fatal(err)

	}
	defer conn.Close()

	wg := new(sync.WaitGroup)
	go initRules(conn, computerName)
	for _, directory := range directories {
		wg.Add(1)
		go checkDir(directory, interval, wg, conn)
	}
	wg.Wait()
}
