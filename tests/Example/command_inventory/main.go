package main

import (
	"fmt"

	"github.com/gustavooferreira/slackcmd/pkg/cmd_inventory"
	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/menu"
)

func cmd1(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO1!"
}
func cmd2(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO2!"
}
func cmd3(rc entities.RequestContext, options []string) string {
	fmt.Printf("Request Context %+v\n", rc)
	return "YOLO3!"
}

func main() {
	mainMenu := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateCommandEntry("cmd1", cmd1, "short help cm1", "long help cmd1"))

	submenu1 := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateSubMenuEntry("submenu1", submenu1, "short help submenu1", "long help submenu1"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd2", cmd2, "short help cmd2", "long help cmd2"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))

	submenu2 := menu.NeWMenu()
	mainMenu.AddEntry(menu.CreateSubMenuEntry("submenu2", submenu2, "short help submenu2", "long help submenu2"))
	submenu1.AddEntry(menu.CreateCommandEntry("cmd3", cmd3, "short help cmd3", "long help cmd3"))

	ci := cmd_inventory.NewCommandInventory("isp", mainMenu)

	// ------------------------------------------

	fmt.Printf("%+v\n", mainMenu)
	fmt.Println("--------------------------------------\n")

	cmdArr := []string{"submenu", "cmd1"}
	// cmdArr := []string{"submenu"}
	// cmdArr := []string{"submenu", "cmd2", "yolo"}

	a, b, err := ci.HelpFunc(cmdArr)

	fmt.Println(a, b, err)

	fmt.Println("--------------------------------------\n")

	f, err := ci.Match(cmdArr)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f(entities.RequestContext{}, []string{}))
	}

	fmt.Println("--------------------------------------\n")

	result, err := ci.Tree(nil, 0)

	fmt.Println(result)
}
