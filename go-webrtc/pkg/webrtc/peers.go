package webrtc

import (
	"github.com/pion/webrtc/v3"
	"github.com/thuta/gowebrtc/pkg/chat"
)

type Room struct {
	Peers *Peers
	Hub   *chat.Hub
}

type Peers struct {
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

func (p *Peers) DispatchKeyFrame() {
}
