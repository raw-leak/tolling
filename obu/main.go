package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
	"tolling/types"

	"github.com/gorilla/websocket"
)

const (
	sendInterval = time.Second * 1 // 1 second
	wsEndpoint   = "ws://127.0.0.1:30000/ws"
)

func getCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()

	return n + f
}

func gemerateOBUIDs(max int) []int {
	ids := make([]int, max)

	for i := 0; i < max; i++ {
		ids[i] = rand.Intn(math.MaxInt)
	}

	return ids
}

func sendOBUData(conn *websocket.Conn, data *types.OBUData) error {
	return conn.WriteJSON(data)
}

func main() {
	// WS connection
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	// generate IDS
	obuIDs := gemerateOBUIDs(20)

	for {
		for i := 0; i < len(obuIDs); i++ {
			data := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   getCoord(),
				Long:  getCoord(),
			}

			if err := conn.WriteJSON(data); err != nil {
				log.Fatal("Error on sending DATA to WS", err)
			}

			fmt.Printf("%+v\n", data)

			time.Sleep(sendInterval)
		}

	}
}
