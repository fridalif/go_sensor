package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"webinterface/views"
	"webinterface/models"
	"os"
)
  
var router *gin.Engine

func main() {
	//Инициализация логирования
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer file.Close()
	
	if err != nil {
    	log.Fatalln("ERROR:Ошибка при открытии файла:", err)
	}
	log.SetOutput(file)

	//Инициализация базы данных
	db, err := models.InitDB()
	if err != nil {
		log.Fatalln("ERROR: Ошибка подключения к базе данных:", err)
	}

	//Инициализация роутера
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func (c *gin.Context) {views.Index(c, db)},)
  
	err = router.Run(":9000")
	if err != nil {
    	log.Fatalln("ERROR: Ошибка запуска сервера:", err)
	}
}