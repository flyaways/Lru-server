package main

import (
	"flag"
	"fmt"
	"time"

	lru "github.com/flyaways/golang-lru"
	"github.com/flyaways/golang-lru/simplelru"
	server "github.com/flyaways/lru-server"
	"github.com/gin-gonic/gin"
)

var (
	size    = flag.Int("s", 10000, "-s 10000")
	port    = flag.String("p", ":8080", "-p :8080")
	mode    = flag.String("m", "debug", "-m debug/release")
	cPolicy = flag.Int("t", 1, "-t 0/1/2/3")
	cache   simplelru.LRUCache
)

const (
	sLRU = iota
	safeLRU
	twoQ
	tARC
)

func NewCache(cPolicy, size int) (simplelru.LRUCache, error) {
	switch cPolicy {
	case sLRU:
		return simplelru.NewLRU(8, func(key interface{}, value interface{}) {
			fmt.Println(time.Now().Format(time.RFC3339Nano))
		})

	case safeLRU:
		/* return lru.NewWithEvict(8, func(key interface{}, value interface{}) {
			fmt.Println(time.Now().Format(time.RFC3339Nano))
		}) */
		return lru.New(size)

	case twoQ:
		return lru.New2Q(size)

	case tARC:
		return lru.NewARC(size)

	default:
		panic(cPolicy)
	}
}

func main() {
	flag.Parsed()

	cache, err := NewCache(*cPolicy, *size)
	if err != nil {
		panic(err)
	}

	gin.SetMode(*mode)
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	server.Version1(engine.Group("/api/v1"), cache)

	engine.Run(*port)
}
