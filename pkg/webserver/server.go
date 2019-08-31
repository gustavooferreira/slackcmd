package webserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nlopes/slack"

	"github.com/gustavooferreira/slackcmd/pkg/entities"
	"github.com/gustavooferreira/slackcmd/pkg/permissions"
)

type HandlerFunction func(context entities.RequestContext) string

type SlashCmdServer struct {
	mux           *http.ServeMux
	Port          uint
	SigningSecret string
}

func NewSlashCmdServer(mux *http.ServeMux, port uint, signingSecret string) SlashCmdServer {
	if mux == nil {
		mux = http.NewServeMux()
	}
	return SlashCmdServer{mux: mux, Port: port, SigningSecret: signingSecret}
}

func (scs *SlashCmdServer) ListenAndServe() {
	addr := fmt.Sprintf(":%d", scs.Port)
	http.ListenAndServe(addr, scs.mux)
}

func (scs *SlashCmdServer) RegisterCommand(cmd string, httpPath string, perm *permissions.Permissions, f HandlerFunction) {
	scs.mux.HandleFunc(httpPath, func(w http.ResponseWriter, r *http.Request) {
		slashCommand(w, r, cmd, perm, f, scs.SigningSecret)
	})
}

func slashCommand(w http.ResponseWriter, r *http.Request, cmd string, perm *permissions.Permissions, f HandlerFunction, signingSecret string) {
	// Make sure it's a POST request
	switch r.Method {
	case http.MethodPost:
		// pass
	default:
		fmt.Fprintf(w, "Only POST method is supported.")
		return
	}

	// Verify -------------------------------------------------

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
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
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	_ = s
	// -------------------------------------------------

	rc := parseForm(r)

	// debug
	fmt.Printf("%+v\n", rc)

	// Validate request was sent by slack
	// Check headers for this!
	slackSig := r.Header.Get("X-Slack-Signature")
	slackReqTS := r.Header.Get("X-Slack-Request-Timestamp")

	if slackSig == "" || slackReqTS == "" {
		fmt.Fprintf(w, "Not authorized")
		return
	}

	fmt.Println("SlackSig: ", slackSig, "SlackReqTS: ", slackReqTS)

	// Validate request was sent by the correct command
	if rc.Cmd != cmd {
		fmt.Fprintf(w, "BUUUUH")
		return
	}

	// Check user permissions (is he allowed to issue the command on the channel and on the workspace?)
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

func parseForm(r *http.Request) (rc entities.RequestContext) {
	r.ParseForm()

	reqCmd := r.Form["command"]
	if len(reqCmd) == 1 {
		rc.Cmd = reqCmd[0]
	}

	reqTeamID := r.Form["team_id"]
	if len(reqTeamID) == 1 {
		rc.TeamID = reqTeamID[0]
	}

	reqChannelID := r.Form["channel_id"]
	if len(reqChannelID) == 1 {
		rc.ChannelID = reqChannelID[0]
	}

	reqUserID := r.Form["user_id"]
	if len(reqUserID) == 1 {
		rc.UserID = reqUserID[0]
	}

	reqText := r.Form["text"]
	if len(reqText) == 1 {
		rc.Text = reqText[0]
	}

	reqResponseURL := r.Form["response_url"]
	if len(reqResponseURL) == 1 {
		rc.ResponseURL = reqResponseURL[0]
	}

	reqTriggerID := r.Form["trigger_id"]
	if len(reqTriggerID) == 1 {
		rc.TriggerID = reqTriggerID[0]
	}

	return rc
}
