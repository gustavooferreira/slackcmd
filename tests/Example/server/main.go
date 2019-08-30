package main

import (
	"fmt"

	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

func main() {
	perm := webserver.Permissions{}

	scs := webserver.NewSlashCmdServer(nil, 8080)
	scs.RegisterCommand("/isp", "/slack/isp", perm, Hello)
	scs.ListenAndServe()
}

func Hello(rc webserver.RequestContext) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "Hello ma friend!"
}
