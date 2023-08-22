package main

import (
	"fmt"
	"math"

	"github.com/xthexder/go-jack"
)

const (
	channels     = 2     // Number of channels to process
	PCMMax       = 1.0   // Maximum float value for PCM audio
	ReferenceFS  = -26.0 // Reference value for FS
	ReferenceSPL = 94.0  // Reference value for SPL
	BufferSize   = 1024  // Number of samples to calculate dB SPL
)

// Global declarations for JACK audio ports and client
var (
	PortsIn     []*jack.Port
	PortsOut    []*jack.Port
	JackClient  *jack.Client
	isConnected bool // Flag to check if the client is connected
)

func computeRMS(samples []jack.AudioSample) float64 {
	var sum float64
	for _, sample := range samples {
		sum += float64(sample * sample)
	}
	return math.Sqrt(sum / float64(len(samples)))
}

func rmsToDBFS(rms float64) float64 {
	if rms == 0.0 {
		return -math.Inf(0) // or some other predefined value for silence
	}
	return 20.0 * math.Log10(rms/PCMMax)
}

func dBFS_to_dBSPL(dbfs float64) float64 {
	return (dbfs - ReferenceFS) + ReferenceSPL
}

var samples []uint32

// process handles the audio data for each frame.
// For now, it prints the non-zero samples from the first channel.
func process(nframes uint32) int {
	port := PortsIn[0]

	// Check if the client is connected before processing
	if isConnected {

		samplesIn := port.GetBuffer(nframes)
		// Compute and print if samplesIn size meets or exceeds BufferSize
		if len(samplesIn) >= BufferSize {
			rms := computeRMS(samplesIn)
			dbFS := rmsToDBFS(rms)
			dbSPL := dBFS_to_dBSPL(dbFS)

			fmt.Printf("RMS: %f, dB FS: %f, dB SPL: %f\n", rms, dbFS, dbSPL)
		}
	}
	return 0
}

// processXX is a callback for when port connections change.
func processXX(x jack.PortId, y jack.PortId, z bool) {
	isConnected = z // Use z to determine connection status
	if isConnected {
		fmt.Println("connected")
	} else {
		fmt.Println("disconnected")
	}
}

func main() {
	// Open a new JACK client named "Go Passthrough"
	var status int
	JackClient, status = jack.ClientOpen("Go Passthrough", jack.NoStartServer)
	if status != 0 {
		fmt.Println("Status:", jack.StrError(status))
		return
	}
	defer JackClient.Close()

	// Register input ports based on the number of channels
	for i := 0; i < channels; i++ {
		portIn := JackClient.PortRegister(fmt.Sprintf("in_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsInput, 0)
		PortsIn = append(PortsIn, portIn)
	}

	// Register output ports based on the number of channels
	for i := 0; i < channels; i++ {
		portOut := JackClient.PortRegister(fmt.Sprintf("out_%d", i), jack.DEFAULT_AUDIO_TYPE, jack.PortIsOutput, 0)
		PortsOut = append(PortsOut, portOut)
	}

	// Connect the "Go Passthrough" client to a specified device
	JackClient.Connect("Scarlett 18i8 3rd Gen Pro Monitor:monitor_AUX0", "Go Passthrough:in_0")

	// Set a callback for when port connections change
	if code := JackClient.SetPortConnectCallback(processXX); code != 0 {
		fmt.Println("Failed to set process callback:", jack.StrError(code))
		return
	}

	// Set the main processing callback for the audio data
	if code := JackClient.SetProcessCallback(process); code != 0 {
		fmt.Println("Failed to set process callback:", jack.StrError(code))
		return
	}

	// Create a channel to wait for shutdown
	shutdown := make(chan struct{})
	JackClient.OnShutdown(func() {
		fmt.Println("Shutting down")
		close(shutdown)
	})

	// Activate the JACK client
	if code := JackClient.Activate(); code != 0 {
		fmt.Println("Failed to activate client:", jack.StrError(code))
		return
	}

	fmt.Println(JackClient.GetName())
	<-shutdown // Wait here until a shutdown signal is received
}
