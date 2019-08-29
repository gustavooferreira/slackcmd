package main

import (
	"fmt"

	"github.com/gustavooferreira/slackcmd/pkg/core"
)

func cmd1([]string) string {
	return "YOLO1!"
}
func cmd2([]string) string {
	return "YOLO2!"
}
func cmd3([]string) string {
	return "YOLO3!"
}

func main() {

	menu := core.NeWMenu()

	submenu1 := core.NeWMenu()
	submenu1.AddEntry(core.CreateCommandEntry("cmd1", cmd1, "short help", "long help"))
	submenu1.AddEntry(core.CreateCommandEntry("cmd2", cmd2, "short help", "long help"))

	menu.AddEntry(core.CreateSubMenuEntry("submenu", submenu1, "short help submenu", "long help submenu"))

	fmt.Printf("%+v\n", menu)
}
