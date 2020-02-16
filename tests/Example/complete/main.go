package main

import (
	"fmt"
	"os"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/inventory"
	"github.com/gustavooferreira/slackcmd/pkg/security"
	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

const banner string = `Command banner.`

func main() {
	mainMenu := inventory.NewMenu()
	mainMenu.AddCommandEntry("cmd1", "short help cm1", "long help cmd1", cmd1)

	submenu1 := mainMenu.AddSubMenuEntry("submenu1", "short help submenu1", "long help submenu1")
	submenu1.AddCommandEntry("cmd2", "short help cmd2", "long help cmd2", cmd2)
	submenu1.AddCommandEntry("cmd3", "short help cmd3", "long help cmd3", cmd3)

	submenu2 := mainMenu.AddSubMenuEntry("submenu2", "short help submenu2", "long help submenu2")
	submenu2.AddCommandEntry("cmd3", "short help cmd3", "long help cmd3", cmd3)

	ci := inventory.NewCommandInventory("isp", banner, mainMenu)

	chPerm := map[string][]string{"GK9U15UJ3": []string{"U451E8XQ8"}}
	perm := security.NewPermissions("T02GEFU92", []string{}, chPerm)

	signingSecret := os.Getenv("SLACK_SS")

	scs := webserver.NewSlashCmdServer(nil, 8080)
	scs.RegisterCommand("/isp", "/slack/isp", &perm, signingSecret, func(rc entities.RequestContext) string {
		return ci.Parse(rc)
	})

	scs.ListenAndServe()
}

func cmd1(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	fmt.Println("Options:", options)
	return "YOLO1!"
}

func cmd2(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	fmt.Println("Options:", options)
	return "YOLO2!"
}

func cmd3(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	fmt.Println("Options:", options)
	return "YOLO3!"
}
