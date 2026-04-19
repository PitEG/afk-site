package main

import (
	"fmt"
	"time"
	"database/sql"
	"math/big"
	_ "modernc.org/sqlite"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
)

func main() {
	r := gin.Default()

	store := persistence.NewInMemoryStore(time.Second)
	score := big.NewInt(0)
	startTime := time.Now().Unix()

	r.GET("/time", cache.CachePageAtomic(store, time.Second, func(c *gin.Context) {
		c.String(200, "score: "+score.String())
  }))

	go incrementScore(score,startTime)

	r.Run(":8080")

	fmt.Println("hi")
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
