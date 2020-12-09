package helpers

var (
	dynamicHost = ""
)

// SetDynamicHost sets host to be used in helpers
func SetDynamicHost(host string) {
	dynamicHost = host
}

//GetDynamicHost gets host
func GetDynamicHost() string {
	return dynamicHost
}
