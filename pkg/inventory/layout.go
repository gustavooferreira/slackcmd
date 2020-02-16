package inventory

import "github.com/gustavooferreira/slackcmd/pkg/entities"

// CmdFunction is the function that executes the command.
// It gets the Request context and a list of arguments passed to the command
type CmdFunction func(entities.RequestContext, []string) string

type menuEntryType int

const (
	menuEntryType_SubMenu menuEntryType = iota
	menuEntryType_Command
)

func (met menuEntryType) String() string {
	return [...]string{"SubMenu", "Command"}[met]
}

type menuEntry struct {
	Name                 string
	HelpShortDescription string
	HelpLongDescription  string
	Type                 menuEntryType
	Cmd                  CmdFunction
	SubMenu              *Menu
}

type Menu struct {
	Entries map[string]menuEntry
}

func NewMenu() Menu {
	menuMap := make(map[string]menuEntry)
	return Menu{Entries: menuMap}
}

func (m *Menu) AddSubMenuEntry(name string, helpsd string, helpld string) Menu {
	menuInst := NewMenu()
	entry := menuEntry{Name: name, HelpShortDescription: helpsd, HelpLongDescription: helpld, Type: menuEntryType_SubMenu, SubMenu: &menuInst}
	m.Entries[name] = entry
	return menuInst
}

func (m *Menu) AddCommandEntry(name string, helpsd string, helpld string, cmd CmdFunction) {
	entry := menuEntry{Name: name, HelpShortDescription: helpsd, HelpLongDescription: helpld, Type: menuEntryType_Command, Cmd: cmd}
	m.Entries[name] = entry
}
