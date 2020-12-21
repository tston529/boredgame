package engine/mmrender

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"gopkg.in/yaml.v2"
)

var client *model.Client4
var webSocketClient *model.WebSocketClient

var myUser *model.User
var myTeam *model.Team
var myRecipient *model.User
var myChannel *model.Channel
var gamePost *model.Post
var debuggingChannel *model.Channel

// MattermostData holds login-adjacent data for Mattermost.
type MattermostData struct {
	User      string
	Pass      string
	ServerURL string
	TeamName  string
}

// LoadMattermostData reads Mattermost credentials from a yaml file
// and returns a struct containing this unmarshalled data.
func LoadMattermostData(filename string) MattermostData {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error opening file")
		os.Exit(1)
	}

	data := MattermostData{}
	err = yaml.Unmarshal(file, &data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	return data
}

// StartMattermostClient creates a client struct, logs in the user
// from given arguments, and sets up a service to gracefully shut down.
func StartMattermostClient(serverAddr string, username string, password string) {
	SetupGracefulShutdown()
	client = model.NewAPIv4Client(serverAddr)
	userLogin(username, password)
}

// SetupGracefulShutdown closes websockets on detection of
// program exiting.
func SetupGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			if webSocketClient != nil {
				webSocketClient.Close()
			}
			os.Exit(0)
		}
	}()
}

func userLogin(email string, password string) {
	if user, resp := client.Login(email, password); resp.Error != nil {
		println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		os.Exit(1)
	} else {
		myUser = user
	}
}

// FindTeam creates a request for the specified Team name and
// quietly updates the stored team if found, exiting with
// an error on failure.
func FindTeam(teamName string) {
	if team, resp := client.GetTeamByName(teamName, ""); resp.Error != nil {
		println("error: failed to find team '" + teamName + "'. Perhaps you are not part of it?")
		os.Exit(1)
	} else {
		myTeam = team
	}
}

// GetChannel creates a request for the specified Channel via
// its name and quietly updates the stored channel if found,
// exiting with an error on failure.
func GetChannel(channelName string) {
	cn := strings.ReplaceAll(channelName, " ", "-")
	if channel, resp := client.GetChannelByName(cn, myTeam.Id, ""); resp.Error != nil {
		println("Failed to get channel '" + channelName + "'")
		fmt.Printf("%v", resp)
		os.Exit(1)
	} else {
		myChannel = channel
		myRecipient = myUser
		return
	}
}

// GetDirectMessageChannel searches for a user via given username
// and quietly updates the recipient struct if found, exiting
// with an error message on failure.
func GetDirectMessageChannel(recipient string) {
	user, resp := client.GetUserByUsername(recipient, "")
	if resp.Error != nil {
		fmt.Printf("Couldn't find user with username '%s'\n", recipient)
		os.Exit(1)
	}
	myRecipient = user
	if channel, resp := client.CreateDirectChannel(myUser.Id, user.Id); resp.Error != nil {
		println("Failed to get channel with user '" + recipient + "'")
		println(resp.Error)
		os.Exit(1)
	} else {
		myChannel = channel
	}
}

// PostMessage sends a message to the stored channel,
// quietly updating the post struct with the successful
// response data, and exiting with an error message on failure.
func PostMessage(msg string) {
	newPost := model.Post{
		UserId:    myRecipient.Id,
		ChannelId: myChannel.Id,
		Message:   msg,
	}
	if post, resp := client.CreatePost(&newPost); resp.Error != nil {
		println("Failed to post the message :(")
		os.Exit(1)
	} else {
		gamePost = post
	}
}

// SendNextFrame updates the post stored via PostMessage with the
// string passed in.
func SendNextFrame(msg string) {
	gamePost.Message = msg
	if _, r := client.UpdatePost(gamePost.Id, gamePost); r.Error != nil {
		println("Failed to update post :(")
		println("%s", r.Error)
	}
}
