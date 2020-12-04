package interceptors

import (
	"net/http"
	"time"

	"github.com/paypermint/appkit"
	"github.com/urfave/negroni"
)

// LoggingInterceptor is a middleware that logs the request as it goes in and the response as it goes out.
type LoggingInterceptor struct {
	appCtx *appkit.AppContext
}

// NewLoggingInterceptor returns a new instance of LoggingInterceptor
func NewLoggingInterceptor(appctx *appkit.AppContext) *LoggingInterceptor {
	return &LoggingInterceptor{appCtx: appctx}
}
func (rec *LoggingInterceptor) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	traceID := appkit.TraceIDFromHTTPRequest(r)

	ctxlogger := rec.appCtx.Logger.New("method", r.Method, "url", r.URL.Path, "traceid", traceID)
	ctxlogger.Info("Started request")

	// profileID := util.ProfileIDFromHTTPRequest(r)
	// apiVersion := util.VersionFromHTTPRequest(r)

	rw.Header().Set("Cache-Control", "no-cache, no-store")
	// rw.Header().Set("Payabbhi-Version", apiVersion)
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// //for checkout pages profileID will not be present in request header
	// if strings.HasPrefix(r.URL.Path, "/api/v1/invc") || strings.HasPrefix(r.URL.Path, "/api/v1/emandate") || strings.HasPrefix(r.URL.Path, "/api/v1/profiles") || strings.HasPrefix(r.URL.Path, "/api/v1/profile_identity") {
	// 	if traceID == "" || apiVersion == "" {
	// 		ctxlogger.Info("Completed request", "error_message", "missing traceid or Payabbhi-Version in headers", "status", http.StatusInternalServerError, "time_taken", time.Since(start))
	// 		util.RenderAPIErrorJSON(rec.appCtx, rw)
	// 		return
	// 	}
	// } else {
	// 	if profileID == "" || traceID == "" || apiVersion == "" {
	// 		ctxlogger.Info("Completed request", "error_message", "missing Profile-Id or traceid or Payabbhi-Version in headers", "status", http.StatusInternalServerError, "time_taken", time.Since(start))
	// 		util.RenderAPIErrorJSON(rec.appCtx, rw)
	// 		return
	// 	}
	// }

	if r.Method == "OPTIONS" {
		ctxlogger.Info("Completed request", "status", http.StatusOK, "time_taken", time.Since(start))
		return
	}

	next(rw, r)
	res := rw.(negroni.ResponseWriter)
	ctxlogger.Info("Completed request", "status", res.Status(), "time_taken", time.Since(start))
}
