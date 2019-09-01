package inventory

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/shlex"
	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/layout"
)

// Errors -----

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

// -----------

type CommandInventory struct {
	Name   string
	Banner string
	Menu   layout.Menu
}

func NewCommandInventory(name string, banner string, menu layout.Menu) CommandInventory {
	return CommandInventory{Name: name, Banner: banner, Menu: menu}
}

func (ci CommandInventory) lookup(cmdArr []string) (me layout.MenuEntry, err error) {

	if len(cmdArr) == 0 {
		return me, &NoCommandError{Msg: "no command supplied"}
	}

	lookup := layout.MenuEntry{"", ci.Menu, "", ""}

	for _, arg := range cmdArr {
		if value, ok := lookup.Target.(layout.Menu); !ok {
			return me, &CommandNotFoundError{Msg: "command not found", Cmd: strings.Join(cmdArr, " ")}
		} else {
			if result, ok := value.Entries[arg]; ok {
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

func (ci CommandInventory) Match(cmdArr []string) (layout.CmdFunction, error) {
	result, err := ci.lookup(cmdArr)
	if err != nil {
		return nil, err
	}

	if _, ok := result.Target.(layout.Menu); ok {
		return nil, &CommandIncompleteError{Msg: "command incomplete", Cmd: strings.Join(cmdArr, " ")}
	}

	return result.Target.(layout.CmdFunction), nil
}

func (ci CommandInventory) Tree(cmdArr *[]string, depth int) (string, error) {
	result := []string{}

	if depth < -1 {
		return "", errors.New("depth needs to be equal or greater than -1")
	}

	var startPoint layout.Menu

	if cmdArr == nil {
		startPoint = ci.Menu
	} else {
		lookupResult, err := ci.lookup(*cmdArr)

		if err != nil {
			return "", err
		}

		if value, ok := lookupResult.Target.(layout.Menu); ok {
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

	result = append(result, fmt.Sprintf("◈ /%s %s %s", ci.Name, subCmd, depthStr))

	if cmdArr == nil {
		result = append(result, "│")
		result = append(result, "├── tree")
		result = append(result, "├── help")
	} else {
		result = append(result, "│")
		result = append(result, "(...)")
		result = append(result, "│")
	}

	result = append(result, recursive_tree(startPoint, depth, []bool{})...)

	return strings.Join(result, "\n"), nil
}

func recursive_tree(point layout.Menu, depth int, bars []bool) []string {
	result := []string{}

	length := len(point.Entries)

	var lastEntryItem bool
	i := 0

	for _, entry := range point.Entries {
		if i+1 == length {
			lastEntryItem = true
		} else {
			lastEntryItem = false
		}

		result = append(result, generateLine(entry.Name, bars, lastEntryItem))

		if value, ok := entry.Target.(layout.Menu); ok {
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

func (ci CommandInventory) Parse(rc entities.RequestContext) string {
	cmdArr, err := shlex.Split(rc.Text)
	if err != nil {
		return "*Error*: Problem while splitting command"
	}

	if len(cmdArr) == 0 {
		return fmt.Sprintf("```%s```", ci.Banner)
	}

	var sep int
	for i, elem := range cmdArr {
		if elem == "--" {
			sep = i
		}
	}

	var commands []string
	var options []string

	if sep == 0 {
		commands = cmdArr[:]
	} else {
		commands = cmdArr[:sep]
		options = cmdArr[sep+1:]
	}

	if commands[0] == "help" {
		return ci.helpHandler(commands[1:], options)
	} else if commands[0] == "tree" {
		return ci.treeHandler(options)
	} else {
		return ci.actionHandler(rc, commands, options)
	}
}

func (ci CommandInventory) helpHandler(helpCmd []string, options []string) string {
	if len(helpCmd) == 0 {
		return ci.Banner
	}

	_, helpld, err := ci.HelpFunc(helpCmd)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	return helpld
}

func (ci CommandInventory) treeHandler(options []string) string {
	var path *[]string
	var depth = -1
	var err error

	if len(options) == 1 || len(options) == 2 {
		pathStr := options[0]
		pathF := strings.Fields(pathStr)

		if len(pathF) != 0 {
			path = &pathF
		}
	}

	if len(options) == 2 {
		depth, err = strconv.Atoi(options[1])
		if err != nil {
			return "ERROR!!! depth needs to be an int"
		}

		if depth < 0 {
			return "ERROR!!! depth needs to be greater or equal to 0"
		}
	}

	result, _ := ci.Tree(path, depth)
	resp := fmt.Sprintf("```%s```", result)
	return resp
}

func (ci CommandInventory) actionHandler(rc entities.RequestContext, commands []string, options []string) string {
	f, err := ci.Match(commands)

	if err != nil {
		return fmt.Sprintln(err)
	} else {
		return f(rc, options)
	}
}
