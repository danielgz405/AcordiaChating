package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/dg/acordia/models"
	"github.com/dg/acordia/repository"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//Validate if the user is authorized to make requests, it's important to work with the token.
		return true
	},
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) HandleWebSocket(JWTSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}
		client := NewClient(hub, socket)

		// Get the value of the parameter sent in the URL
		params := mux.Vars(r)
		tokenString := strings.TrimSpace(params["Authorization"])
		_, err = ValidateTokenAndGetProfile(JWTSecret, tokenString, r.Context())
		if err != nil {
			http.Error(w, "Error validating token", http.StatusUnauthorized)
			return
		}
		client.id = tokenString
		client.channel = params["Channel"]

		hub.register <- client

		go func() {
			for {
				_, _, err := client.socket.ReadMessage()
				if err != nil {
					hub.unregister <- client
					break
				}
			}
		}()

		go client.Write()
	}
}
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.onConnect(client)
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

func (hub *Hub) onConnect(client *Client) {
	fmt.Println("Connecting")

	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	hub.clients = append(hub.clients, client)
}

func (hub *Hub) onDisconnect(client *Client) {
	fmt.Println("Disconnecting")

	client.socket.Close()
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
		}
	}

	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]
}

func (hub *Hub) Broadcast(message interface{}, channels []string) {
	fmt.Println("send message")
	fmt.Println(hub.clients)

	data, _ := json.Marshal(message)
	for _, client := range hub.clients {
		if ValidateChannel(client.channel, channels) {
			client.outbound <- data
		}
	}
}

func ValidateTokenAndGetProfile(JWTSecret string, tokenString string, ctx context.Context) (*models.Profile, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
		userId := claims.UserId.Hex()
		profile, err := repository.GetUserById(ctx, userId)
		if err != nil {
			return nil, err
		}
		return profile, nil
	} else {
		return nil, err
	}
}

func ValidateChannel(channel string, channels []string) bool {
	for _, currchannel := range channels {
		if currchannel == channel {
			return true
		}
	}
	return false
}
