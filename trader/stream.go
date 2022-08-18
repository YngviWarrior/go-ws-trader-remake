package trader

import (
	"encoding/json"
	"errors"
	"fmt"
	entities "gowstrader/entities"
	mysql "gowstrader/mysql"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/websocket"
)

type Client struct {
	Id              int64
	SocketConn      *websocket.Conn
	IsAuthenticated bool
	SubscribedGames []*entities.GameWithBet
}

var ClientHub = make(map[int64]*Client)

func (c *Client) getUserSession(session string) (id int64, err error) {
	err = errors.New("not found id")

	searchUser := strings.Index(session, `"User"`)
	if searchUser < 0 {
		err = errors.New("not found string User")
		return
	}

	searchID := strings.Index(session[searchUser:], `s:2:"id"`)
	if searchID < 0 {
		err = errors.New("not found string id")
		return
	}

	splitDoublePoint := strings.Split(session[searchUser:][searchID+8:], ":")
	if len(splitDoublePoint) < 3 {
		err = errors.New("not found split two points")
		return
	}

	splitEnd := strings.Split(splitDoublePoint[2], `"`)
	if len(splitEnd) < 2 {
		err = errors.New("not found split \"")
		return
	}

	id, errorInt := strconv.ParseInt(splitEnd[1], 10, 64)
	if errorInt != nil {
		return
	}

	err = nil
	return
}

func (c *Client) SubscribeUser(msg *entities.Read, cache *memcache.Client, socket *websocket.Conn) (response entities.Write, err error) {
	response.Id = "1337"
	response.Endpoint = "/login"
	response.Error = true
	response.ErrorCode = "6"
	response.ErrorMessage = "Token invalid."
	response.Response = entities.GameWithBet{}

	r, err := regexp.Compile(`^[a-zA-Z0-9,\-]{10,32}$`)

	if err != nil {
		fmt.Println("Erro getSessionPHP regexp.compile", err)
		return
	}

	findMatch := r.FindStringSubmatch(msg.Parameters.Token)

	if len(findMatch) == 0 {
		err = errors.New("authtoken not are rules from regexp")
		return
	}

	authtoken := findMatch[0]
	cachedInfo, err := cache.Get("memc.sess.key." + authtoken)

	if err != nil {
		return
	}

	id, err := c.getUserSession(string(cachedInfo.Value))

	if err != nil {
		return
	}

	c.Id = id
	c.IsAuthenticated = true
	c.SocketConn = socket

	ClientHub[id] = c

	response.Error = false
	response.ErrorCode = ""
	response.ErrorMessage = ""

	return
}

func (c *Client) SubscribeGameInfo(dbConn *mysql.SqlConn, msg *entities.Read, cache *memcache.Client) (response entities.Write, err error) {
	response.Id = "7543452322"
	response.Endpoint = "/user/game/subscribe-game-info"
	response.Error = true
	response.ErrorCode = "4"
	response.ErrorMessage = "Internal Error"
	response.Response = entities.GameWithBet{}

	if !c.IsAuthenticated {
		response.ErrorCode = "7"
		response.ErrorMessage = "forbidden user not authenticated"
		return
	}

	idGame, _ := strconv.ParseInt(msg.Parameters.IdGame, 10, 64)

	game, err := dbConn.GetGame(idGame)

	if err != nil {
		fmt.Println(err)
	}

	if (game == entities.GameWithBet{}) {
		return
	}

	response.Id = "7543452322"
	response.Endpoint = "/user/game/subscribe-game-info"
	response.Error = false
	response.ErrorCode = ""
	response.ErrorMessage = ""
	response.Response = game

	c.SubscribedGames = append(c.SubscribedGames, &game)
	*ClientHub[c.Id] = *c

	return
}

func sendNotification(subsGame *entities.GameWithBet, c *websocket.Conn) {
	var response entities.Write

	if subsGame.GameIDStatus == 4 {
		response.Id = "updategameinfodone"
		response.Endpoint = "/user/game/update-game-info-done"
	} else {
		response.Id = "subscribeupdate2"
		response.Endpoint = "/user/game/update-game-info"
	}

	response.Error = false
	response.ErrorCode = "0"
	response.ErrorMessage = ""
	response.Response = *subsGame

	json, err := json.Marshal(response)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = c.WriteMessage(websocket.TextMessage, json)

	if err != nil {
		if !errors.Is(err, syscall.EPIPE) {
			log.Println("write kindle: ", err)
		}
	}
}

func hasUpdate(user *Client, gamesInfo *[]entities.Memcached, gameBetList []*entities.GameWithBet, dbConn *mysql.SqlConn) {
	for i := range user.SubscribedGames {
		for _, cacheGame := range *gamesInfo {
			gamesIds := strings.Split(cacheGame.IdGames, ",")

			idStatus, _ := strconv.ParseInt(cacheGame.GameIdStatus, 0, 64)
			for _, id := range gamesIds {
				if id == fmt.Sprintf("%v", ClientHub[user.Id].SubscribedGames[i].GameID) && ClientHub[user.Id].SubscribedGames[i].GameIDStatus != idStatus {
					for _, v := range gameBetList {
						if v.GameID == ClientHub[user.Id].SubscribedGames[i].GameID {
							ClientHub[user.Id].SubscribedGames[i] = v
						}
					}

					// ClientHub[user.Id].SubscribedGames[i].GameIDStatus = idStatus

					sendNotification(ClientHub[user.Id].SubscribedGames[i], ClientHub[user.Id].SocketConn)
				}
			}
		}
	}
}

func SubscribeUpdate(c *websocket.Conn, cache *memcache.Client, dbConn *mysql.SqlConn) {
	for range time.Tick(time.Second * 1) {
		resp, err := cache.Get("GamesUpdate")

		if err != nil {
			fmt.Println(err)
			return
		}

		if resp != nil && resp.Key == "GamesUpdate" && len(string(resp.Value)) > 0 {
			var gamesInfo []entities.Memcached
			err := json.Unmarshal(resp.Value, &gamesInfo)

			if err != nil {
				fmt.Println("Marshal Json: " + err.Error())
			}

			for _, user := range ClientHub {
				go func(user *Client) {
					if user.IsAuthenticated {
						var list []int64

						for _, v := range user.SubscribedGames {
							list = append(list, v.GameID)
						}

						gameBetList, _ := dbConn.GetInfoGameListWithBet(list, user.Id)

						hasUpdate(user, &gamesInfo, gameBetList, dbConn)
					}
				}(user)
			}

			i := memcache.Item{}
			i.Key = "GamesUpdate"
			i.Value = []byte{}
			err = cache.Set(&i)

			if err != nil {
				fmt.Println(err)
			}

		}

	}
}

func Subscribe(c *websocket.Conn, cache *memcache.Client, client *Client, dbConn *mysql.SqlConn) {
	for {
		// IsConnected()

		_, message, err := c.ReadMessage()

		if err != nil {
			fmt.Println(err)
		}

		var msg entities.Read
		_ = json.Unmarshal(message, &msg)

		switch msg.Endpoint {
		case "/login":
			response, err := client.SubscribeUser(&msg, cache, c)

			if err != nil {
				fmt.Println(err)
			}

			bytes, err := json.Marshal(response)

			if err != nil {
				log.Fatal("encode error:", err)
			}

			_ = c.WriteMessage(websocket.TextMessage, bytes)
		case "/user/game/subscribe-game-info":
			response, err := client.SubscribeGameInfo(dbConn, &msg, cache)

			if err != nil {
				fmt.Println(err)
			}

			bytes, err := json.Marshal(response)

			if err != nil {
				log.Fatal("encode error:", err)
			}

			_ = c.WriteMessage(websocket.TextMessage, bytes)
		}
	}
}

func IsConnected() {
	for _, v := range ClientHub {
		for {
			v.SocketConn.SetCloseHandler(func(code int, text string) error {
				fmt.Fprintf(os.Stderr, "websocket connection closed(%d, %s)\n", code, text)

				// from default CloseHandler
				message := websocket.FormatCloseMessage(code, "")
				_ = v.SocketConn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))

				v.IsAuthenticated = false

				return nil
			})
		}
	}
}
