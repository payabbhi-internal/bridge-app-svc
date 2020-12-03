package main

import "github.com/paypermint/mskit"

const (
	serviceName = "bridge-app-svc"
)

func main() {
	mskit.Init(serviceName)
}
