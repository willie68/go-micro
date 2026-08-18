package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/willie68/go-micro/internal/services/logging"
)

var loggerNames = []string{
	"api-service",
	"worker-pool",
	"database-sync",
	"cache-manager",
	"auth-handler",
}

func randomLogger() string {
	return loggerNames[rand.Intn(len(loggerNames))]
}

func main() {
	// Parse command line flags
	port := flag.Int("port", 12201, "GELF server port")
	host := flag.String("host", "127.0.0.1", "GELF server host")
	interval := flag.Int("interval", 5, "Interval in seconds between log messages")
	flag.Parse()

	// Initialize logger
	logging.Init(logging.Config{
		VictoriaLogsURL: "http://localhost:9428",
		Level:           "debug",
		GelfURL:         "",
	}, "gelfling-test")
	logger := logging.New("TestSender")

	fmt.Printf("TestSender started. Sending logs to %s:%d every %d seconds\n", *host, *port, *interval)
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create ticker for sending logs
	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	// Counter for message numbering
	counter := 1

	// Send initial log
	loggerName := randomLogger()
	logger.Info(fmt.Sprintf("Message #%d - TestSender started [%s]", counter, loggerName),
		"service", "TestSender",
		"logger", loggerName,
		"message_number", counter,
	)
	counter++

	// Main loop: send logs until signal received
	for {
		select {
		case <-sigChan:
			fmt.Println("\nShutdown signal received. Exiting...")
			loggerName := randomLogger()
			logger.Info(fmt.Sprintf("Message #%d - TestSender stopped [%s]", counter, loggerName),
				"service", "TestSender",
				"logger", loggerName,
				"message_number", counter,
			)
			time.Sleep(500 * time.Millisecond) // Allow last log to be sent
			return

		case <-ticker.C:
			loggerName := randomLogger()
			logger.Info(fmt.Sprintf("Message #%d - Test log entry [%s]", counter, loggerName),
				"service", "TestSender",
				"logger", loggerName,
				"message_number", counter,
				"timestamp", time.Now().Unix(),
			)
			counter++
		}
	}
}
