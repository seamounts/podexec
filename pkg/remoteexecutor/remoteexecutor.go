package remoteexecutor

import (
	"net/url"

	"github.com/seamounts/pod-exec/pkg/kube"
	"k8s.io/client-go/tools/remotecommand"
)

// RemoteExecutor defines the interface accepted by the Exec command - provided for test stubbing
type RemoteExecutor interface {
	Execute(method string, url *url.URL) error
}

type PtyHandler interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Next() *remotecommand.TerminalSize
	Done()
}

// DefaultRemoteExecutor is the standard implementation of remote command execution
type DefaultRemoteExecutor struct{}

func NewDefaultExecutor() *DefaultRemoteExecutor {
	return &DefaultRemoteExecutor{}
}

func (*DefaultRemoteExecutor) Execute(method string, url *url.URL, pty PtyHandler) error {
	defer func() {
		pty.Done()
	}()

	config, err := kube.Config()
	if err != nil {
		return err
	}

	exec, err := remotecommand.NewSPDYExecutor(config, method, url)
	if err != nil {
		return err
	}
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
}
