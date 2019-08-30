package main

import (
	"fmt"
	"strings"

	"github.com/gustavooferreira/slackcmd/pkg/core"
	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

func main() {
	menu := core.NeWMenu()
	menu.AddEntry(core.CreateCommandEntry("cmd1", cmd1, "short help cm1", "long help cmd1"))
	submenu1 := core.NeWMenu()
	menu.AddEntry(core.CreateSubMenuEntry("submenu1", submenu1, "short help submenu1", "long help submenu1"))
	submenu1.AddEntry(core.CreateCommandEntry("cmd2", cmd2, "short help cmd2", "long help cmd2"))
	submenu1.AddEntry(core.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
	submenu2 := core.NeWMenu()
	menu.AddEntry(core.CreateSubMenuEntry("submenu2", submenu2, "short help submenu2", "long help submenu2"))
	submenu1.AddEntry(core.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
	ci := core.NewCommandInventory("isp", menu)

	perm := webserver.Permissions{}

	scs := webserver.NewSlashCmdServer(nil, 8080)
	scs.RegisterCommand("/isp", "/slack/isp", perm, func(rc webserver.RequestContext) string {
		cmdArr := strings.Split(rc.Text, " ")

		result, _ := ci.Logic(rc, cmdArr)
		return result
	})

	scs.ListenAndServe()
}

func cmd1(rc webserver.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO1!"
}
func cmd2(rc webserver.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO2!"
}
func cmd3(rc webserver.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO3!"
}
