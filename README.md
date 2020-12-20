# mattermost-game-engine
ASCII game engine whose board state will be rendered through calls to the Mattermost API

NOTE: Requires Go v1.13 or above, there's specific error handling in the Mattermost api package that was only released starting in that version.

Demo included -> currently called `mm_game.go`, it's basically a half-baked Pac-Man clone.
There's a bunch to learn from it though, including what I'd consider boilerplate code for getting a game to render on mattermost:
```go
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

if *mmUser != "" || *mmChannel != "" {
    cli = false
}

if !cli {
    if *mmPreformatted {
        preBeginWrap = "```\n" // Until I care enough to write better code, I defined these as globals.
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


TODO:
### Engine
- [x] Make every board state exportable as a single string, such that it would be rendered the same in Notepad.exe as it would in a terminal.
  - [x] Render board tiles as single string
  - [x] Remove the need for escape sequences
    - [x] Move message box generator out of game logic and into engine
    - [x] Message box should be inserted into and rendered as part of current frame
    - [x] Separate the mattermost renderer from the cli renderer (which may still use escape sequences)
- [x] Move wall handling out from engine to game logic - not every ascii game will have the need for walls/impassable/passable tile separation
  - [x] Game tile/actor metadata stored in yaml file under a "Data" section allowing for game-side extensibility. This include the pacman "passable tile" trait.
- [x] Mattermost support
  - [x] Ensure logging in/basic rendering works
  - [x] Mattermost credentials via yaml file
  - [x] Move all mattermost handling into mattermost renderer package
  - [x] Handle Channels
  - [x] Handle Direct Messages
- [ ] Proper Go package handling

###  Demo Pac-Man clone
- [ ] Player movement should happen automatically, arrow keys should only be used for changing directions
- [x] Slow down player movement
- [ ] Base enemy movement logic
  - [x] Random movement pattern
  - [ ] Movement pattern based on player's position
- [x] Large puck logic (points, make enemies vulnerable, alter enemy logic)
- [ ] Respawn points
  - [x] Enemy respawn after being eaten
  - [ ] Player respawn after being touched by enemy
