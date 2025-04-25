// func mergeVideoWithBackgroundVideo(
// 	inputCtx *astiav.FormatContext,
// 	outputCtx *astiav.FormatContext,
// 	backgroundVideoPath string,
// 	mixRatio float64, // Ratio of background video (0 = only incoming, 1 = only background)
// ) error {
// 	if mixRatio < 0 || mixRatio > 1 {
// 		return fmt.Errorf("mixRatio must be between 0 and 1, got: %f", mixRatio)
// 	}

// 	// Open the background video file
// 	bgVideoCtx := astiav.AllocFormatContext()
// 	defer bgVideoCtx.Free()

// 	if err := bgVideoCtx.OpenInput(backgroundVideoPath, nil, nil); err != nil {
// 		log.Printf("Failed to open background video file: %v", err)
// 		return err
// 	}
// 	defer bgVideoCtx.CloseInput()

// 	// Find the video stream in the background video
// 	var bgVideoStream *astiav.Stream
// 	for _, stream := range bgVideoCtx.Streams() {
// 		if stream.CodecParameters().CodecType() == astiav.MediaTypeVideo {
// 			bgVideoStream = stream
// 			break
// 		}
// 	}
// 	if bgVideoStream == nil {
// 		return fmt.Errorf("no video stream found in background video file")
// 	}

// 	// Prepare the decoders for both the incoming video and the background video
// 	bgDecoder := astiav.AllocCodecContext(bgVideoStream.Codec())
// 	defer bgDecoder.Free()
// 	if err := bgDecoder.FromCodecParameters(bgVideoStream.CodecParameters()); err != nil {
// 		log.Printf("Failed to configure decoder for background video: %v", err)
// 		return err
// 	}
// 	if err := bgDecoder.Open(nil); err != nil {
// 		log.Printf("Failed to open background video decoder: %v", err)
// 		return err
// 	}

// 	// Prepare encoder for output video
// 	videoEncoder := astiav.AllocCodecContext(astiav.FindEncoder(astiav.CodecIDH264))
// 	defer videoEncoder.Free()
// 	videoEncoder.SetWidth(bgDecoder.Width())
// 	videoEncoder.SetHeight(bgDecoder.Height())
// 	videoEncoder.SetPixelFormat(bgDecoder.PixelFormat())
// 	videoEncoder.SetTimeBase(astiav.NewRational(1, 30)) // Assuming 30 FPS; adjust as needed
// 	if err := videoEncoder.Open(nil); err != nil {
// 		log.Printf("Failed to open video encoder: %v", err)
// 		return err
// 	}

// 	// Create new video stream in the output context
// 	videoStream := outputCtx.NewStream(nil)
// 	if err := videoStream.CodecParameters().FromCodecContext(videoEncoder); err != nil {
// 		log.Printf("Failed to configure output video stream: %v", err)
// 		return err
// 	}

// 	// Set up video mixing
// 	packet := astiav.AllocPacket()
// 	defer packet.Free()
// 	frame := astiav.AllocFrame()
// 	defer frame.Free()
// 	mixedFrame := astiav.AllocFrame()
// 	defer mixedFrame.Free()

// 	for {
// 		// Read from the background video file and loop if necessary
// 		if err := bgVideoCtx.ReadFrame(packet); err != nil {
// 			if errors.Is(err, astiav.ErrEof) {
// 				// Loop the background video
// 				bgVideoCtx.SeekFrame(0, astiav.SeekFlagAny)
// 				continue
// 			}
// 			log.Printf("Error reading background video frame: %v", err)
// 			break
// 		}

// 		// Process background video frames
// 		if packet.StreamIndex() == bgVideoStream.Index() {
// 			if err := bgDecoder.SendPacket(packet); err != nil {
// 				log.Printf("Error sending background video packet to decoder: %v", err)
// 				break
// 			}

// 			for {
// 				if err := bgDecoder.ReceiveFrame(frame); err != nil {
// 					if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 						break
// 					}
// 					log.Printf("Error receiving decoded background video frame: %v", err)
// 					break
// 				}

// 				// Here we would fetch the corresponding input video frame to mix
// 				inputFrame := astiav.AllocFrame()
// 				defer inputFrame.Free()
// 				// Assume inputFrame is retrieved and matches frame attributes

// 				// Mix the frames based on mixRatio
// 				if err := mixVideoFrames(frame, inputFrame, mixedFrame, mixRatio); err != nil {
// 					log.Printf("Error mixing frames: %v", err)
// 					break
// 				}

// 				// Send the mixed frame to the encoder
// 				if err := videoEncoder.SendFrame(mixedFrame); err != nil {
// 					log.Printf("Error sending mixed frame to encoder: %v", err)
// 					break
// 				}

// 				// Write encoded packets to output
// 				for {
// 					if err := videoEncoder.ReceivePacket(packet); err != nil {
// 						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 							break
// 						}
// 						log.Printf("Error receiving encoded video packet: %v", err)
// 						break
// 					}

// 					packet.SetStreamIndex(videoStream.Index())
// 					if err := outputCtx.WriteInterleavedFrame(packet); err != nil {
// 						log.Printf("Error writing video packet to output: %v", err)
// 						break
// 					}
// 					packet.Unref()
// 				}
// 			}
// 		}
// 		packet.Unref()
// 	}

// 	return nil
// }

// // mixVideoFrames blends two video frames together using the specified mix ratio
// func mixVideoFrames(background *astiav.Frame, input *astiav.Frame, mixed *astiav.Frame, mixRatio float64) error {
// 	// Ensure the frame data matches
// 	if background.Width() != input.Width() || background.Height() != input.Height() || background.Format() != input.Format() {
// 		return fmt.Errorf("frame attributes do not match for mixing")
// 	}

// 	// Allocate memory for the mixed frame
// 	mixed.AllocData()
// 	mixed.SetWidth(background.Width())
// 	mixed.SetHeight(background.Height())
// 	mixed.SetFormat(background.Format())

// 	// Mix pixel data for blending
// 	for i := 0; i < len(background.Data(0)); i++ {
// 		mixed.Data(0)[i] = byte(
// 			float64(background.Data(0)[i])*mixRatio +
// 				float64(input.Data(0)[i])*(1-mixRatio),
// 		)
// 	}

// 	return nil
// }
