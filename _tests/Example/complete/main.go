package main

import (
	"fmt"
	"os"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/inventory"
	"github.com/gustavooferreira/slackcmd/pkg/layout"
	"github.com/gustavooferreira/slackcmd/pkg/security"
	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

var banner string = `Command banner.`

func main() {
	mainMenu := layout.NewMenu()
	mainMenu.AddEntry(layout.CreateCommandEntryFromStruct(Cmd1{}))
	submenu1 := layout.NewMenu()
	mainMenu.AddEntry(layout.CreateSubMenuEntry("submenu1", submenu1, "short help submenu1", "long help submenu1"))
	submenu1.AddEntry(layout.CreateCommandEntry("cmd2", cmd2, "short help cmd2", "long help cmd2"))
	submenu1.AddEntry(layout.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
	submenu2 := layout.NewMenu()
	mainMenu.AddEntry(layout.CreateSubMenuEntry("submenu2", submenu2, "short help submenu2", "long help submenu2"))
	submenu1.AddEntry(layout.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
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

type Cmd1 struct {
}

func (c Cmd1) GetName() string {
	return "cmdUltra"
}

func (c Cmd1) GetSHelp() string {
	return "YOLO SHORT HELP"
}

func (c Cmd1) GetLHelp() string {
	return "LOOOOONG HELP MSG!"
}

func (c Cmd1) Exec() layout.CmdFunction {
	return func(rc entities.RequestContext, options []string) string {
		fmt.Printf("Request Context %+v\n", rc)
		fmt.Println("Options:", options)
		return "YOLO1!"
	}
}
