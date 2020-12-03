package main

import (
	"flag"
	"github.com/paypermint/mskit"
)

const (
	serviceName = "bridge-app-svc"
)

func main() {
	mskit.Init(serviceName)
	mskit.RequireFlags(mskit.F_GRPC | mskit.F_LOG | mskit.F_HEALTH | mskit.F_KEY | mskit.F_REGION)x
	mskit.LoadValidators()
	flag.Parse()

	config := mskit.GetConfig()

	log := mskit.NewLogger(config.Log)
	sctx := mskit.NewServiceContext(config, log)
	defer sctx.Cleanup()

	go mskit.StartHealthCheckEndpoint(sctx)
}
