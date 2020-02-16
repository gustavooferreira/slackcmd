package inventory

import (
	"testing"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
)

func TestMenu(t *testing.T) {
	menu := NewMenu()
	entriesCount := len(menu.Entries)

	if entriesCount != 0 {
		t.Errorf("Menu entries was incorrect, got: %d, want: %d.", entriesCount, 0)
	}
}

func TestMenuEntry(t *testing.T) {
	cmd1 := func(rc entities.RequestContext, options []string) string {
		return "cmd1"
	}

	mainMenu := NewMenu()
	submenu1 := mainMenu.AddSubMenuEntry("submenu1", "short help submenu1", "long help submenu1")
	submenu1.AddCommandEntry("cmd1", "short help cmd1", "long help cmd1", cmd1)

	entryName := mainMenu.Entries["submenu1"].SubMenu.Entries["cmd1"].Name

	if entryName != "cmd1" {
		t.Errorf("Menu entry name doesn't match, got: %s, want: %s.", entryName, "cmd1")
	}
}
