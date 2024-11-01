package webrtc_plugin

import (
	"bufio"
	"fmt"
	"os"

	"webrtc_utils"

	"github.com/pion/webrtc/v4"
)

type WebRTCPlugin struct {
	Connection *webrtc.PeerConnection
	DataCh     *webrtc.DataChannel
	DataChOpen chan struct{}
}

func NewWebRTCPlugin() (*WebRTCPlugin, error) {
	api := webrtc.NewAPI()
	connection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{
				URLs:       []string{"turn:60.204.212.12:3478"},
				Username:   "bobby",
				Credential: "bobby",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	dataChOpen := make(chan struct{})

	plugin := &WebRTCPlugin{
		Connection: connection,
		DataChOpen: dataChOpen,
	}

	connection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Println("ICE Connection State:", state)
	})

	connection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			fmt.Println("Candidate found", candidate)
		} else {
			// 全部 ICE 候选者已收集完毕
			return
		}
	})

	gatherComplete := make(chan struct{})
	connection.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		if state == webrtc.ICEGatheringStateComplete {
			gatherComplete <- struct{}{}
		}
	})

	// 阻塞等待gatherComplete
	// <-gatherComplete
	return plugin, nil
}

func (plugin *WebRTCPlugin) CreateDataChannel(label string) error {
	dataCh, err := plugin.Connection.CreateDataChannel(label, &webrtc.DataChannelInit{})
	if err != nil {
		return err
	}
	plugin.DataCh = dataCh

	dataCh.OnOpen(func() {
		fmt.Println("DataChannel Opened")
		plugin.DataChOpen <- struct{}{}
	})

	return nil
}

func (plugin *WebRTCPlugin) SendData(data string) error {
	<-plugin.DataChOpen
	return plugin.DataCh.Send([]byte(data))
}

func (plugin *WebRTCPlugin) ReceiveData() (string, error) {
	<-plugin.DataChOpen
	reader := bufio.NewReader(os.Stdin)
	var buf string
	_, err := fmt.Fscanf(reader, "%s", &buf)
	return buf, err
}

func (plugin *WebRTCPlugin) SetRemoteSDP(sdp string) error {
	if err := webrtc_utils.ValidateSDP(sdp); err != nil {
		return err
	}
	remoteSDP, err := webrtc_utils.DecodeSDP(sdp)
	if err != nil {
		return err
	}
	return plugin.Connection.SetRemoteDescription(*remoteSDP)
}

func (plugin *WebRTCPlugin) GetLocalSDP() (string, error) {
	initOffer, err := plugin.Connection.CreateOffer(nil)
	if err != nil {
		return "", err
	}
	err = plugin.Connection.SetLocalDescription(initOffer)
	if err != nil {
		return "", err
	}
	return webrtc_utils.EncodeSDP(plugin.Connection.LocalDescription())
}
