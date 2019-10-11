package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/siddontang/go-log/log"
	"hash/crc32"
	"net/url"
	"strings"
	"sync"
	"time"
)

type subRequest struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type cbMessage struct {
	Type          string `json:"type"`
	Side          string `json:"side"`
	ProductId     string `json:"product_id"`
	OrderId       string `json:"order_id"`
	OrderType     string `json:"order_type"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	Funds         string `json:"funds"`
	RemainingSize string `json:"remaining_size"`
}

var tokens sync.Map

func getTokenByProductId(productId string) string {
	token, found := tokens.Load(productId)
	if found {
		return token.(string)
	}

	email := strings.ToLower(productId) + "@trader.com"
	password := "123456"
	t, err := getToken(email, password)
	if err != nil {
		panic(err)
	}
	t = strings.Replace(t, "\"", "", 2)
	tokens.Store(productId, t)
	return t
}

func coinbaseWs() {
	//cancelOrders("BTC-USDT", "buy")
	//cancelOrders("BTC-USDT", "sell")
	//time.Sleep(time.Second)

	u := url.URL{Scheme: "wss", Host: "ws-feed.pro.coinbase.com", Path: "", RawQuery: ""}
	d := websocket.Dialer{TLSClientConfig: &tls.Config{RootCAs: nil, InsecureSkipVerify: true}, HandshakeTimeout: 10 * time.Second}
	c, _, err := d.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}

	req := subRequest{
		Type:       "subscribe",
		ProductIds: []string{"BTC-USD", "LTC-USD", "ETH-USD", "BCH-USD", "EOS-USD"},
		Channels:   []string{"full"},
	}
	buf, _ := json.Marshal(req)

	err = c.WriteMessage(websocket.TextMessage, buf)
	if err != nil {
		log.Fatalf("write error: %v", err)
		_ = c.Close()
	}

	workers := newWorkers(10)

	defer func() {
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Error("write close error: %v", err)
		}
		_ = c.Close()
	}()
	for {
		select {
		default:
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Error("read error: %v", err)
				return
			}
			//log.Info(string(message))

			if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
				err := c.WriteMessage(websocket.PongMessage, []byte(time.Now().String()))
				if err != nil {
					log.Error("wsWrite", err)
					return
				}
			}

			var msg cbMessage
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Error(err)
				continue
			}

			if strings.HasSuffix(msg.ProductId, "USD") {
				msg.ProductId = strings.Replace(msg.ProductId, "USD", "USDT", 1)
			}

			hash := hashCode(msg.OrderId)
			workers[hash%10].msgCh <- &msg
		}
	}
}

type worker struct {
	msgCh chan *cbMessage
}

func newWorkers(n int) (workers []*worker) {
	for i := 0; i < n; i++ {
		workers = append(workers, newWorker())
	}
	return workers
}

func newWorker() *worker {
	w := &worker{msgCh: make(chan *cbMessage, 1000)}
	go func() {
		orderMap := map[string]*gbeOrder{}
		for {
			select {
			case msg := <-w.msgCh:
				switch msg.Type {
				case "received":
					if msg.OrderType == "market" {
						_, err := placeOrder(getTokenByProductId(msg.ProductId), msg.ProductId, msg.Size, msg.Price, msg.Funds, msg.Side, "market")
						if err != nil {
							log.Error(err)
							continue
						}
					} else {
						o, err := placeOrder(getTokenByProductId(msg.ProductId), msg.ProductId, msg.Size, msg.Price, msg.Funds, msg.Side, "limit")
						if err != nil {
							log.Error(err)
							continue
						}
						orderMap[msg.OrderId] = o
					}
				case "done":
					o, found := orderMap[msg.OrderId]
					if !found {
						continue
					}
					err := cancelOrder(getTokenByProductId(msg.ProductId), o.Id)
					if err != nil {
						log.Error(err)
					}
					delete(orderMap, msg.OrderId)
				}
			}
		}
	}()
	return w
}

func hashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}
