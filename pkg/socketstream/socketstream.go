package socketstream

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/seamounts/pod-exec/pkg/webshell"
	"gopkg.in/square/go-jose.v2/json"
	"k8s.io/client-go/tools/remotecommand"
)

var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 5
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	upgrader.Error = func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Printf("web socket conn err, code: %d, reason: %s", status, reason)
	}
	return upgrader
}()

type SocketStream struct {
	wsConn   wsconnecter
	sizeChan chan *remotecommand.TerminalSize
	done     chan struct{}
}

func NewSocketStream(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*SocketStream, error) {
	conn, err := newWSConn(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	socketStream := &SocketStream{
		wsConn:   conn,
		sizeChan: make(chan *remotecommand.TerminalSize),
		done:     make(chan struct{}),
	}

	return socketStream, nil
}

// Done done, must call Done() before connection close, or Next() would not exits.
func (t *SocketStream) Done() {
	close(t.done)
}

func (s *SocketStream) Next() *remotecommand.TerminalSize {
	for {
		select {
		case tsize := <-s.sizeChan:
			return tsize
		case <-s.done:
			return nil
		}
	}
}

func (s *SocketStream) Read(p []byte) (n int, err error) {
	_, message, err := s.wsConn.ReadMessage()
	if err != nil {
		log.Printf("read message err: %v", err)
		return copy(p, webshell.EndOfTransmission), err
	}

	tmsg := &webshell.TerminalMessage{}
	if err := json.Unmarshal(message, tmsg); err != nil {
		log.Printf("Unmarshal err: %v", err)
		return copy(p, webshell.EndOfTransmission), err
	}
	log.Printf("read message  %v", tmsg)
	switch tmsg.Operation {
	case webshell.OPERATION_STDIN:
		if len(tmsg.Data) == 0 {
			log.Printf("stdin  len data == 0")
			return 0, nil
		}
		log.Printf("stdin  %v", tmsg.Data)
		return copy(p, tmsg.Data), nil

	case webshell.OPERATION_RESIZE:
		log.Printf("resize  %v", tmsg.Data)
		tsize := &remotecommand.TerminalSize{
			Width:  tmsg.Cols,
			Height: tmsg.Rows,
		}
		s.sizeChan <- tsize
		return 0, nil

	case webshell.OPERATION_PING:
		log.Printf("ping  %v", tmsg.Data)
		return 0, nil

	default:
		log.Printf("default  %v", tmsg.Data)
		return copy(p, webshell.EndOfTransmission), fmt.Errorf("unknown message type '%s'", tmsg.Operation)
	}
}

func (s *SocketStream) Write(p []byte) (n int, err error) {
	tmsg := &webshell.TerminalMessage{
		Operation: webshell.OPERATION_STDOUT,
		Data:      string(p),
	}
	msg, err := json.Marshal(tmsg)
	if err != nil {
		return 0, err
	}

	if err := s.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		return 0, err
	}

	return len(p), nil
}
