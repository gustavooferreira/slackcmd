package core

type menu struct {
	menu map[string]menuEntry
}

func NeWMenu() menu {
	menuMap := make(map[string]menuEntry)
	return menu{menu: menuMap}
}

func (mm *menu) AddEntry(entry menuEntry) {
	entryName := entry.EntryName
	mm.menu[entryName] = entry
}

type menuEntry struct {
	EntryName            string
	Target               interface{}
	HelpShortDescription string
	HelpLongDescription  string
}

type CmdFunction func([]string) string

func CreateCommandEntry(entryName string, cmd CmdFunction, helpsd string, helpld string) menuEntry {
	return menuEntry{EntryName: entryName, Target: cmd, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}

func CreateSubMenuEntry(entryName string, menu menu, helpsd string, helpld string) menuEntry {
	return menuEntry{EntryName: entryName, Target: menu, HelpShortDescription: helpsd, HelpLongDescription: helpld}
}
