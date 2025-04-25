// package main

// import (
// 	"fmt"

// 	"github.com/asticode/go-astiav"
// )

// // applyChatOverlay applies a chat overlay stream onto the main video stream.
// func applyChatOverlay(mainInputCtx, chatInputCtx, outputCtx *astiav.FormatContext) error {
// 	// Initialize filter graph
// 	filterGraph := astiav.AllocFilterGraph()
// 	if filterGraph == nil {
// 		return fmt.Errorf("failed to allocate filter graph")
// 	}
// 	defer filterGraph.Free()

// 	// Define filter chain: overlay chat on main video
// 	mainInput := "[main]"
// 	chatInput := "[chat]"
// 	output := "[out]"
// 	overlayFilter := fmt.Sprintf("%s%s%s", mainInput, chatInput, output)

// 	// Add inputs for main and chat streams
// 	mainInputs := astiav.AllocFilterInOut()
// 	if mainInputs == nil {
// 		return fmt.Errorf("failed to allocate filter in/out for main")
// 	}
// 	defer mainInputs.Free()

// 	chatInputs := astiav.AllocFilterInOut()
// 	if chatInputs == nil {
// 		return fmt.Errorf("failed to allocate filter in/out for chat")
// 	}
// 	defer chatInputs.Free()

// 	// Configure the overlay filter
// 	if err := filterGraph.Parse(overlayFilter, nil); err != nil {
// 		return fmt.Errorf("failed to parse filter graph: %w", err)
// 	}

// 	// Configure filter graph
// 	if err := filterGraph.Configure(); err != nil {
// 		return fmt.Errorf("failed to configure filter graph: %w", err)
// 	}

// 	// Process frames
// 	for {
// 		// Read main video frame
// 		mainFrame := astiav.AllocFrame()
// 		if mainFrame == nil {
// 			return fmt.Errorf("failed to allocate frame for main input")
// 		}
// 		defer mainFrame.Free()

// 		if err := mainInputCtx.ReadFrame(mainFrame); err != nil {
// 			break
// 		}

// 		// Read chat overlay frame
// 		chatFrame := astiav.AllocFrame()
// 		if chatFrame == nil {
// 			return fmt.Errorf("failed to allocate frame for chat input")
// 		}
// 		defer chatFrame.Free()

// 		if err := chatInputCtx.ReadFrame(chatFrame); err != nil {
// 			break
// 		}

// 		// Send frames to filter graph
// 		if err := filterGraph.SendFrame(mainFrame); err != nil {
// 			return fmt.Errorf("failed to send main frame to filter graph: %w", err)
// 		}
// 		if err := filterGraph.SendFrame(chatFrame); err != nil {
// 			return fmt.Errorf("failed to send chat frame to filter graph: %w", err)
// 		}

// 		// Retrieve processed frame
// 		overlayFrame := astiav.AllocFrame()
// 		if overlayFrame == nil {
// 			return fmt.Errorf("failed to allocate overlay frame")
// 		}
// 		defer overlayFrame.Free()

// 		if err := filterGraph.ReceiveFrame(overlayFrame); err == nil {
// 			// Write processed frame to output
// 			if err := outputCtx.WriteFrame(overlayFrame); err != nil {
// 				return fmt.Errorf("failed to write overlay frame to output: %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }
