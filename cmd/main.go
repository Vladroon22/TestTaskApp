package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Vladroon22/multiplicator/internal/handlers"
	"github.com/Vladroon22/multiplicator/internal/server"
)

func parseRTP() (float64, error) {
	if len(os.Args) < 2 {
		return 0, fmt.Errorf("value is not found")
	}

	for _, arg := range os.Args[1:] {
		if len(arg) > 4 && arg[:4] == "-rtp" {
			numberPart := arg[4:]

			if numberPart == "" {
				return 0, fmt.Errorf("after -rtp must be a number")
			}

			value, err := strconv.ParseFloat(numberPart, 64)
			if err != nil {
				return 0, fmt.Errorf("wrong number's format: %s", numberPart)
			}

			if value <= 0 || value > 1.0 {
				return 0, fmt.Errorf("rtp must be in (0; 1.0]") // 0.000 - how many after ','
			}

			return value, nil
		}
	}

	return 0, fmt.Errorf("flag -rtp{value} not found")
}

func main() {
	rtp, err := parseRTP()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Example: go run . -rtp0.5")
		os.Exit(1)
	}

	router := http.NewServeMux()
	router.HandleFunc("/get", handlers.Multiplicate(rtp))

	port := ":64333"
	srv := server.NewServer(port, router)

	stopServerChan := make(chan error, 1)

	log.Printf("HTTP-server runs on %s with rtp=%f", port, rtp)
	go func() {
		if err := srv.Start(); err != nil {
			stopServerChan <- err
		}
	}()

	stopOSChan := make(chan os.Signal, 1)
	signal.Notify(stopOSChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-stopServerChan:
		log.Fatalf("HTTP-server error: %v", err)
		close(stopServerChan)
	case sig := <-stopOSChan:
		log.Printf("Received %v signal, shutting down...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP-server shutdown error: %v", err)
		}
	}
}
