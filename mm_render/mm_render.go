package mm_render

import (
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
)

var Client *model.Client4
var webSocketClient *model.WebSocketClient

var MyUser *model.User
var myTeam *model.Team
var Recipient *model.User
var MyChannel *model.Channel
var GamePost *model.Post
var debuggingChannel *model.Channel

type MattermostData struct {
	User      string
	Pass      string
	ServerUrl string
	TeamName  string
}

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

func StartMattermostClient(serverAddr string, username string, password string) {
	SetupGracefulShutdown()
	Client = model.NewAPIv4Client(serverAddr)
	UserLogin(username, password)
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
		MyUser = user
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

func GetChannel(channelName string) {
	cn := strings.ReplaceAll(channelName, " ", "-")
	if channel, resp := Client.GetChannelByName(cn, myTeam.Id, ""); resp.Error != nil {
		println("Failed to get channel '" + channelName + "'")
		fmt.Printf("%v", resp)
		os.Exit(1)
	} else {
		MyChannel = channel
		Recipient = MyUser
		return
	}
}

func GetDirectMessageChannel(recipient string) {
	user, resp := Client.GetUserByUsername(recipient, "")
	if resp.Error != nil {
		fmt.Printf("Couldn't find user with username '%s'\n", recipient)
		os.Exit(1)
	}
	Recipient = user
	if channel, resp := Client.CreateDirectChannel(MyUser.Id, user.Id); resp.Error != nil {
		println("Failed to get channel with user '" + recipient + "'")
		println(resp.Error)
		os.Exit(1)
	} else {
		MyChannel = channel
	}
}

func PostMessage(msg string) {
	// not sure how userId comes into play with a non-DM channel
	newPost := model.Post{
		UserId:    Recipient.Id,
		ChannelId: MyChannel.Id,
		Message:   msg,
	}
	if post, resp := Client.CreatePost(&newPost); resp.Error != nil {
		println("Failed to post the message :(")
		os.Exit(1)
	} else {
		GamePost = post
	}
}

func SendNextFrame(msg string) {
	GamePost.Message = msg
	if _, r := Client.UpdatePost(GamePost.Id, GamePost); r.Error != nil {
		println("Failed to update post :(")
		println("%s", r.Error)
	}
}
