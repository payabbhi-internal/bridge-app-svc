package interceptors

import (
	"net/http"
	"runtime"

	"github.com/paypermint/appkit"
	"github.com/paypermint/bridge-app-svc/util"
)

//RecoveryInterceptor defines Recovery structure
type RecoveryInterceptor struct {
	appCtx *appkit.AppContext
}

// NewRecoveryInterceptor returns a new instance of Recovery
func NewRecoveryInterceptor(appctx *appkit.AppContext) *RecoveryInterceptor {
	return &RecoveryInterceptor{appCtx: appctx}
}

func (rec *RecoveryInterceptor) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1024*8)
			stack = stack[:runtime.Stack(stack, false)]

			rec.appCtx.Logger.Crit("PANIC :", "stacktrace", string(stack))
			util.RenderAPIErrorJSON(rec.appCtx, rw)
		}
	}()

	next(rw, r)
}
