package models

import (
	"gorm.io/driver/postgres"
	"errors"
	"gorm.io/gorm"
)

type IncludedComputer struct {
	gorm.Model
	Name string
	Address string
}

type Layer struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Rule struct {
	gorm.Model
	NetLayer Layer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Definition map[string]interface{}
}



func initDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=webinterface port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Не удалось подключиться к базе данных")
	}
	err = db.AutoMigrate(&IncludedComputer{}, &Layer{}, &Rule{})
	if err != nil {
		return nil, errors.New("Не удалось создать таблицы в базе данных")
	}
	return db, nil
}