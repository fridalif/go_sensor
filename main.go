package main

import (
	"os"
	"fmt"
	"log"
	"github.com/google/gopacket/pcap"
	"sync"
)

func sniffer(iface Interface, wg *sync.WaitGroup) {
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
	wg := new(sync.WaitGroup)

	//Запуск прослушивания интерфейсов
	for _, device := range devices {
		wg.Add(1)
		go sniffer(device, wg)
	}
	wg.Wait()

	//Закрытие логирования
	log.Println("Info: Программа завершена")
}