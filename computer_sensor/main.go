package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	//Инициализация логирования
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("ERROR:Ошибка при открытии файла:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	file, err := os.Open("testfile.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatal(err)
	}

	hashSum := hasher.Sum(nil)

	fmt.Printf("SHA-1: %x\n", hashSum)
}
