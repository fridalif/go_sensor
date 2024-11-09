package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
    "log"
    "fmt"
    "time"
	"github.com/gorilla/websocket"
    "webinterface/models"
)


var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true 
    },
}

func WSHandler(c *gin.Context, db *gorm.DB) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
        log.Println("Ошибка при установлении WebSocket-соединения:", err)
        return
    }
    defer conn.Close()

    for {
        
        /*messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Ошибка при чтении сообщения:", err)
            break
        }
        fmt.Printf("Получено сообщение: %s\n", message)*/

        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            break
        }
    }
}


func Index(c *gin.Context, db *gorm.DB) {
    
    /* testovie dannie
    netLayer := models.Layer{Name:"IPv4",}
    result := db.Create(&netLayer)
    if result.Error != nil {
        fmt.Println(result.Error)
    }
    rule := models.Rule{
        NetlayerID: netLayer.ID,
        Netlayer:   netLayer,
        SrcIp:      "244.178.44.111",
        DstIp:      "244.178.44.111",
        TTL:        64,
        Checksum:   0,
        SrcPort:    "*",
        DstPort:    "*",
        PayloadContains: "Hello, World!",
    }
    result = db.Create(&rule)
    if result.Error != nil {
        fmt.Println(result.Error)
    }
    includedComputer := models.IncludedComputer{
        Name:    "Computer1",
        Address: "244.178.44.111",
    }
    result = db.Create(&includedComputer)
    if result.Error != nil {
        fmt.Println(result.Error)
    }
    alert := models.Alert{
        ComputerID: includedComputer.ID,
        Computer:   includedComputer,
        RuleID:     rule.ID,
        Rule:       rule,
        Timestamp:  time.Now(),
    }
    result = db.Create(&alert)
    if result.Error != nil {
        fmt.Println(result.Error)
    }*/
	c.HTML(
	  http.StatusOK,
	  "index.html",
	  gin.H{
		"title": "Home Page",
		"db":    db,
	  },
	)
}
