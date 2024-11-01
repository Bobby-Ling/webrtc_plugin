package main

import (
	"fmt"
	"log"
	"os"

	"webrtc_plugin"
)

func main() {
	log.SetOutput(os.Stdout)

	plugin, err := webrtc_plugin.NewWebRTCPlugin()
	if err != nil {
		log.Fatal("Failed to initialize WebRTC plugin:", err)
	}

	err = plugin.CreateDataChannel("data")
	if err != nil {
		log.Fatal("Failed to create data channel:", err)
	}

	localSDP, err := plugin.GetLocalSDP()
	if err != nil {
		log.Fatal("Failed to get local SDP:", err)
	}
	fmt.Println("Local SDP:", localSDP)

	var remoteSDP string
	fmt.Printf("Input Remote SDP:")
	fmt.Scan(&remoteSDP)

	err = plugin.SetRemoteSDP(remoteSDP)
	if err != nil {
		log.Fatal("Failed to set remote SDP:", err)
	}

	err = plugin.SendData("Hello WebRTC!")
	if err != nil {
		log.Fatal("Failed to send data:", err)
	}
}
