package tcp

import (
	"github.com/davyxu/cellnet"
	cellevent "github.com/davyxu/cellnet/event"
	cellpeer "github.com/davyxu/cellnet/peer"
	cellqueue "github.com/davyxu/cellnet/queue"
	xframe "github.com/davyxu/x/frame"
)

type Peer struct {
	cellpeer.Hooker
	*cellpeer.SessionManager
	xframe.PropertySet
	cellpeer.SocketOption
	cellpeer.Protect
	Queue *cellqueue.Queue
	Recv  func(ses *Session) (ev *cellevent.RecvMsgEvent, err error)
	Send  func(ses *Session, ev *cellevent.SendMsgEvent) error
}

func (self *Peer) Peer() *Peer {
	return self
}

func newPeer() *Peer {
	return &Peer{
		SessionManager: cellpeer.NewSessionManager(),
	}
}

// SessionID根据各种实现不一样(例如网关), 应该在具体实现里获取
func SessionID(ses cellnet.Session) int64 {
	if ses == nil {
		return 0
	}

	type idfetcher interface {
		ID() int64
	}

	if f, ok := ses.(idfetcher); ok {
		return f.ID()
	}

	return 0
}

func SessionPeer(ses cellnet.Session) *Peer {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		return tcpSes.Peer
	}

	return nil
}

func SessionParent(ses cellnet.Session) interface{} {
	if ses == nil {
		return nil
	}

	if tcpSes, ok := ses.(*Session); ok {
		return tcpSes.parent
	}

	return nil
}

func SessionQueuedCall(ses cellnet.Session, callback func()) {
	peer := SessionPeer(ses)
	if peer == nil {
		return
	}

	if peer.Queue == nil {
		callback()
	} else {
		peer.Queue.Post(callback)
	}
}