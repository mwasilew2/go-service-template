package main

import "golang.org/x/exp/slog"

const (
	componentNameAppClient = "clientCmd"
)

type clientCmd struct {
	// cli options
	httpAddr string `help:"address of the http server which the client should connect to" default:":8080"`
	grpcAddr string `help:"address of the grpc server which the client should connect to" default:":8081"`

	// Dependencies
	logger *slog.Logger
}

func (c *clientCmd) Run(cmdCtx *cmdContext) error {
	c.logger = cmdCtx.Logger.With("component", componentNameAppClient)
	c.logger.Info("starting client", "httpAddr", c.httpAddr, "grpcAddr", c.grpcAddr)

	return nil
}
