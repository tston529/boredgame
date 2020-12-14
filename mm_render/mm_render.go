package mm_render

import (
	"os"
	"os/signal"

	"github.com/mattermost/mattermost-server/model"
)

var Client *model.Client4
var webSocketClient *model.WebSocketClient

var myUser *model.User
var myTeam *model.Team
var debuggingChannel *model.Channel

func StartMattermostClient(serverAddr string) *model.Client4 {
	return model.NewAPIv4Client(serverAddr)
}

func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}

			os.Exit(0)
		}
	}()
}

func UserLogin(email string, password string) {
	if user, resp := Client.Login(email, password); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		os.Exit(1)
	} else {
		myUser = user
	}
}

func FindTeam(teamName string) {
	if team, resp := Client.GetTeamByName(teamName, ""); resp.Error != nil {
		println("We failed to get the initial load")
		println("or we do not appear to be a member of the team '" + teamName + "'")
		os.Exit(1)
	} else {
		myTeam = team
	}
}
