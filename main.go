// zoe - 2026
// this server was made for a very specfic usecase feel free too use and alter
package main

import (
	"fmt"
	"io"
	"ivygo/poisonIvy"
	"log"
	"net/http"
	"os"
	"time"
)

var printerClient *poisonIvy.Client
var isPrinterConnected bool

func printHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	//  if printer isn't connected tell the iPad immediately
	if !isPrinterConnected || printerClient == nil {
		http.Error(w, "Printer is currently offline. Please wait for reconnection.", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	tempFile := "incoming_print.png"
	f, err := os.Create(tempFile)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(f, r.Body)
	f.Close()

	imgData, err := poisonIvy.PrepareImage(tempFile)
	if err != nil {
		log.Printf("Image prep error: %v", err)
		http.Error(w, "Invalid Image Format", http.StatusBadRequest)
		return
	}

	fmt.Printf("Printing image: %d bytes processed\n", len(imgData))

	printerClient.OutboundQ <- poisonIvy.GetBaseMessage(0, true, false)
	time.Sleep(500 * time.Millisecond)

	printerClient.KeepAlive() // boop the sensors (this fixes a weird bug where fails too print after we reset paper)
	time.Sleep(300 * time.Millisecond)

	printerClient.OutboundQ <- poisonIvy.GetPrintReadyMessage(len(imgData))
	time.Sleep(500 * time.Millisecond)

	chunkSize := 990
	for i := 0; i < len(imgData); i += chunkSize {
		end := i + chunkSize
		if end > len(imgData) {
			end = len(imgData)
		}
		printerClient.OutboundQ <- imgData[i:end]
	}

}

// init connetcion
func connectToPrinter(port string) {
	for {
		if !isPrinterConnected {
			fmt.Printf("Attempting to connect to printer on %s...\n", port)

			printerClient = poisonIvy.NewClient()
			err := printerClient.Connect(port)

			if err != nil {
				fmt.Printf("Connection failed: %v. Retrying in 5 seconds...\n", err)
				isPrinterConnected = false
			} else {
				fmt.Println("Successfully connected! Configuring printer...")
				isPrinterConnected = true

				// Disable Sleep
				disableSleep := poisonIvy.GetBaseMessage(259, false, false)
				disableSleep[8] = 0
				printerClient.OutboundQ <- disableSleep

			}
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	const printerPort = "COM8"
	const webPort = "105"

	go connectToPrinter(printerPort)

	go func() {
		ticker := time.NewTicker(45 * time.Second)
		for range ticker.C {
			if isPrinterConnected && printerClient != nil {
				printerClient.KeepAlive()
			}
		}
	}()

	fmt.Printf("--- IVY PRINT SERVER ONLINE ---\n")
	fmt.Printf("Listening on http://0.0.0.0:%s/print\n", webPort)

	http.HandleFunc("/print", printHandler)

	err := http.ListenAndServe(":"+webPort, nil)
	if err != nil {
		fmt.Printf("FATAL: Could not start web server: %v\n", err)
		fmt.Println("Is another instance of the server already running?")
		os.Exit(1)
	}
}
