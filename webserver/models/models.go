package models

import (
	"gorm.io/driver/postgres"
	"errors"
	"gorm.io/gorm"
	"time"
)


type IncludedComputer struct {
    gorm.Model
    Name    string
    Address string
}

type Layer struct {
    gorm.Model
    Name string `gorm:"unique;not null"`
}

type Rule struct {
    gorm.Model
	NetlayerID uint
    Netlayer   Layer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SrcIp      string
	DstIp      string
	TTL        int64
	Checksum   int64
	SrcPort    string "SrcPort"
	DstPort    string "DstPort"
	PayloadContains string "PayloadContains"
}

type Alert struct {
	gorm.Model
	ComputerID uint
	Computer   IncludedComputer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RuleID uint
	Rule    Rule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Timestamp time.Time
}



func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=webinterface port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Не удалось подключиться к базе данных")
	}
	err = db.AutoMigrate(&IncludedComputer{}, &Layer{}, &Rule{}, &Alert{})
	if err != nil {
		return nil, errors.New("Не удалось создать таблицы в базе данных")
	}
	return db, nil
}