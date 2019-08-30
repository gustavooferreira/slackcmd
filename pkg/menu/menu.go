package menu

import "github.com/gustavooferreira/slackcmd/pkg/entities"

type Menu struct {
	Menu map[string]MenuEntry
}

func NeWMenu() Menu {
	menuMap := make(map[string]MenuEntry)
	return Menu{Menu: menuMap}
}

func (mm *Menu) AddEntry(entry MenuEntry) {
	entryName := entry.EntryName
	mm.Menu[entryName] = entry
}

type MenuEntry struct {
	EntryName            string
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
	return MenuEntry{EntryName: cmdStruct.GetName(), Target: cmdStruct.Exec(),
		HelpShortDescription: cmdStruct.GetSHelp(), HelpLongDescription: cmdStruct.GetLHelp()}
}

func CreateCommandEntry(entryName string, cmd CmdFunction, helpsd string, helpld string) MenuEntry {
	return MenuEntry{EntryName: entryName, Target: cmd, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}

func CreateSubMenuEntry(entryName string, menu Menu, helpsd string, helpld string) MenuEntry {
	return MenuEntry{EntryName: entryName, Target: menu, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}
