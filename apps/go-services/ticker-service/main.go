package main

import (
	"encoding/json"
	"log"
	"net"

	pb "github.com/Inesh-Reddy/hft-3-gomicroservices/apps/go-services/ticker-service/proto/ticker"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)
type TickerServer struct {
	pb.UnimplementedTickerServiceServer
}

type BinanceTicker struct {
	E int64		`json:"E"`
	S string	`json:"s"`
	C string	`json:"c"`
	V string	`json:"v"`
}

func (s *TickerServer) StreamTicker(req *pb.TickerRequest, stream pb.TickerService_StreamTickerServer) error {
	url := "wss://stream.binance.com:9443/ws/" + req.Symbol + "@ticker"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return err
		}

		var t BinanceTicker
		if err := json.Unmarshal(msg, &t); err != nil {
			continue
		}

		update := &pb.TickerUpdate{
			Exchange:  "binance",
			Symbol:    t.S,
			Price:     t.C,
			Volume:    t.V,
			EventTime: t.E,
		}

		if err := stream.Send(update); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTickerServiceServer(grpcServer, &TickerServer{})
	log.Println("🚀 Go Ticker Service running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// func main(){
// 	fmt.Println(`Ticker Service running......`)

// 	ctx := context.Background()
// 	log.Println(`Context:`, ctx);
// 	redis := redis.ConnectToRedis()
// 	redis.Ping(ctx)
// 	wes:=ws.ConnectToWs()
// 	wes.PingHandler();
// 	_,data,err:=wes.ReadMessage()
// 	if err != nil {
// 		log.Fatal("erroredading data from ws:", err)
// 	}
// 	log.Println(string(data))


// }