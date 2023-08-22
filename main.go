package main

import (
	"fmt"
	"github.com/xthexder/go-jack"
)

var channels int = 2 // Number of channels to process

// Global declarations for JACK audio ports and client
var PortsIn []*jack.Port
var PortsOut []*jack.Port
var JackClient *jack.Client
var isConnected bool // A flag to check if the client is connected

// process handles the audio data for each frame.
// For now, it prints the non-zero samples from the first channel.
func process(nframes uint32) int {
	port := PortsIn[0]

	// Check if the client is connected before processing
	if isConnected {
		samplesIn := port.GetBuffer(nframes)

		// Print non-zero samples for debugging
		for i, sample := range samplesIn {
			if sample != 0 {
				fmt.Println("sample no", i, " :", sample)
			}
		}
	}

	return 0
}

// processXX is a callback for when port connections change.
func processXX(x jack.PortId, y jack.PortId, z bool) {
	isConnected = true
	fmt.Println("connected")
}

func main() {
	var status int
	// Open a new JACK client named "Go Passthrough"
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
