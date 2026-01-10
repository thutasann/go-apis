package webrtc

import (
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/thuta/gowebrtc/pkg/chat"
)

type Room struct {
	Peers *Peers
	Hub   *chat.Hub
}

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrackLocalStaticRTP
}

type PeerConnectionState struct {
	PeerConnection *webrtc.PeerConnection
	Websocket      *ThreadSafeWriter
}

type ThreadSafeWriter struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

func (p *Peers) DispatchKeyFrame() {
}
