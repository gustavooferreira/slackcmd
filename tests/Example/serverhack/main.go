package main

import (
	"fmt"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/security"
	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

func main() {
	perm := security.Permissions{}

	scs := webserver.NewSlashCmdServer(nil, 8080)
	scs.RegisterCommand("/isp", "/slack/isp", &perm, "signSecret", Hello)
	scs.ListenAndServe()
}

func Hello(rc entities.RequestContext) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "Hello ma friend!"
}
