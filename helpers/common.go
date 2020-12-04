package helpers

var (
	dynamicHost          = ""
	identityPublicAPIkey = ""
)

// SetDynamicHost sets host to be used in helpers
func SetDynamicHost(host string) {
	dynamicHost = host
}

//GetDynamicHost gets host
func GetDynamicHost() string {
	return dynamicHost
}

// SetIdentityPublicAPIKey sets host to be used in helpers
func SetIdentityPublicAPIKey(host string) {
	identityPublicAPIkey = host
}

//GetIdentityPublicAPIKey gets host
func GetIdentityPublicAPIKey() string {
	return identityPublicAPIkey
}
