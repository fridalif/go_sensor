package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
    "log"
 //   "fmt"
  //  "time"
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

    var alerts []models.Alert
    if err := db.Preload("Computer").Preload("Rule").Preload("Rule.Netlayer").Order("timestamp desc").Find(&alerts).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }

    for _, alert := range alerts {
        if err := conn.WriteJSON(alert); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            return
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
