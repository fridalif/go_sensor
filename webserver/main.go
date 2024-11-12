package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"webinterface/views"
	"webinterface/models"
	"os"
	"strconv"
	"github.com/gookit/config/v2"
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
	router.GET("/ws", func (c *gin.Context) {views.WSHandler(c, db)},)
	router.GET("/sensorconn", func (c *gin.Context) {views.GetRules(c, db)},)
	router.POST("/api/add_rule", func(c *gin.Context){views.AddRule(c,db)},)
	router.DELETE("/api/delete_rule", func(c *gin.Context){views.DeleteRule(c,db)},)
	

	err = config.LoadFiles("config.json")
	if err != nil {
		log.Fatalln("ERROR: Ошибка загрузки конфига:", err)
	}
	port := config.Int("port")
	address := config.String("address")
	if address == "" {
		address = "127.0.0.1"
	}
	if port == 0 {
		port = 9000
	}
	err = router.Run(address + ":" + strconv.Itoa(port))
	if err != nil {
    	log.Fatalln("ERROR: Ошибка запуска сервера:", err)
	}
}