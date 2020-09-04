package socketstream

import "net/http"

type wsconnecter interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

func newWSConn(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (wsconnecter, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
