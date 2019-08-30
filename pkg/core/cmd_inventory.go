package core

import (
	"errors"
	"fmt"
	"strings"
)

type NoCommandError struct {
	Msg string
}

func (e NoCommandError) Error() string {
	return e.Msg
}

type CommandNotFoundError struct {
	Msg string
	Cmd string
}

func (e CommandNotFoundError) Error() string {
	return e.Msg
}

type EntryNotFoundError struct {
	Msg          string
	Cmd          string
	ValidOptions []string
}

func (e EntryNotFoundError) Error() string {
	return e.Msg
}

type CommandIncompleteError struct {
	Msg string
	Cmd string
}

func (e CommandIncompleteError) Error() string {
	return e.Msg
}

type CommandInventory struct {
	AppName string
	Menu    menu
}

func NewCommandInventory(appName string, menu menu) CommandInventory {
	return CommandInventory{AppName: appName, Menu: menu}
}

func (ci CommandInventory) lookup(cmdArr []string) (me menuEntry, err error) {

	if len(cmdArr) == 0 {
		return me, &NoCommandError{Msg: "no command supplied"}
	}

	lookup := menuEntry{"", ci.Menu, "", ""}

	for _, arg := range cmdArr {
		if value, ok := lookup.Target.(menu); !ok {
			return me, &CommandNotFoundError{Msg: "command not found", Cmd: strings.Join(cmdArr, " ")}
		} else {
			if result, ok := value.menu[arg]; ok {
				lookup = result
			} else {
				return me, &EntryNotFoundError{Msg: "entry not found", Cmd: strings.Join(cmdArr, " "), ValidOptions: []string{"yolo1", "yolo2"}}
			}
		}
	}
	return lookup, nil
}

func (ci CommandInventory) HelpFunc(cmdArr []string) (helpsd string, helpld string, err error) {
	result, err := ci.lookup(cmdArr)
	if err != nil {
		return "", "", err
	}

	return result.HelpShortDescription, result.HelpLongDescription, nil
}

func (ci CommandInventory) Match(cmdArr []string) (CmdFunction, error) {
	result, err := ci.lookup(cmdArr)
	if err != nil {
		return nil, err
	}

	if _, ok := result.Target.(menu); ok {
		return nil, &CommandIncompleteError{Msg: "command incomplete", Cmd: strings.Join(cmdArr, " ")}
	}

	return result.Target.(CmdFunction), nil
}

func (ci CommandInventory) Tree(cmdArr *[]string, depth int) (string, error) {
	result := []string{}

	if depth < -1 {
		return "", errors.New("depth needs to be equal or greater than -1")
	}

	var startPoint menu

	if cmdArr == nil {
		startPoint = ci.Menu
	} else {
		lookupResult, err := ci.lookup(*cmdArr)

		if err != nil {
			return "", err
		}

		if value, ok := lookupResult.Target.(menu); ok {
			startPoint = value
		} else {
			return "", errors.New("error HERE!!")
		}
	}

	var depthStr string
	var subCmd string

	if depth == -1 {
		depthStr = ""
	} else {
		depthStr = fmt.Sprintf("[%d]", depth)
	}

	if cmdArr == nil {
		subCmd = ""
	} else {
		subCmd = fmt.Sprintf("{%s}", strings.Join(*cmdArr, "-"))
	}

	result = append(result, fmt.Sprintf("◈ /%s %s %s", ci.AppName, subCmd, depthStr))
	result = append(result, "│")

	result = append(result, recursive_tree(startPoint, depth, []bool{})...)

	// result.extend(self._recursive_tree(start_point, depth, []))

	return strings.Join(result, "\n"), nil
}

func recursive_tree(point menu, depth int, bars []bool) []string {
	result := []string{}

	length := len(point.menu)

	var lastEntryItem bool
	i := 0

	for _, entry := range point.menu {
		if i+1 == length {
			lastEntryItem = true
		} else {
			lastEntryItem = false
		}

		result = append(result, generateLine(entry.EntryName, bars, lastEntryItem))

		if value, ok := entry.Target.(menu); ok {
			d := 0

			if depth == -1 {
				d = -1
			} else {
				d--
			}

			if d != 0 {
				result = append(result, recursive_tree(value, d, append(bars, !lastEntryItem))...)
			}
		}
		i++
	}
	return result
}

func generateLine(name string, bars []bool, lastEntryItem bool) string {
	var char string

	if lastEntryItem {
		char = "└"
	} else {
		char = "├"
	}

	var msgArr []string

	for _, v := range bars {
		if v {
			msgArr = append(msgArr, "│   ")
		} else {
			msgArr = append(msgArr, "    ")
		}
	}

	return fmt.Sprintf("%s%s── %s", strings.Join(msgArr, ""), char, name)
}
