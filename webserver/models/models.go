package models

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gookit/config/v2"
	"gorm.io/driver/postgres"
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
	NetlayerID      uint
	Netlayer        Layer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SrcIp           string
	DstIp           string
	TTL             int64
	Checksum        int64
	SrcPort         string
	DstPort         string
	PayloadContains string
}

type RuleComputer struct {
	gorm.Model
	HashSum string `gorm:"unique;not null" json:"hash_sum"`
}

type AlertComputer struct {
	gorm.Model
	ComputerID uint
	Computer   IncludedComputer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RuleID     uint
	Rule       RuleComputer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Timestamp  time.Time
}

type Alert struct {
	gorm.Model
	ComputerID uint
	Computer   IncludedComputer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RuleID     uint
	Rule       Rule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Timestamp  time.Time
}

func InitDB() (*gorm.DB, error) {
	err := config.LoadFiles("config.json")
	if err != nil {
		log.Fatalln("ERROR: Ошибка загрузки конфига:", err)
	}
	dbPort := config.Int("db_port")
	dbAddress := config.String("db_address")
	if dbAddress == "" {
		dbAddress = "localhost"
	}
	if dbPort == 0 {
		dbPort = 5432
	}
	dbName := config.String("db_name")
	if dbName == "" {
		dbName = "webinterface"
	}
	dbUser := config.String("db_user")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := config.String("db_password")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	dsn := "host=" + dbAddress + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + strconv.Itoa(dbPort) + " sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, errors.New("не удалось подключиться к базе данных")
	}
	err = db.AutoMigrate(&IncludedComputer{}, &Layer{}, &Rule{}, &Alert{}, &RuleComputer{}, &AlertComputer{})
	if err != nil {
		return nil, errors.New("не удалось создать таблицы в базе данных")
	}
	return db, nil
}
