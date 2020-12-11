package util

import (
	"net/http"

	"html/template"

	"github.com/paypermint/bridge-app-svc/models"

	"github.com/paypermint/appkit"
)

const (
	codeMappingNotFound    = "Error code mapping not found for "
	categoryInvalidRequest = "invalid_request_error"
	categoryGatewayError   = "gateway_error"
	categoryAPIError       = "api_error"
	apiErrorMsg            = "There is some problem with the server"
)

//RenderJSON renders JSON
func RenderJSON(appctx *appkit.AppContext, w http.ResponseWriter, status int, v interface{}) error {
	return appctx.Renderer.JSON(w, status, v)
}

//RenderErrorJSON renders error JSON
func RenderErrorJSON(appctx *appkit.AppContext, w http.ResponseWriter, status int, message, field string) error {
	businessErr := models.APIBusinessError{
		Type:    categoryInvalidRequest,
		Message: message,
		Field:   field,
	}
	return RenderJSON(appctx, w, status, models.Error{Err: businessErr})
}

//RenderAPIErrorJSON renders internal error json
func RenderAPIErrorJSON(appctx *appkit.AppContext, w http.ResponseWriter) error {
	businessErr := models.APIBusinessError{
		Type:    categoryAPIError,
		Message: apiErrorMsg,
	}
	return RenderJSON(appctx, w, http.StatusInternalServerError, models.Error{Err: businessErr})
}

//RenderGatewayErrorJSON renders gateway error JSON
func RenderGatewayErrorJSON(appctx *appkit.AppContext, w http.ResponseWriter, status int, message string) error {
	businessErr := models.APIBusinessError{
		Type:    categoryGatewayError,
		Message: message,
	}
	return RenderJSON(appctx, w, status, models.Error{Err: businessErr})
}

// RenderTemplate ...
func RenderTemplate(appCtx *appkit.AppContext, w http.ResponseWriter, template *template.Template, binding interface{}) {
	err := template.Execute(w, binding)
	if err != nil {
		appCtx.Logger.Error("Template Execution failed", "error_message", err.Error())
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
}
