package inventory

import (
	"testing"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
)

func TestInventoryLookup(t *testing.T) {
	cmd1 := func(rc entities.RequestContext, options []string) string {
		return "cmd1"
	}

	mainMenu := NewMenu()
	submenu1 := mainMenu.AddSubMenuEntry("submenu1", "short help submenu1", "long help submenu1")
	submenu1.AddCommandEntry("cmd1", "short help cmd1", "long help cmd1", cmd1)

	ci := NewCommandInventory("cmd", "my banner", "v1.5.2", mainMenu)

	me, err := ci.lookup([]string{"submenu1", "cmd1"})
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	t.Logf("MenuEntry: %+v", me)

	entryName := me.Name

	if entryName != "cmd1" {
		t.Errorf("Menu entry name doesn't match, got: %s, want: %s.", entryName, "cmd1")
	}
}

func TestInventoryTree(t *testing.T) {
	cmd1 := func(rc entities.RequestContext, options []string) string {
		return "cmd1"
	}

	cmd2 := func(rc entities.RequestContext, options []string) string {
		return "cmd2"
	}

	mainMenu := NewMenu()
	mainMenu.AddCommandEntry("cmd1", "short help cmd1", "long help cmd1", cmd1)
	mainMenu.AddCommandEntry("nil", "short help cmdNil", "long help cmdNil", nil)
	submenu1 := mainMenu.AddSubMenuEntry("submenu1", "short help submenu1", "long help submenu1")
	submenu1.AddCommandEntry("cmd2", "short help cmd2", "long help cmd2", cmd2)
	submenu1.AddCommandEntry("nil", "short help cmdNil", "long help cmdNil", nil)
	submenu2 := mainMenu.AddSubMenuEntry("submenu2", "short help submenu2", "long help submenu2")
	submenu3 := submenu2.AddSubMenuEntry("submenu3", "short help submenu3", "long help submenu3")
	submenu3.AddCommandEntry("cmd2", "short help cmd2", "long help cmd2", cmd2)
	submenu3.AddCommandEntry("nil", "short help cmdNil", "long help cmdNil", nil)

	ci := NewCommandInventory("cmd", "my banner", "v1.5.2", mainMenu)

	treeStr, err := ci.Tree(nil, -1)
	if err != nil {
		t.Errorf("error %s", err.Error())
	}

	t.Logf("Tree:\n%s", treeStr)

	// if entryName != "cmd1" {
	// 	t.Errorf("Menu entry name doesn't match, got: %s, want: %s.", entryName, "cmd1")
	// }
}
