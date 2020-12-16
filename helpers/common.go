package helpers

var (
	dynamicHost      = ""
	sapUserCredsPath = ""
	sapURL           = ""
)

// SetDynamicHost sets host to be used in helpers
func SetDynamicHost(host string) {
	dynamicHost = host
}

//GetDynamicHost gets host
func GetDynamicHost() string {
	return dynamicHost
}

// SetSapUserCredsPath sets SapUserCredsPath to be used in helpers
func SetSapUserCredsPath(path string) {
	sapUserCredsPath = path
}

//GetSapUserCredsPath gets SapUserCredsPath
func GetSapUserCredsPath() string {
	return sapUserCredsPath
}

// SetSapURL sets sap url to be used in helpers
func SetSapURL(url string) {
	sapURL = url
}

//GetSapURL gets sap url
func GetSapURL() string {
	return sapURL
}
