package inventory

import "github.com/gustavooferreira/slackcmd/pkg/entities"

// CmdFunction is the function that executes the command.
// It gets the Request context and a list of arguments passed to the command
type CmdFunction func(entities.RequestContext, []string) string

type CmdInterface interface {
	GetName() string
	GetSHelp() string
	GetLHelp() string
	CmdHandler() CmdFunction
}

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
	Type                 MenuEntryType
	Cmd                  CmdFunction
	SubMenu              *menu
}

type menu struct {
	Entries map[string]menuEntry
}

func NewMenu() Menu {
	menuMap := make(map[string]menuEntry)
	return menu{Entries: menuMap}
}

func (m *menu) AddSubMenuEntry(name string, helpsd string, helpld string) menu {
	menuInst := NewMenu()
	entry := menuEntry{Name: name, HelpShortDescription: helpsd, HelpLongDescription: helpld, Type: menuEntryType_SubMenu, SubMenu: &menuInst}
	m.Entries[name] = entry
	return menuInst
}

func (m *menu) AddCommandEntry(name string, helpsd string, helpld string, cmd CmdFunction) {
	entry := menuEntry{Name: name, HelpShortDescription: helpsd, HelpLongDescription: helpld, Type: menuEntryType_Command, Cmd: cmd}
	m.Entries[name] = entry
}

func (m *menu) AddCommandEntryFromStruct(cmdI CmdInterface) {
	entry := menuEntry{Name: cmdI.GetName(),
		HelpShortDescription: cmdI.GetSHelp(),
		HelpLongDescription:  cmdI.GetLHelp(),
		Type:                 menuEntryType_Command,
		Cmd:                  cmdI.CmdHandler()}
	m.Entries[cmdI.GetName()] = entry
}
