package inventory

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/google/shlex"
	"github.com/gustavooferreira/slackcmd/pkg/entities"
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
	Name    string
	Banner  string
	Version string
	Menu    Menu
}

func NewCommandInventory(name string, banner string, version string, menu Menu) CommandInventory {
	return CommandInventory{Name: name, Banner: banner, Version: version, Menu: menu}
}

func (ci CommandInventory) lookup(cmdArr []string) (me menuEntry, err error) {

	if len(cmdArr) == 0 {
		return me, &NoCommandError{Msg: "no command supplied"}
	}

	lookup := menuEntry{Name: "", HelpShortDescription: "", HelpLongDescription: "", Type: menuEntryType_SubMenu, SubMenu: &ci.Menu}

	for _, arg := range cmdArr {
		if lookup.Type == menuEntryType_SubMenu {
			if result, ok := lookup.SubMenu.Entries[arg]; ok {
				lookup = result
			} else {
				return me, &EntryNotFoundError{Msg: "entry not found", Cmd: strings.Join(cmdArr, " "), ValidOptions: []string{"yolo1", "yolo2"}}
			}
		} else {
			return me, &CommandNotFoundError{Msg: "command not found", Cmd: strings.Join(cmdArr, " ")}
		}
	}
	return lookup, nil
}

func (ci CommandInventory) match(cmdArr []string) (CmdFunction, error) {
	result, err := ci.lookup(cmdArr)
	if err != nil {
		return nil, err
	}

	if result.Type == menuEntryType_SubMenu {
		return nil, &CommandIncompleteError{Msg: "command incomplete", Cmd: strings.Join(cmdArr, " ")}
	}

	return result.Cmd, nil
}

// Tree returns a drawing of the menu
// cmdArr is the starting point and depth is how deep the recursion should go
func (ci CommandInventory) tree(cmdArr *[]string, depth int) (string, error) {
	result := []string{}

	if depth < -1 {
		return "", errors.New("depth needs to be greater or equal to -1")
	}

	var startPoint Menu

	if cmdArr == nil {
		startPoint = ci.Menu
	} else {
		lookupResult, err := ci.lookup(*cmdArr)
		if err != nil {
			return "", err
		}

		if lookupResult.Type != menuEntryType_SubMenu {
			return "", errors.New("starting point should be a menu")
		}
		startPoint = *lookupResult.SubMenu
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
		result = append(result, "├── version")
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

func recursive_tree(point Menu, depth int, bars []bool) []string {
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

		if entry.Type == menuEntryType_SubMenu {
			d := 0

			if depth == -1 {
				d = -1
			} else {
				d--
			}

			if d != 0 {
				result = append(result, recursive_tree(*entry.SubMenu, d, append(bars, !lastEntryItem))...)
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

func (ci CommandInventory) Parse(rc entities.RequestContext, resp io.Writer) {
	cmdArr, err := shlex.Split(rc.Text)
	if err != nil {
		fmt.Fprintf(resp, "*Error*: Problem while splitting command")
		return
	}

	if len(cmdArr) == 0 {
		fmt.Fprintf(resp, "```%s```", ci.Banner)
		return
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
		ci.handlerHelp(commands[1:], options, resp)
	} else if commands[0] == "version" {
		ci.handlerVersion(resp)
	} else if commands[0] == "tree" {
		ci.handlerTree(options, resp)
	} else {
		ci.handlerAction(rc, commands, options, resp)
	}
}

// TODO: check if command is "help tree" show help stuff for tree command!
// Do the same for help version
func (ci CommandInventory) handlerHelp(helpCmd []string, options []string, resp io.Writer) {
	// TODO: show help for the help command instead
	if len(helpCmd) == 0 {
		fmt.Fprintf(resp, "```%s```", ci.Banner)
		return
	}

	result, err := ci.lookup(helpCmd)
	if err != nil {
		HandlerError(resp, err.Error())
		return
	}

	cmd := strings.Join(helpCmd, " - ")
	helpld := result.HelpLongDescription

	broadcast := false

	for _, opt := range options {
		optionsArr, err := strings.Split(opt, "=")
		if err != nil {
			HandlerError(resp, "error while splitting options")
			return
		}
		if len(optionsArr) != 2 {
			HandlerError(resp, "error while splitting options, too many '='")
			return
		}

		if optionsArr[0] == "broadcast" && optionsArr[1] == "all" {
			broadcast = true
		}
	}

	respType := "ephemeral"
	if broadcast {
		respType = "in_channel"
	}

	msg := `{
    "response_type": "%s",
	"blocks": [
		{"type": "divider"},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "> :grey_question: *Path:* %s\n>\n>%s"
			}
		},
        {"type": "divider"}
	]}`
	fmt.Fprintf(resp, msg, respType, cmd, helpld)
}

func (ci CommandInventory) handlerVersion(resp io.Writer) {
	msg := `{
    "response_type": "ephemeral",
	"blocks": [
		{"type": "divider"},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "> :package: *Version:* %s"
			}
		},
        {"type": "divider"}
	]}`

	fmt.Fprintf(resp, msg, ci.Version)
}

func (ci CommandInventory) handlerTree(options []string, resp io.Writer) {
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
			fmt.Fprint(resp, "ERROR!!! depth needs to be an int")
			return
		}

		if depth < 0 {
			fmt.Fprintf(resp, "ERROR!!! depth needs to be greater or equal to 0")
			return
		}
	}

	result, err := ci.tree(path, depth)
	if err != nil {
		HandlerError(resp, err.Error())
		return
	}

	respType := "ephemeral"
	if 1 == 1 {
		respType = "in_channel"
	}

	msg := `{
    "response_type": "%s",
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "%s"
			}
		}
	]}`
	fmt.Fprintf(resp, msg, respType, fmt.Sprintf("```%s```", result))
}

// TODO: Cmd might be nil, careful with that!
func (ci CommandInventory) handlerAction(rc entities.RequestContext, commands []string, options []string, resp io.Writer) {
	f, err := ci.match(commands)
	if err != nil {
		HandlerError(resp, err.Error())
		return
	}

	f(rc, options, resp)
}

func HandlerError(resp io.Writer, errorMsg string) {
	msg := `{
    "response_type": "ephemeral",
	"blocks": [
		{"type": "divider"},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "> :bangbang: *ERROR:* %s"
			}
		},
        {"type": "divider"}
	]}`

	fmt.Fprintf(resp, msg, errorMsg)
}
