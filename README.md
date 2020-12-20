# mattermost-game-engine
ASCII game engine whose board state will be rendered through calls to the Mattermost API

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
