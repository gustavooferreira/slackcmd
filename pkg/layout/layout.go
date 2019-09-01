package layout

import "github.com/gustavooferreira/slackcmd/pkg/entities"

type Menu struct {
	Entries map[string]MenuEntry
}

func NewMenu() Menu {
	menuMap := make(map[string]MenuEntry)
	return Menu{Entries: menuMap}
}

func (m *Menu) AddEntry(entry MenuEntry) {
	name := entry.Name
	m.Entries[name] = entry
}

type MenuEntry struct {
	Name                 string
	Target               interface{}
	HelpShortDescription string
	HelpLongDescription  string
}

type CmdFunction func(entities.RequestContext, []string) string

type CmdInterface interface {
	GetName() string
	GetSHelp() string
	GetLHelp() string
	Exec() CmdFunction
}

func CreateCommandEntryFromStruct(cmdStruct CmdInterface) MenuEntry {
	return MenuEntry{Name: cmdStruct.GetName(), Target: cmdStruct.Exec(),
		HelpShortDescription: cmdStruct.GetSHelp(), HelpLongDescription: cmdStruct.GetLHelp()}
}

func CreateCommandEntry(name string, cmd CmdFunction, helpsd string, helpld string) MenuEntry {
	return MenuEntry{Name: name, Target: cmd, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}

func CreateSubMenuEntry(name string, menu Menu, helpsd string, helpld string) MenuEntry {
	return MenuEntry{Name: name, Target: menu, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}
