package main

import (
	"github.com/siddontang/go-log/log"
	_ "github.com/siddontang/go-mysql/mysql"
	"time"
)

//const gitBitExAddr = "https://gitbitex.com:8080"

const gitBitExAddr = "http://127.0.0.1:8001"

func main() {
	for {
		log.Info("connecting...")
		coinbaseWs()
		time.Sleep(time.Second)
	}
}
