// package main

// import (
// 	"fmt"
// 	"github.com/asticode/go-astiav"
// )

// func applyTranslucentOverlay(inputCtx *astiav.FormatContext, outputCtx *astiav.FormatContext, color string, alpha float64) error {
// 	// Initialize filter graph
// 	filterGraph := astiav.AllocFilterGraph()
// 	if filterGraph == nil {
// 		return fmt.Errorf("failed to allocate filter graph")
// 	}
// 	defer filterGraph.Free()

// 	// Create filter descriptions
// 	videoInput := "[in]"
// 	colorFilter := fmt.Sprintf("color=color=%s:size=%dx%d:duration=0", color, 1280, 720) // Adjust size as needed
// 	blendFilter := fmt.Sprintf("blend=all_mode=overlay:all_opacity=%f", alpha)
// 	videoOutput := "[out]"

// 	// Define the filter chain
// 	filterChain := fmt.Sprintf("%s,%s,%s", videoInput, colorFilter, blendFilter)

// 	// Add filters to the graph
// 	if err := filterGraph.Parse(filterChain, nil); err != nil {
// 		return fmt.Errorf("failed to parse filter chain: %w", err)
// 	}

// 	// Configure filter graph
// 	if err := filterGraph.Configure(); err != nil {
// 		return fmt.Errorf("failed to configure filter graph: %w", err)
// 	}

// 	// Apply the filter graph to the video frames
// 	for {
// 		// Fetch frame from input
// 		frame := astiav.AllocFrame()
// 		if frame == nil {
// 			return fmt.Errorf("failed to allocate frame")
// 		}
// 		defer frame.Free()

// 		if err := inputCtx.ReadFrame(frame); err != nil {
// 			break
// 		}

// 		// Process frame through filter graph
// 		if err := filterGraph.SendFrame(frame); err != nil {
// 			return fmt.Errorf("failed to send frame to filter graph: %w", err)
// 		}

// 		// Retrieve processed frame
// 		filteredFrame := astiav.AllocFrame()
// 		if filteredFrame == nil {
// 			return fmt.Errorf("failed to allocate filtered frame")
// 		}
// 		defer filteredFrame.Free()

// 		if err := filterGraph.ReceiveFrame(filteredFrame); err == nil {
// 			// Write frame to output
// 			if err := outputCtx.WriteFrame(filteredFrame); err != nil {
// 				return fmt.Errorf("failed to write frame to output: %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }
