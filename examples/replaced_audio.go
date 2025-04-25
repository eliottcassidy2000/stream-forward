// import (
// 	"log"
// 	"os"

// 	"github.com/asticode/go-astiav"
// )

// // replaceAudioWithBackgroundMusic replaces the audio in the incoming stream with a looped background music file
// func replaceAudioWithBackgroundMusic(inputCtx *astiav.FormatContext, outputCtx *astiav.FormatContext, musicFilePath string) error {
// 	// Open the background music file
// 	musicCtx := astiav.AllocFormatContext()
// 	defer musicCtx.Free()

// 	if err := musicCtx.OpenInput(musicFilePath, nil, nil); err != nil {
// 		log.Printf("Failed to open background music file: %v", err)
// 		return err
// 	}
// 	defer musicCtx.CloseInput()

// 	// Find the audio stream in the music file
// 	var musicAudioStream *astiav.Stream
// 	for _, stream := range musicCtx.Streams() {
// 		if stream.CodecParameters().CodecType() == astiav.MediaTypeAudio {
// 			musicAudioStream = stream
// 			break
// 		}
// 	}
// 	if musicAudioStream == nil {
// 		return errors.New("no audio stream found in background music file")
// 	}

// 	// Prepare the decoder for the background music
// 	musicDecoder := astiav.AllocCodecContext(musicAudioStream.Codec())
// 	defer musicDecoder.Free()
// 	if err := musicDecoder.FromCodecParameters(musicAudioStream.CodecParameters()); err != nil {
// 		log.Printf("Failed to configure decoder for background music: %v", err)
// 		return err
// 	}
// 	if err := musicDecoder.Open(nil); err != nil {
// 		log.Printf("Failed to open music decoder: %v", err)
// 		return err
// 	}

// 	// Prepare the encoder for the output audio stream
// 	audioEncoder := astiav.AllocCodecContext(astiav.FindEncoder(astiav.CodecIDAAC))
// 	defer audioEncoder.Free()
// 	audioEncoder.SetSampleRate(musicDecoder.SampleRate())
// 	audioEncoder.SetChannelLayout(musicDecoder.ChannelLayout())
// 	audioEncoder.SetChannels(musicDecoder.Channels())
// 	audioEncoder.SetSampleFormat(astiav.SampleFormatFLTP)
// 	audioEncoder.SetTimeBase(astiav.NewRational(1, musicDecoder.SampleRate()))
// 	if err := audioEncoder.Open(nil); err != nil {
// 		log.Printf("Failed to open audio encoder: %v", err)
// 		return err
// 	}

// 	// Create a new audio stream in the output context
// 	audioStream := outputCtx.NewStream(nil)
// 	if err := audioStream.CodecParameters().FromCodecContext(audioEncoder); err != nil {
// 		log.Printf("Failed to configure output audio stream: %v", err)
// 		return err
// 	}

// 	// Loop through the music file and replace the audio in the input stream
// 	packet := astiav.AllocPacket()
// 	defer packet.Free()
// 	frame := astiav.AllocFrame()
// 	defer frame.Free()

// 	for {
// 		if err := musicCtx.ReadFrame(packet); err != nil {
// 			if errors.Is(err, astiav.ErrEof) {
// 				// Loop the background music
// 				musicCtx.SeekFrame(0, astiav.SeekFlagAny)
// 				continue
// 			}
// 			log.Printf("Error reading background music frame: %v", err)
// 			break
// 		}

// 		if packet.StreamIndex() == musicAudioStream.Index() {
// 			if err := musicDecoder.SendPacket(packet); err != nil {
// 				log.Printf("Error sending music packet to decoder: %v", err)
// 				break
// 			}

// 			for {
// 				if err := musicDecoder.ReceiveFrame(frame); err != nil {
// 					if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 						break
// 					}
// 					log.Printf("Error receiving decoded music frame: %v", err)
// 					break
// 				}

// 				// Send the decoded frame to the encoder
// 				if err := audioEncoder.SendFrame(frame); err != nil {
// 					log.Printf("Error sending music frame to encoder: %v", err)
// 					break
// 				}

// 				// Receive the encoded packets and write them to the output
// 				for {
// 					if err := audioEncoder.ReceivePacket(packet); err != nil {
// 						if errors.Is(err, astiav.ErrEof) || errors.Is(err, astiav.ErrEagain) {
// 							break
// 						}
// 						log.Printf("Error receiving encoded audio packet: %v", err)
// 						break
// 					}

// 					// Write the packet to the output context
// 					packet.SetStreamIndex(audioStream.Index())
// 					if err := outputCtx.WriteInterleavedFrame(packet); err != nil {
// 						log.Printf("Error writing audio packet to output: %v", err)
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
