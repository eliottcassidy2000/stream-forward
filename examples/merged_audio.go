func mergeAudioWithBackgroundMusic(
	inputCtx *astiav.FormatContext,
	outputCtx *astiav.FormatContext,
	musicFilePath string,
	mixRatio float64, // Ratio of background music (0 = only incoming, 1 = only music)
) error {
	if mixRatio < 0 || mixRatio > 1 {
		return fmt.Errorf("mixRatio must be between 0 and 1, got: %f", mixRatio)
	}

	// Open the background music file
	musicCtx := astiav.AllocFormatContext()
	defer musicCtx.Free()

	if err := musicCtx.OpenInput(musicFilePath, nil, nil); err != nil {
		log.Printf("Failed to open background music file: %v", err)
		return err
	}
	defer musicCtx.CloseInput()

	// Find the audio stream in the music file
	var musicAudioStream *astiav.Stream
	for _, stream := range musicCtx.Streams() {
		if stream.CodecParameters().CodecType() == astiav.MediaTypeAudio {
			musicAudioStream = stream
			break
		}
	}
	if musicAudioStream == nil {
		return fmt.Errorf("no audio stream found in background music file")
	}

	// Prepare the decoders for both the incoming audio and the background music
	musicDecoder := astiav.AllocCodecContext(musicAudioStream.Codec())
	defer musicDecoder.Free()
	if err := musicDecoder.FromCodecParameters(musicAudioStream.CodecParameters()); err != nil {
		log.Printf("Failed to configure decoder for background music: %v", err)
		return err
	}
	if err := musicDecoder.Open(nil); err != nil {
		log.Printf("Failed to open music decoder: %v", err)
		return err
	}

	// Prepare encoder for output audio
	audioEncoder := astiav.AllocCodecContext(astiav.FindEncoder(astiav.CodecIDAAC))
	defer audioEncoder.Free()
	audioEncoder.SetSampleRate(musicDecoder.SampleRate())
	audioEncoder.SetChannelLayout(musicDecoder.ChannelLayout())
	audioEncoder.SetChannels(musicDecoder.Channels())
	audioEncoder.SetSampleFormat(astiav.SampleFormatFLTP)
	audioEncoder.SetTimeBase(astiav.NewRational(1, musicDecoder.SampleRate()))
	if err := audioEncoder.Open(nil); err != nil {
		log.Printf("Failed to open audio encoder: %v", err)
		return err
	}

	// Create new audio stream in the output context
	audioStream := outputCtx.NewStream(nil)
	if err := audioStream.CodecParameters().FromCodecContext(audioEncoder); err != nil {
		log.Printf("Failed to configure output audio stream: %v", err)
		return err
	}

	// Set up audio mixing
	packet := astiav.AllocPacket()
	defer packet.Free()
	frame := astiav.AllocFrame()
	defer frame.Free()
	mixedFrame := astiav.AllocFrame()
	defer mixedFrame.Free()

	for {
		// Read from the music file and loop if necessary
		if err := musicCtx.ReadFrame(packet); err != nil {
			if errors.Is(err, astiav.ErrEof) {
				// Loop the background music
				musicCtx.SeekFrame(0, astiav.SeekFlagAny)
				continue
			}
			log.Printf("Error reading background music frame: %v", err)
			break
		}

		// Process music frames
		if packet.StreamIndex() == musicAudioStream.Index() {
			if err := musicDecoder.SendPacket(packet); err != nil {
				log.Printf("Error sending music packet to decoder: %v", err)
				break
			}

			for {
				if err := musicDecoder.ReceiveFrame(frame); err != nil {
					if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
						break
					}
					log.Printf("Error receiving decoded music frame: %v", err)
					break
				}

				// Here we would fetch the corresponding input audio frame to mix
				inputFrame := astiav.AllocFrame()
				defer inputFrame.Free()
				// Assume inputFrame is retrieved and matches frame attributes

				// Mix the frames based on mixRatio
				if err := mixFrames(frame, inputFrame, mixedFrame, mixRatio); err != nil {
					log.Printf("Error mixing frames: %v", err)
					break
				}

				// Send the mixed frame to the encoder
				if err := audioEncoder.SendFrame(mixedFrame); err != nil {
					log.Printf("Error sending mixed frame to encoder: %v", err)
					break
				}

				// Write encoded packets to output
				for {
					if err := audioEncoder.ReceivePacket(packet); err != nil {
						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
							break
						}
						log.Printf("Error receiving encoded audio packet: %v", err)
						break
					}

					packet.SetStreamIndex(audioStream.Index())
					if err := outputCtx.WriteInterleavedFrame(packet); err != nil {
						log.Printf("Error writing audio packet to output: %v", err)
						break
					}
					packet.Unref()
				}
			}
		}
		packet.Unref()
	}

	return nil
}

// mixFrames combines two audio frames with a specified mix ratio
func mixFrames(background *astiav.Frame, input *astiav.Frame, mixed *astiav.Frame, mixRatio float64) error {
	// Ensure the frame data matches
	if background.NbSamples() != input.NbSamples() || background.Format() != input.Format() || background.Channels() != input.Channels() {
		return fmt.Errorf("frame attributes do not match for mixing")
	}

	// Allocate memory for the mixed frame
	mixed.AllocData()
	mixed.SetFormat(background.Format())
	mixed.SetNbSamples(background.NbSamples())
	mixed.SetChannels(background.Channels())

	// Mix samples
	for i := 0; i < len(background.Data(0)); i++ {
		mixed.Data(0)[i] = byte(
			float64(background.Data(0)[i])*mixRatio +
				float64(input.Data(0)[i])*(1-mixRatio),
		)
	}

	return nil
}
