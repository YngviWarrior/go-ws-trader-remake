package main

import (
	websocket "gowstrader/webSocket"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	var cache = memcache.New(os.Getenv("memcached"))
	var addr = os.Args[2]

	websocket.CreateWSServer(addr, cache)
}
