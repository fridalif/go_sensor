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

type AlertMessage struct {
    TableName string `json:"table_name"`
    Data      models.Alert  `json:"data"`       
}

type ComputerMessage struct {
    TableName string `json:"table_name"`
    Data      models.IncludedComputer  `json:"data"`       
}

type RuleMessage struct {
    TableName string `json:"table_name"`
    Data      models.Rule  `json:"data"`
}

var alertsChanel = make(chan models.Alert)
var compChanel = make(chan models.IncludedComputer)
var rulesChanel = make(chan modles.Rule)

var clients = make([]*websocket.Conn, 0)
func closeConn(conn *websocket.Conn) {
    conn.Close()
    for i, c := range clients {
        if c == conn {
            clients = append(clients[:i], clients[i+1:]...)
            break
        }
    }
}

// handler обработки websocket
func WSHandler(c *gin.Context, db *gorm.DB) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
        log.Println("Ошибка при установлении WebSocket-соединения:", err)
        return
    }

    clients = append(clients, conn)

    defer closeConn(conn)
    
    var computers []models.IncludedComputer
    if err := db.Find(&computers).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }
    //Инициализация компютеров
    for _, computer := range computers {
        message := ComputerMessage{
            TableName: "computers",
            Data:      computer,
        }
        if err := conn.WriteJSON(message); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            return
        }
    }
    
    var rules []models.Rule
    if err := db.Preload("Netlayer").Find(&rules).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }
    //Инициализация правил
    for _, rule := range rules {
        message := RuleMessage{
            TableName: "rules",
            Data:      rule,
        }
        if err := conn.WriteJSON(message); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            return
        }
    }

    var alerts []models.Alert
    
    
    if err := db.Preload("Computer").Preload("Rule").Preload("Rule.Netlayer").Order("timestamp desc").Find(&alerts).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }
    //Инициализация алертов
    for _, alert := range alerts {
        message := AlertMessage{
            TableName: "alerts",
            Data:      alert,
        }
        if err := conn.WriteJSON(message); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            return
        }
    }
    
    //Ожидание подключения новых компьютеров или сработок
    for {
        select {
        case alert := <-alertsChanel:
            message := AlertMessage{
                TableName: "new_alerts",
                Data:      alert,
            }
            for _, compConnection := range clients {
                if err := compConenction.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
        case computer := <-compChanel:
            message := ComputerMessage{
                TableName: "new_computers",
                Data:      computer,
            }
            for _, compConnection := range clients {
                if err := compConenction.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
        case rule := <-rulesChanel:
            message := RuleMessage{
                TableName:"new_rule",
                Data: rule,
            }
            for _, compConnection := range clients {
                if err := compConenction.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
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
