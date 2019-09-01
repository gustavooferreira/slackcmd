package webserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/security"
)

type HandlerFunction func(context entities.RequestContext) string

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

func (scs *SlashCmdServer) RegisterCommand(cmd string, httpPath string, perm *security.Permissions, signSecret string, f HandlerFunction) {
	scs.mux.HandleFunc(httpPath, func(w http.ResponseWriter, r *http.Request) {
		slashCommand(w, r, cmd, perm, f, signSecret)
	})
}

func slashCommand(w http.ResponseWriter, r *http.Request, cmd string, perm *security.Permissions, f HandlerFunction, signSecret string) {
	// Make sure it's a POST request
	switch r.Method {
	case http.MethodPost:
		// pass
	default:
		fmt.Fprintf(w, "Only POST method is supported.")
		return
	}

	// Validate signed request
	verifier, err := slack.NewSecretsVerifier(r.Header, signSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rc := entities.RequestContext{Cmd: s.Command, TeamID: s.TeamID, ChannelID: s.ChannelID, UserID: s.UserID,
		Text: s.Text, ResponseURL: s.ResponseURL, TriggerID: s.TriggerID}

	// Validate request was sent by the correct command
	if rc.Cmd != cmd {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "This handler is not serving requests for %q command", rc.Cmd)
		return
	}

	// Check user permissions
	if perm != nil {
		if !perm.ValidateGlobal(rc.UserID) {
			if !perm.ValidateChannel(rc.ChannelID, rc.UserID) {
				fmt.Fprintf(w, "ERROR: You are not allowed to execute this command on this channel")
				return
			}
		}
	}

	// Call handler (with context of the request passed in)
	result := f(rc)
	fmt.Fprint(w, result)
}
