package main

import (
	"encoding/json"
	"log"
	"net"

	pb "github.com/Inesh-Reddy/hft-3-gomicroservices/apps/go-services/ticker-service/proto/ticker"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// BinanceTicker matches the structure of the Binance WebSocket ticker data
type BinanceTicker struct {
	EventType string  `json:"e"` // Event type (e.g., "24hrTicker")
	E         int64   `json:"E"` // Event time
	S         string  `json:"s"` // Symbol
	C         *string `json:"c"` // Last price (optional, pointer to handle null)
	V         *string `json:"v"` // Volume (optional, pointer to handle null)
}

type TickerServer struct {
	pb.UnimplementedTickerServiceServer
}

func (s *TickerServer) StreamTicker(req *pb.TickerRequest, stream pb.TickerService_StreamTickerServer) error {
    log.Printf("Received StreamTicker request for symbol: %s", req.Symbol)
    
    // Ensure symbol is lowercase as Binance WebSocket expects lowercase symbols
    symbol := req.Symbol
    url := "wss://stream.binance.com:9443/ws/" + symbol + "@ticker"
    
    log.Printf("Connecting to Binance WebSocket: %s", url)
    c, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Printf("Failed to connect to WebSocket: %v", err)
        return err
    }
    defer c.Close()

    for {
        // Read message from WebSocket
        messageType, msg, err := c.ReadMessage()
        if err != nil {
            log.Printf("WebSocket read error: %v", err)
            return err
        }

        // Log the raw message for debugging
        log.Printf("Received WebSocket message (type: %d): %s", messageType, string(msg))

        // Skip non-JSON messages
        if messageType != websocket.TextMessage {
            log.Printf("Skipping non-text WebSocket message (type: %d)", messageType)
            continue
        }

        // Check if the message is valid JSON
        if !json.Valid(msg) {
            log.Printf("Invalid JSON received: %s", string(msg))
            continue
        }

        // Unmarshal the JSON message
        var t BinanceTicker
        if err := json.Unmarshal(msg, &t); err != nil {
            log.Printf("Failed to unmarshal JSON: %v, raw message: %s", err, string(msg))
            continue
        }

        // Verify the event type is "24hrTicker"
        if t.EventType != "24hrTicker" {
            log.Printf("Skipping non-ticker event: %s, raw message: %s", t.EventType, string(msg))
            continue
        }

        // Check for required fields
        if t.C == nil || t.V == nil {
            log.Printf("Missing price or volume in ticker data: %+v, raw message: %s", t, string(msg))
            continue
        }

        // Validate EventTime
        if t.E <= 0 || t.E > 4102444800000 { // Approx. year 2100 in milliseconds
            log.Printf("Invalid EventTime: %d, skipping message: %+v, raw message: %s", t.E, t, string(msg))
            continue
        }

        // Log the unmarshaled data
        log.Printf("Unmarshaled ticker data: %+v", t)

        // Create TickerUpdate for gRPC stream
        update := &pb.TickerUpdate{
            Exchange:  "binance",
            Symbol:    t.S,
            Price:     *t.C,
            Volume:    *t.V,
            EventTime: t.E,
        }

        // Send the update to the gRPC client
        if err := stream.Send(update); err != nil {
            log.Printf("Failed to send gRPC update: %v, update: %+v", err, update)
            return err
        }
    }
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Starting gRPC server for TickerService")
	
	grpcServer := grpc.NewServer()
	pb.RegisterTickerServiceServer(grpcServer, &TickerServer{})
	
	log.Println("🚀 Go Ticker Service running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}