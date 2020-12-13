package helpers

var (
	dynamicHost      = ""
	sapUserCredsPath = ""
)

// SetDynamicHost sets host to be used in helpers
func SetDynamicHost(host string) {
	dynamicHost = host
}

//GetDynamicHost gets host
func GetDynamicHost() string {
	return dynamicHost
}

// SetDynamicHost sets host to be used in helpers
func SetSapUserCredsPath(path string) {
	sapUserCredsPath = path
}

//GetDynamicHost gets host
func GetSapUserCredsPath() string {
	return sapUserCredsPath
}
