package main

import (
	"fmt"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/inventory"
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

	mainMenu := inventory.NewMenu()
	mainMenu.AddCommandEntry("cmd1", "short help cm1", "long help cmd1", cmd1)

	submenu1 := mainMenu.AddSubMenuEntry("submenu1", "short help submenu1", "long help submenu1")
	submenu1.AddCommandEntry("cmd2", "short help cmd2", "long help cmd2", cmd2)
	submenu1.AddCommandEntry("cmd3", "short help cmd3", "long help cmd3", cmd3)

	submenu2 := mainMenu.AddSubMenuEntry("submenu2", "short help submenu2", "long help submenu2")
	submenu2.AddCommandEntry("cmd3", "short help cmd3", "long help cmd3", cmd3)

	ci := inventory.NewCommandInventory("isp", "banner!", "v1.5.2", mainMenu)

	// ------------------------------------------

	fmt.Printf("Main Menu: %+v\n", mainMenu)
	fmt.Println("--------------------------------------\n")

	cmdArr := []string{"submenu1", "cmd2"}
	// cmdArr := []string{"submenu"}
	// cmdArr := []string{"submenu", "cmd2", "yolo"}

	a, b, err := ci.HelpFunc(cmdArr)
	fmt.Println(a, "---", b, "---", err)
	fmt.Println("--------------------------------------\n")

	f, err := ci.Match(cmdArr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(f(entities.RequestContext{}, []string{}))
	}
	fmt.Println("--------------------------------------\n")

	result, err := ci.Tree(nil, 0)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println(result)
	}
}
