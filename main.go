package main

import (
	"flag"

	"github.com/paypermint/appkit"
)

const (
	serviceName = "bridge-app-svc"
)

func main() {
	appkit.RequireFlags(appkit.F_WEB | appkit.F_LOG | appkit.F_HEALTH | appkit.F_REGION)
	flag.Parse()
	config := appkit.GetAppConfig()
	log := appkit.NewLogger(config.Log)
	appctx := appkit.NewAppContext(config, log)

	defer appctx.Cleanup()

	go appkit.StartHealthCheckEndpoint(appctx)

}
