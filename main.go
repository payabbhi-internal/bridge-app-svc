package main

import (
	"flag"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/handlers"
	"github.com/paypermint/bridge-app-svc/helpers"
	"github.com/paypermint/bridge-app-svc/interceptors"
	"github.com/unrolled/render"
	"github.com/unrolled/secure"
	"github.com/urfave/negroni"
	"google.golang.org/grpc/grpclog"
)

const (
	serviceName = "bridge-app-svc"
)

var (
	dynamicHost  = flag.String("dynamic-host", "payabbhi.com", "Dynamic host")
	bucketRegion = flag.String("bucket-region", "", "Region for AWS where the bucket for file upload has been created")
)

func main() {
	appkit.RequireFlags(appkit.F_WEB | appkit.F_LOG | appkit.F_HEALTH | appkit.F_REGION)
	flag.Parse()
	config := appkit.GetAppConfig()
	log := appkit.NewLogger(config.Log)
	appctx := appkit.NewAppContext(config, log)

	defer appctx.Cleanup()
	grpclog.SetLogger(appkit.NewGrpcLogger(log))
	handlers.SetAppContext(appctx)
	go appkit.StartHealthCheckEndpoint(appctx)
	helpers.SetDynamicHost(*dynamicHost)
	helpers.SetBucketConfig(*bucketRegion)
	appctx.Renderer = render.New(render.Options{
		IndentJSON: true,
	})

	secureMiddleware := secure.New(secure.Options{
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
	})

	//sets routes
	router := setRouter()
	// Start server
	n := negroni.New()
	n.Use(interceptors.NewRecoveryInterceptor(appctx))
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(interceptors.NewLoggingInterceptor(appctx))
	n.UseHandler(router)
	appkit.StartWeb(appctx, n)
}
