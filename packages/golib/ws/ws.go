package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

func ConnectToWs() *websocket.Conn{
	ws,_,err:=websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/btcusdt@ticker",nil);
	if err != nil{
		log.Fatal("failed to connect ", err);
	}
	return ws
}