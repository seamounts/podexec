package socketstream

import (
	"encoding/json"
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/rancher/machine/libmachine/log"
	"github.com/seamounts/pod-exec/pkg/webshell"
	"k8s.io/client-go/tools/remotecommand"
)

func newMockWSConn(t *testing.T) (*Mockwsconnecter, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	return NewMockwsconnecter(ctrl), ctrl
}

func newMockSocketStream(t *testing.T) *SocketStream {
	socketStream := &SocketStream{
		sizeChan: make(chan *remotecommand.TerminalSize),
		done:     make(chan struct{}),
	}

	return socketStream
}

func TestRead(t *testing.T) {
	ss := newMockSocketStream(t)
	conn, ctrl := newMockWSConn(t)
	defer ctrl.Finish()

	testCases := []struct {
		Message      string
		Data         *webshell.TerminalMessage
		ExpectErr    error
		ApplyMock    func(data *webshell.TerminalMessage)
		ExpactResult func(p []byte, err error) error
	}{
		{
			Message: "unknown message type",
			Data: &webshell.TerminalMessage{
				Operation: "errOperation",
			},
			ApplyMock: func(data *webshell.TerminalMessage) {
				conn.EXPECT().ReadMessage().Return(0, marshalData(data), nil)
				ss.wsConn = conn
			},

			ExpactResult: func(p []byte, err error) error {
				if err.Error() != fmt.Errorf("unknown message type '%s'", "errOperation").Error() {
					return fmt.Errorf("Test failed %s", err.Error())
				}
				return nil
			},
		},
		{
			Message: "stdin operation",
			Data: &webshell.TerminalMessage{
				Operation: webshell.OPERATION_STDIN,
				Data:      "test stdin operation",
			},
			ApplyMock: func(data *webshell.TerminalMessage) {
				conn.EXPECT().ReadMessage().Return(0, marshalData(data), nil)
				ss.wsConn = conn
			},

			ExpactResult: func(p []byte, err error) error {
				if err != nil {
					return fmt.Errorf("Test failed %s", err.Error())
				}

				if string(p) != "test stdin operation" {
					return fmt.Errorf("Test failed read data [%s] not match input data", p)
				}

				return nil
			},
		},

		{
			Message: "ping operation",
			Data: &webshell.TerminalMessage{
				Operation: webshell.OPERATION_PING,
				Data:      "test ping operation",
			},
			ApplyMock: func(data *webshell.TerminalMessage) {
				conn.EXPECT().ReadMessage().Return(0, marshalData(data), nil)
				ss.wsConn = conn
			},
			ExpactResult: func(p []byte, err error) error {
				if err != nil {
					return fmt.Errorf("Test failed %s", err.Error())
				}

				return nil
			},
		},
	}

	for _, tc := range testCases {
		log.Infof("test case [%s]", tc.Message)
		tc.ApplyMock(tc.Data)

		pbytes := make([]byte, len(tc.Data.Data))
		_, err := ss.Read(pbytes)
		if err = tc.ExpactResult(pbytes, err); err != nil {
			t.Fatal(err)
		}
	}
}

func TestWrite(t *testing.T) {
	ss := newMockSocketStream(t)
	conn, ctrl := newMockWSConn(t)
	defer ctrl.Finish()

	testCases := []struct {
		Message      string
		Data         string
		ExpectErr    error
		ApplyMock    func(data []byte)
		ExpactResult func(err error) error
	}{
		{
			Message: "write msg success",
			Data:    "write msg success",
			ApplyMock: func(data []byte) {
				tmsg := &webshell.TerminalMessage{
					Operation: webshell.OPERATION_STDOUT,
					Data:      string(data),
				}
				msg, _ := json.Marshal(tmsg)

				conn.EXPECT().WriteMessage(websocket.TextMessage, msg).Return(nil)
				ss.wsConn = conn
			},
			ExpactResult: func(err error) error {
				if err != nil {
					return fmt.Errorf("Test failed %s", err.Error())
				}

				return nil
			},
		},
		{
			Message: "write msg failed",
			Data:    "write msg failed",
			ApplyMock: func(data []byte) {
				tmsg := &webshell.TerminalMessage{
					Operation: webshell.OPERATION_STDOUT,
					Data:      string(data),
				}
				msg, _ := json.Marshal(tmsg)

				conn.EXPECT().WriteMessage(websocket.TextMessage, msg).Return(fmt.Errorf("write msg failed"))
				ss.wsConn = conn
			},
			ExpactResult: func(err error) error {

				if err != nil && err.Error() == "write msg failed" {
					return nil
				}
				return fmt.Errorf("Test failed %s", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		log.Infof("test case [%s]", tc.Message)
		tc.ApplyMock([]byte(tc.Data))

		_, err := ss.Write([]byte(tc.Data))
		if err = tc.ExpactResult(err); err != nil {
			t.Fatal(err)
		}
	}
}

func marshalData(data *webshell.TerminalMessage) []byte {
	dataBytes, _ := json.Marshal(data)

	return dataBytes
}
