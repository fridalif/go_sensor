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

var hashRules []string

func initRules() {
	hashRules = append(hashRules, "70f32f84fe08d19204d9e31f7a885451ed9af344")
}

func checkFile(filePath string) {
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
		if fmt.Sprintf("%x", hashSum) == rule {
			fmt.Printf("File %s matches rule %s\n", filePath, rule)
		}
	}
}

func checkDir(checkingDir string, interval int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		err := filepath.Walk(checkingDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				go checkFile(path)
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
		log.Fatal(err)
	}
	defer conn.Close()
	if err := conn.WriteMessage(websocket.TextMessage, []byte("Hello")); err != nil {
		log.Fatal(err)
	}

	wg := new(sync.WaitGroup)
	initRules()
	for _, directory := range directories {
		wg.Add(1)
		go checkDir(directory, interval, wg)
	}
	wg.Wait()
}
