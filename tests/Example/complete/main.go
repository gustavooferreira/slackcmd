package main

import (
	"fmt"
	"os"

	"github.com/gustavooferreira/slackcmd/pkg/cmd_inventory"
	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/menu"
	"github.com/gustavooferreira/slackcmd/pkg/permissions"
	"github.com/gustavooferreira/slackcmd/pkg/webserver"
)

var banner string = `ISP command provides general utilities.

You can issue:

/isp tree

To see a tree of commands you can use


/isp help [path to command]

Shows a help message.
`

func main() {
	mainMenu := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateCommandEntryFromStruct(Cmd1{}))
	submenu1 := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateSubMenuEntry("submenu1", submenu1, "short help submenu1", "long help submenu1"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd2", cmd2, "short help cmd2", "long help cmd2"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
	submenu2 := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateSubMenuEntry("submenu2", submenu2, "short help submenu2", "long help submenu2"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))
	ci := cmd_inventory.NewCommandInventory("isp", banner, mainMenu)

	chPerm := map[string][]string{"GK9U15UJ3": []string{"U451E8XQ8"}}

	perm := permissions.Permissions{TeamID: "T02GEFU92", Global: []string{}, Channel: chPerm}

	signingSecret := os.Getenv("SLACK_SS")

	scs := webserver.NewSlashCmdServer(nil, 8080, signingSecret)
	scs.RegisterCommand("/isp", "/slack/isp", &perm, func(rc entities.RequestContext) string {
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

func (c Cmd1) Exec() menu.CmdFunction {
	return func(rc entities.RequestContext, options []string) string {
		fmt.Printf("Request Context %+v\n", rc)
		fmt.Println("Options:", options)
		return "YOLO1!"
	}
}
