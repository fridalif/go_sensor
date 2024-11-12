package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
    "log"
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

//каналы для синхронизации
var alertsChanel = make(chan models.Alert)
var compChanel = make(chan models.IncludedComputer)
var rulesChanel = make(chan models.Rule)
var deleteRules = make(chan int)


var clients = make([]*websocket.Conn, 0)
var sensors = make([]*websocket.Conn, 0)


func closeSensor(conn *websocket.Conn) {
    conn.Close()
    for i, c := range sensors {
        if c == conn {
            sensors = append(sensors[:i], sensors[i+1:]...)
            break
        }
    }
}
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
                if err := compConnection.WriteJSON(message); err != nil {
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
                if err := compConnection.WriteJSON(message); err != nil {
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
                if err := compConnection.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
            for _, sensor := range sensors {
                if err := sensor.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
        case deletedID := <-deleteRules:
            message := map[string]interface{}{
                "TableName":"delete_rule",
                "Id": deletedID,
            }
            for _, compConnection := range clients {
                if err := compConnection.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
            for _, sensor := range sensors {
                if err := sensor.WriteJSON(message); err != nil {
                    log.Println("Ошибка при отправке сообщения:", err)
                    return
                }
            }
        }
    }
    
}

func GetRules(c *gin.Context, db *gorm.DB) {
    var rules []models.Rule
    if err := db.Preload("Netlayer").Find(&rules).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Ошибка при установлении WebSocket-соединения:", err)
        return
    }
    sensors = append(sensors, conn)
    defer closeSensor(conn)
    var newComp = map[string]interface{}{}
    if err := conn.ReadJSON(&newComp); err != nil {
        log.Println("Ошибка при чтении сообщения:", err)
        return
    }
    var computers []models.IncludedComputer
    if err := db.Find(&computers).Error; err != nil {
        log.Println("Ошибка при получении записей:", err)
        return
    }
    var found bool = false
    var address string = c.Request.RemoteAddr
    for _, computer := range computers {
        if computer.Address == address {
            found = true
            break
        }
    }
    if !found {
        newComputerModel := models.IncludedComputer{
            Name:    newComp["name"].(string),
            Address: address,
        }
        if err := db.Create(&newComputerModel).Error; err != nil {
            log.Println("Ошибка при создании записи:", err)
            return
        }
        compChanel <- newComputerModel
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
    //вечный цикл
    for {
        var newAlert = map[string]interface{}{}
        if err := conn.ReadJSON(&newAlert); err != nil {
            log.Println("Ошибка при чтении сообщения:", err)
            
        }
        
    }
}
func Index(c *gin.Context, db *gorm.DB) {
 	c.HTML(
	  http.StatusOK,
	  "index.html",
	  gin.H{
		"title": "Home Page",
	  },
	)
}
