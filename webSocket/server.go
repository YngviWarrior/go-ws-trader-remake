package ws

import (
	"fmt"
	"gowstrader/mysql"
	trader "gowstrader/trader"
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func upgradeServer(cache *memcache.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade err: ", err)
			return
		}

		var client = trader.Client{}
		var dbConn = mysql.SqlConn{}
		dbConn.CreateConnection()

		// go trader.IsConnected(c)
		go trader.SubscribeUpdate(c, cache, &dbConn)
		go trader.Subscribe(c, cache, &client, &dbConn)
	}
}

func CreateWSServer(addr string, cache *memcache.Client) {
	var router = mux.NewRouter()

	router.HandleFunc("/", upgradeServer(cache))
	err := http.ListenAndServe(addr, router)

	fmt.Println(err)
}
