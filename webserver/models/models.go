package models

import (
	"gorm.io/driver/postgres"
	"errors"
	"gorm.io/gorm"
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
	IHL        int64 "IHL"
	Protocol   string
	TTL        int64
	TOS        int64
	Checksum   int64
	SrcPort    string "SrcPort"
	DstPort    string "DstPort"
	Seq        int64 "Seq"
	Ack        int64 "Ack"
	DataOffset int64 "DataOffset"
	FIN        bool "FIN"
	SYN        bool "SYN"
	RST        bool "RST"
	PSH        bool "PSH"
	ACK        bool "ACK"
	URG        bool "URG"
	ECE        bool "ECE"
	CWR        bool "CWR"

	PayloadContains string "PayloadContains"

}



func InitDB() (*gorm.DB, error) {
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