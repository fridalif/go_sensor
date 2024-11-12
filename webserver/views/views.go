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
var deleteRules = make(chan uint)


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
    fmt.Println(clients)
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
                    fmt.Println("error")
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
        fmt.Println(computer.Address)
        fmt.Println(address)
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

    var myComputer models.IncludedComputer
    if err := db.Where("address = ?", address).First(&myComputer).Error; err != nil {
        log.Println("ERROR: Компьютер не найден:", err)
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
    //вечный цикл
    for {
        var newAlert = map[string]interface{}{}
        if err := conn.ReadJSON(&newAlert); err != nil {
            log.Println("Ошибка при чтении сообщения:", err)   
        }
        ruleId, exists := newAlert["rule_id"].(uint)
        if !exists {
            log.Println("ERROR: Не удалось получить ID правила")
            continue
        }
        var rule models.Rule
        if err := db.Where("id = ?", ruleId).First(&rule).Error; err != nil {
            log.Println("ERROR: Правило не найдено:", err)
            continue
        }
         
        newAlertModel := models.Alert{
            ComputerID: myComputer.ID,
            Computer:   myComputer,
            RuleID:     rule.ID,
            Rule:       rule,
            Timestamp:  time.Now(),
        }
        if timestamp, exists:= newAlert["timestamp"].(time.Time); exists {
            newAlertModel.Timestamp = timestamp
        }
        if err := db.Create(&newAlertModel).Error; err != nil {
            log.Println("ERROR: Ошибка при создании записи:", err)
            continue
        }
        alertsChanel <- newAlertModel
    }
}

func AddRule(c *gin.Context, db *gorm.DB) {
    if c.Request.Method != "POST" {
        c.AbortWithStatus(http.StatusMethodNotAllowed)
        return
    }
    var ruleInterface map[string]interface{}
    if err := c.BindJSON(&ruleInterface); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    var rule models.Rule
    var netlayer models.Layer
    if err := db.Where("name = ?", ruleInterface["netlayer"].(string)).First(&netlayer).Error; err != nil {
        response := gin.H{
            "status":  "error",
            "message": "Сетевой уровень не найден",
            "data":    nil,
        }
        c.JSON(http.StatusBadRequest, response)
        return
    }
    rule.Netlayer = netlayer
    rule.NetlayerID = netlayer.ID
    if srcIp, exists := ruleInterface["src_ip"].(string); exists{
        rule.SrcIp = srcIp
    } else {
        rule.SrcIp = "*"
    }
    if dstIp, exists := ruleInterface["dst_ip"].(string); exists{
        rule.DstIp = dstIp
    } else{
        rule.DstIp = "*"
    }
    if TTL, exists := ruleInterface["TTL"].(int64); exists{
        rule.TTL = TTL
    } else{
        rule.TTL = -1
    }
    if checksum, exists := ruleInterface["checksum"].(int64);exists{
        rule.Checksum = checksum
    } else{
        rule.Checksum = -1
    }
    if srcPort, exists := ruleInterface["src_port"].(string);exists{
        rule.SrcPort = srcPort
    } else{
        rule.SrcPort = "*"
    }
    if dstPort, exists := ruleInterface["dst_port"].(string); exists{
        rule.DstPort = dstPort
    } else{
        rule.DstPort = "*"
    }
    if payloadContains, exists := ruleInterface["payload_contains"].(string); exists{
        rule.PayloadContains = payloadContains
    } else{
        rule.PayloadContains = "*"
    }
    if err := db.Create(&rule).Error; err != nil {
        log.Println("ERROR: ошибка добавления в базу данных", err)
        c.AbortWithStatus(http.StatusInternalServerError)
        return
    }
    rulesChanel <- rule
    response := gin.H{
        "status":  "success",
        "message": "Правило успешно добавлено",
    }

    c.JSON(http.StatusOK, response)
}

func DeleteRule(c* gin.Context, db *gorm.DB){
    if c.Request.Method != "DELETE" {
        c.AbortWithStatus(http.StatusMethodNotAllowed)
        return
    }
    var ruleInterface map[string]interface{}
    if err := c.BindJSON(&ruleInterface); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    }
    id, exists := ruleInterface["rule_id"].(uint) 
    if !exists{
        response := gin.H{
            "status":"error",
            "message":"Не указан id правила",
            "data":nil,
        }
        c.JSON(http.StatusBadRequest, response)
        return
    }
    var rule models.Rule
    result := db.First(&rule, id)
    if result.Error != nil{
        response := gin.H{
            "status":"error",
            "message":"Нет правила с таким id",
            "data":nil,
        }
        c.JSON(http.StatusBadRequest, response)
        return
    }
    db.Delete(&rule)
    deleteRules <- rule.ID
    response := gin.H{
        "status":"Success",
        "message":"Удалено",
        "data":nil,
    }
    c.JSON(http.StatusOK, response)
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
