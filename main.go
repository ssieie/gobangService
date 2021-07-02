package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "gobangSercice/model"
	"net/http"
)

//func init()  {
//	err := model.RedisInit()
//	if err != nil {
//		panic(err.Error())
//	}
//}

type connection struct {
	ws *websocket.Conn

	send chan []byte

	h *hup
}

type hup struct {
	connections map[*connection]bool
}

type piecePoint struct {
	x     int
	y     int
	color int
}

var pieceData = make([]*piecePoint, 10)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ping(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

		}
	}(ws)

	for {
		// 读取数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "ping" {
			message = []byte("pong")
			err = ws.WriteMessage(mt, message)
			if err != nil {
				break
			}
			break
		}

		point := &piecePoint{}

		err = json.Unmarshal(message, point)
		if err != nil {
			panic(err)
			//log.Println(err.Error())
			//break
		}

		pieceData = append(pieceData, point)

		content, err := json.Marshal(pieceData)
		if err != nil {
			panic(err)
			//log.Println(err.Error())
			//break
		}

		fmt.Println(content)

		err = ws.WriteMessage(mt, content)
		if err != nil {
			break
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ping", ping)

	err := r.Run("127.0.0.1:8888")
	if err != nil {
		return
	}
}
