package app

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo"
	"github.com/seamounts/pod-exec/pkg/remoteexecutor"
	"github.com/seamounts/pod-exec/pkg/socketstream"
	"github.com/spf13/cobra"
)

const (
	POST = "POST"
)

func NewAppCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "hydra agent login test",
		Long: "hydra agent login test",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	return cmd
}

func run() {
	e := echo.New()
	e.GET("/", serveTerminal)
	e.GET("/webshell/:cluster", serveWsTerminal)
	e.Static("/static", "./frontend")

	e.Start(":8080")
}

func serveWsTerminal(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	websocketStream, err := socketstream.NewSocketStream(w, r, nil)
	if err != nil {
		return err
	}

	u, err := url.Parse("https://xxx.com:443")
	if err != nil {
		return err
	}
	u.Path = fmt.Sprintf("/cluster-management/%s/shell", c.Param("cluster"))

	return remoteexecutor.NewDefaultExecutor().Execute(POST, u, websocketStream)
}

func serveTerminal(c echo.Context) error {
	return c.File("/frontend/terminal.html")
}
