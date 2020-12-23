# boredgame
ASCII game engine whose board state will be rendered through calls to the Mattermost API  

[![Build status](https://img.shields.io/appveyor/build/tston529/boredgame?style=for-the-badge)](https://ci.appveyor.com/project/tston529/boredgame/branch/main)

NOTE: Requires Go v1.13 or above, there's specific error handling in the Mattermost api package that was only released starting in that version.

installing: 
`go get -u github.com/tston529/boredgame`

To run on mattermost, define a \*.yaml file with the structure:
```yaml
user: "<username>"
pass: "<password>"
serverurl: "<mattermost server url>"
teamname: "<mattermost team>"
```

game data is also a yaml file, in the structure:
```yaml
---
map:
  filename: "path/to/map/file.txt"
  data:
    x: !!int # width of map
    y: !!int # height of map
    # more custom game-specific data is allowed here if desired
tiles:
  tile_name_1:
    ascii: !!str
    data:
      # optional; custom game-specific data goes here
  ...
  tile_name_n:
    ascii: !!str
    data:
      # optional; custom game-specific data goes here
actors:
  actor_name_1:
    ascii: !!str
    data:
      # optional; custom game-specific data goes here
  ...
  actor_name_n:
    ascii: !!str
    data:
      # optional; custom game-specific data goes here

# only necessary for emoji-based games, this defines tilesets used in message boxes.
# However, all are mandatory if emoji-based message boxes are desired.
message: 
  blank: !!str # used for filling in blank space.
  msg_vert: !!str # vertical edge of box
  msg_horiz: !!str # horizontal edge of box
  corner: !!str # corner of box (emoji equivalent of ascii art '+')

  # message will be rendered using emoji as well, make sure they
  # are labeled consistently (e.g. ":scrabble_a:, :scrabble_b:, :scrabble_c:, etc.")
  alpha_prefix: !!str 
```

Boilerplate code for getting a game to render on mattermost:
```go
import (
  "github.com/tston529/boredgame" // invoke methods from `engine`
  "github.com/tston529/boredgame/mmrender" // invoke methods from `mmrender`

  // this is less-important, util is a subpackage for dumping helper functions not directly related to the engine
  // "github.com/tston529/boredgame/util" // invoke methods from `util`
)
// Handle command line args, namely destination channel or user. The rest of the mattermost 
// credentials (username, password, url, team) are set in a yaml file read in a call to
// `LoadMattermostData(filename string) MattermostData`.
mmUser := flag.String("user", "", "The user to receive the DM of the game")
mmChannel := flag.String("channel", "", "The channel to receive the game message")

// I made this engine with custom emojis in mind. In your game's main yaml data file, you are
// to set the text to be rendered for each tile and actor, under the yaml header `ASCII`.
// These can be emoji tags (e.g. ":my_custom_emoji:")
// If you choose not to use emoji and have this flag set to false, the message will turn out 
// messy since the font used is not monospaced, so make sure this flag is set (true by default) 
// when not using emoji.
mmPreformatted := flag.Bool("pre", true, "Whether to wrap each frame in backticks to be rendered as preformatted text on Mattermost.")
flag.Parse()

if *mmUser != "" && *mmChannel != "" {
    fmt.Println("Can't specify both user and channel, choose one or the other.")
    os.Exit(1)
}

// determines what rendering strategy to use (`cli` uses escape sequences to smoothly update each frame)
var cli bool
if *mmUser != "" || *mmChannel != "" {
    cli = false
}

if !cli {
    if *mmPreformatted {
        preBeginWrap = "```\n" // Until I care enough to write better code, I defined these as globals in my tests.
        preEndWrap = "\n```"
    } else {
        preBeginWrap = ""
        preEndWrap = ""
    }
    mmData := mm_render.LoadMattermostData("./path/to/mattermost-credentials.yml")

    mm_render.StartMattermostClient(mmData.ServerUrl, mmData.User, mmData.Pass)
    if *mmUser != "" {
        mm_render.GetDirectMessageChannel(*mmUser)
    } else if *mmChannel != "" {
        mm_render.FindTeam(mmData.TeamName)
        mm_render.GetChannel(*mmChannel)
    }

  // This PostMessage is only to be called once -> it saves the post metadata under the hood
  // which it will use to update the message with the newly-generated frame.
    mm_render.PostMessage("Starting game")
}

gameData = engine.LoadGameData("./path/to/game-data.yml")

/*
   ...
   I wrote the majority of my game's loops as anonymous goroutines, I invoke them here.
   ...
*/

// Core loop; continue rendering each frame until all player's lives are lost.
for !exit {
    if cli {
        // Feel free to copy-paste this line until I write a smarter line generator.
        // Until then, change the "27" at the end to the height of your string (total lines rendered)
        fmt.Printf("\x1b[0E\x1b7%s%s\x1b[K\x1b[2G\x1b[27A", gameMap, player1.Hud())
    } else {
        // This is partly where making the preformatting wraps global comes in useful.
        mm_render.SendNextFrame(fmt.Sprintf("%s%s%s%s", preBeginWrap, gameMap, player1.Hud(), preEndWrap))
    }
    // render at 10fps. Just right for a text-based board-style game, especially if you're concerned about Mattermost rate limits. 
    time.Sleep(100 * time.Millisecond)
}
```
