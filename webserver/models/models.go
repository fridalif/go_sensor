package models

import (
	"gorm.io/driver/postgres"
	"errors"
	"gorm.io/gorm"
	"log"
)

//Инициализация логирования
file, err := os.OpenFile("../database.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
defer file.Close()

if err != nil {
	log.Fatalln("ERROR:Ошибка при открытии файла:", err)
}
log.SetOutput(file)

func initDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=webinterface port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Не удалось подключиться к базе данных")
	}

	return db, nil
}