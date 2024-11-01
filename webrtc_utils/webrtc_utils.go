package webrtc_utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pion/webrtc/v4"
)

type Configuration struct {
	ICEServers []webrtc.ICEServer `json:"iceServers"`
}

func ParseICEServersFromJSON(iceServersJSON string) ([]webrtc.ICEServer, error) {
	var config Configuration
	err := json.Unmarshal([]byte(iceServersJSON), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ICE servers JSON: %v", err)
	}
	return config.ICEServers, nil
}

func EncodeSDP(sdp *webrtc.SessionDescription) (string, error) {
	sdpJSON, err := json.Marshal(sdp)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	g, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", err
	}
	defer g.Close()
	if _, err = g.Write(sdpJSON); err != nil {
		return "", err
	}

	if err = g.Close(); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func DecodeSDP(in string) (*webrtc.SessionDescription, error) {
	buf, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	sdpBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var sdp webrtc.SessionDescription
	err = json.Unmarshal(sdpBytes, &sdp)
	if err != nil {
		return nil, err
	}

	return &sdp, nil
}

func ValidateSDP(input string) error {
	buf, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return err
	}
	r, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return err
	}
	defer r.Close()

	sdpBytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var sdp webrtc.SessionDescription
	err = json.Unmarshal(sdpBytes, &sdp)
	if err != nil {
		return err
	}

	return nil
}
