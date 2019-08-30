package webserver

import (
	"fmt"
	"net/http"
)

type RequestContext struct {
	Text        string
	ResponseURL string
	TriggerID   string
	TeamID      string
	ChannelID   string
	UserID      string
}

type HandlerFunction func(context RequestContext) string

type SlashCmdServer struct {
	mux  *http.ServeMux
	Port uint
}

func NewSlashCmdServer(mux *http.ServeMux, port uint) SlashCmdServer {
	if mux == nil {
		mux = http.NewServeMux()
	}
	return SlashCmdServer{mux: mux, Port: port}
}

func (scs *SlashCmdServer) ListenAndServe() {
	addr := fmt.Sprintf(":%d", scs.Port)
	http.ListenAndServe(addr, scs.mux)
}

func (scs *SlashCmdServer) RegisterCommand(cmd string, httpPath string, permissions Permissions, f HandlerFunction) {
	scs.mux.HandleFunc(httpPath, func(w http.ResponseWriter, r *http.Request) {
		slashCommand(w, r, cmd, f)
	})
}

func slashCommand(w http.ResponseWriter, r *http.Request, cmd string, f HandlerFunction) {

	// Make sure it's a POST request
	switch r.Method {
	case http.MethodPost:
		// pass
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
		return
	}

	r.ParseForm()

	// Validate request was sent by slack
	// Check headers for this!

	// Validate request was sent by the correct command
	reqCmd := r.Form["command"]
	if len(reqCmd) != 1 || reqCmd[0] != cmd {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	// Check user permissions (is he allowed to issue the command on the channel and on the workspace?)
	reqTeamID := r.Form["team_id"]
	if len(reqTeamID) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	reqChannelID := r.Form["channel_id"]
	if len(reqChannelID) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	reqUserID := r.Form["user_id"]
	if len(reqUserID) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	reqText := r.Form["text"]
	if len(reqText) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	reqResponseURL := r.Form["response_url"]
	if len(reqResponseURL) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	reqTriggerID := r.Form["trigger_id"]
	if len(reqTriggerID) != 1 {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	// Call handler (with context of the request passed in)
	rc := RequestContext{reqText[0], reqResponseURL[0], reqTriggerID[0], reqTeamID[0], reqChannelID[0], reqUserID[0]}

	result := f(rc)
	fmt.Fprint(w, result)
}
