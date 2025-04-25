package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/datarhei/gosrt"
)
//hp1a-55bj-6w4m-h8ef-8heb
const youtubeRTMP = "rtmp://a.rtmp.youtube.com/live2/hp1a-55bj-6w4m-h8ef-8heb"

// StreamProcessor handles the management of input streams (SRT), and the logic for sending streams to YouTube and Twitch.
type StreamProcessor struct {
	targets  []string
	f   string
	srtPortLow  int
	srtPortHigh int
}

// NewStreamProcessor creates and initializes a StreamProcessor
func NewStreamProcessor(targets []string, f string, srtPortLow, srtPortHigh int) *StreamProcessor {
	return &StreamProcessor{
		targets:  targets,
		f:   f,
		srtPortLow:  srtPortLow,
		srtPortHigh: srtPortHigh,
	}
}

// manageStreams accepts incoming SRT streams, manages them and sends them to both YouTube and Twitch
func (sp *StreamProcessor) manageStreams() error {
	// Creating an SRT listener for each SRT port in the range
	for port := sp.srtPortLow; port <= sp.srtPortHigh; port++ {
		// Create the SRT listener for each port
		srtListener, err := srt.Listen("srt",":" + strconv.Itoa(port), srt.DefaultConfig())
		if err != nil {
			log.Printf("Failed to listen on SRT port %d: %v", port, err)
			continue
		}

		log.Printf("Listening for incoming SRT streams on port %d", port)

		// Accept the incoming SRT connection
		go sp.handleSRTInputs(srtListener)
	}

	return nil
}

// handleSRTInputs processes incoming SRT connections and sends them to YouTube and Twitch
func (sp *StreamProcessor) handleSRTInputs(listener srt.Listener) {
	for {
		// Accept a new SRT connection
		conn, err := listener.Accept2()
		if err != nil {
			log.Printf("Failed to accept SRT connection: %v", err)
			return
		}
		log.Printf("New SRT connection accepted")
		fmt.Println(conn)

		// Here you would set up the stream to handle the video/audio streams (setupInputStreams)
		// and then send it to YouTube and Twitch (setupOutputStreams)

		// Start sending to YouTube and Twitch in separate goroutines (to handle both streams)
		//go sp.sendToYouTube(conn)
		//go sp.sendToTwitch(conn)
	}
}

// sendToYouTube sends the stream to YouTube, using the stream's URL and key
func (sp *StreamProcessor) sendToYouTube(conn srt.Conn) {
	// Handle sending the stream to YouTube using HLS
	log.Println("Sending stream to YouTube via HLS")
	// Step 2: Connect to YouTube RTMP server
	client, err := rtmp.NewClient(youtubeRTMP)
	if err != nil {
		log.Fatalf("Error connecting to YouTube RTMP: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected to YouTube RTMP")

	// Step 3: Stream data from SRT to RTMP
	for {
		data := make([]byte, 1316) // Standard TS packet size
		n, err := srt.Read(data)
		if err != nil {
			log.Printf("Error reading SRT stream: %v", err)
			//time.Sleep(1 * time.Second) // Prevent tight loop on errors
			continue
		}

		err = client.Write(data[:n])
		if err != nil {
			log.Printf("Error sending data to YouTube: %v", err)
			break
		}
	}
}

// sendToTwitch sends the stream to Twitch, transcoding from H.265 to H.264
func (sp *StreamProcessor) sendToTwitch(conn srt.Conn) {
	// Handle transcoding the stream and sending to Twitch via RTMP
	log.Println("Sending stream to Twitch via RTMP")
	// Setup and send the stream, ensure H.264 transcoding for Twitch compatibility
	// Call transcoding and streaming logic here
}

// main function to initialize the StreamProcessor and start managing streams
func main() {
	var targets [] string 
	f := "rtmp://twitch_url/live/"
	srtPortLow := 5000
	srtPortHigh := 5005

	// Create the stream processor
	sp := NewStreamProcessor(targets, f, srtPortLow, srtPortHigh)

	// Start managing SRT streams
	if err := sp.manageStreams(); err != nil {
		log.Fatalf("Error managing streams: %v", err)
	}

	// Keep the main function alive
	select {}
}

