package webbsocket

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type newClient struct {
	UserID string
	Conn   *websocket.Conn
}

var UserConnections = make(map[string]*newClient)

var ErrRecipientNotFound = errors.New("recipient not found")

func HandleWebsocket(c *gin.Context) {
	conn, err := InitializeWebsocketWIthUser(c)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, msgBytes, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if messageType == websocket.TextMessage {
			var msg map[string]string
			if err := json.Unmarshal(msgBytes, &msg); err != nil {
				log.Println(err)
				continue
			}

			toUserID, found := msg["to"]
			message, foundMessage := msg["message"]

			if found && foundMessage {
				err := SendMessageToUser(toUserID, message)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

}

func InitializeWebsocketWIthUser(c *gin.Context) (*websocket.Conn, error) {
	userID := c.Param("UserID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	UserConnections[userID] = &newClient{
		UserID: userID,
		Conn:   conn,
	}

	return conn, nil
}

func SendMessageToUser(userID string, msg string) error {
	recipientConn, found := UserConnections[userID]
	if !found {
		return ErrRecipientNotFound
	}

	message := map[string]string{
		"message": msg,
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if err := recipientConn.Conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		return err
	}

	return nil
}
