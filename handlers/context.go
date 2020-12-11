package handlers

import "github.com/paypermint/appkit"

var appCtx *appkit.AppContext

//SetAppContext sets the application context in the handlers
func SetAppContext(ac *appkit.AppContext) {
	appCtx = ac
}
