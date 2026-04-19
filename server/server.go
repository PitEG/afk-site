package main

import (
	"fmt"
	"log"
	"time"
	"database/sql"
	"math/big"
	_ "modernc.org/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func main() {
	r := setupRouter()

	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	store := persistence.NewInMemoryStore(time.Second)
	score := big.NewInt(0)
	startTime := time.Now().Unix()

	go incrementScore(score,startTime)

	r.GET("/time", cache.CachePageAtomic(store, time.Second, func(c *gin.Context) {
		c.String(200, "score: "+score.String())
  }))
	r.GET("/ping", handlePing)

	return gin.Default()
}

func handlePing(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	conn, err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Println("upgrade error: ", err)
		return
	}
	defer func() { _ = conn.Close() }()

	conn.SetReadDeadline(time.Now().Add(5*time.Second))
	/*
	conn.SetPongHandler(func(string) error { 
		conn.SetReadDeadline(time.Now().Add(5*time.Second))
		return nil 
	})
	*/

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("failed reading: ",err)
			break
		}
		conn.SetReadDeadline(time.Now().Add(5*time.Second))
		log.Println("message received:", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}

func incrementScore(score *big.Int, startTime int64) {
	for {
		time.Sleep(time.Second)
		score.Sub(big.NewInt(time.Now().Unix()),big.NewInt(startTime))
	}
}

func test_db() {
	db, err := sql.Open("sqlite", "thing.db")
	if err != nil {
		fmt.Println("error")
	}
	defer db.Close()

	_, err = db.Exec(`create table dada(id integer not null primary key, name text);`)

	if err != nil {
		fmt.Println(err)
	}

	db.Exec(`insert into dada(id,name) values (10,"andrew");`)
}
