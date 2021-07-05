package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gobangSercice/middleware"
	_ "gobangSercice/model"
	"log"
	"net/http"
)

type connection struct {
	ws *websocket.Conn

	role uint8
}

var Users []*connection

var currentRoleStatus = map[string]bool{
	"white": false,
	"black": false,
}

var currentRound = 0

type piecePoint struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Color int `json:"color"`
	Role  int `json:"role"`
}

type Message struct {
	Code  uint8       `json:"code"`
	Role  uint8       `json:"role"`
	Round int         `json:"round"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

var pieceData = make([]*piecePoint, 0)

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

	link := &connection{
		ws:   ws,
		role: assigningRoles(),
	}

	Users = append(Users, link)

	for {
		// 读取数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		if string(message) == "ping" {

			for i := 0; i < len(Users); i++ {

				msg := Message{
					Code:  1,
					Role:  Users[i].role,
					Msg:   getUserStatus(),
					Round: currentRound,
					Data:  pieceData,
				}
				content, err := json.Marshal(msg)
				if err != nil {
					log.Println(err.Error())
					break
				}

				err = Users[i].ws.WriteMessage(mt, content)
				if err != nil {
					// 客户端退出了socket
					role := Users[i].role
					switch role {
					case 1:
						currentRoleStatus["white"] = false
					case 2:
						currentRoleStatus["black"] = false
					}
					Users = append(Users[:i], Users[i+1:]...)
					fmt.Println(err)
					break
				}
			}

		} else {

			point := &piecePoint{}

			err = json.Unmarshal(message, point)

			//currentRound = point.Color

			if err != nil {
				log.Println(err.Error())
				break
			}

			pieceData = append(pieceData, point)

			editCurrentRound(point.Color)

			for i := 0; i < len(Users); i++ {
				msg := Message{
					Code:  0,
					Role:  Users[i].role,
					Msg:   getUserStatus(),
					Round: currentRound,
					Data:  pieceData,
				}
				content, err := json.Marshal(msg)
				if err != nil {
					log.Println(err.Error())
					break
				}

				err = Users[i].ws.WriteMessage(mt, content)
				if err != nil {
					// 客户端退出了socket
					role := Users[i].role
					switch role {
					case 1:
						currentRoleStatus["white"] = false
					case 2:
						currentRoleStatus["black"] = false
					}
					Users = append(Users[:i], Users[i+1:]...)
					fmt.Println(err)
					break
				}
			}

		}
	}
}

func editCurrentRound(role int) {
	//fmt.Println("-----")
	//fmt.Println(role)
	if role == 1 {
		currentRound = 2
	} else {
		currentRound = 1
	}
}

func getUserStatus() string {
	if !currentRoleStatus["white"] && !currentRoleStatus["black"] {
		return "游戏结束"
	} else if !currentRoleStatus["black"] {
		return "黑方已经退出游戏或未加入请等待"
	} else if !currentRoleStatus["white"] {
		return "白方已经退出游戏或未加入请等待"
	}
	return ""
}

// 分配角色 1白 2黑 3观众
func assigningRoles() uint8 {
	if len(Users) == 0 {
		currentRoleStatus["white"] = true
		return 1
	}
	if currentRoleStatus["white"] && !currentRoleStatus["black"] {
		currentRoleStatus["black"] = true
		currentRound = 1
		return 2
	}
	return 3
}

func main() {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.GET("/ping", ping)

	r.GET("/refresh", Refresh)

	err := r.Run("0.0.0.0:8881")
	if err != nil {
		return
	}
}

func Refresh(context *gin.Context) {
	pwd := context.DefaultQuery("pwd", "")

	if pwd == "888888" {
		currentRoleStatus["black"] = false
		currentRoleStatus["white"] = false

		pieceData = pieceData[0:0]
		currentRound = 0
		Users = Users[0:0]
	}

	context.JSON(http.StatusOK, gin.H{
		"msg": pwd,
	})
}
