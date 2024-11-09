package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true 
    },
}

func WSHandler(с *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
        log.Println("Ошибка при установлении WebSocket-соединения:", err)
        return
    }
    defer conn.Close()

    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("Ошибка при чтении сообщения:", err)
            break
        }
        fmt.Printf("Получено сообщение: %s\n", message)

        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Println("Ошибка при отправке сообщения:", err)
            break
        }
    }
}


func Index(c *gin.Context, db *gorm.DB) {
	c.HTML(
	  http.StatusOK,
	  "index.html",
	  gin.H{
		"title": "Home Page",
		"db":    db,
	  },
	)
}
