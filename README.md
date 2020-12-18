# mattermost-game-engine
ASCII game engine whose board state will be rendered through calls to the Mattermost API

TODO:
### Engine
- [ ] Make every board state exportable as a single string, such that it would be rendered the same in Notepad.exe as it would in a terminal.
  - [x] Render board tiles as single string
  - [ ] Remove all escape sequences
    - [ ] Move message box generator out of game logic and into engine
- [x] Move wall handling out from engine to game logic - not every ascii game will have the need for walls/impassable/passable tile separation
  - [x] Game tile/actor metadata stored in yaml file under a "Data" section allowing for game-side extensibility. This include the pacman "passable tile" trait.
- [ ] Mattermost support
  - [x] Ensure logging in/rendering works
  - [ ] Mattermost credentials via yaml file

###  Demo Pac-Man clone
- [ ] Player movement should happen automatically, arrow keys should only be used for changing directions
- [ ] Slow down player movement
- [ ] Base enemy movement logic
  - [ ] Random movement pattern
  - [ ] Movement pattern based on player's position
- [ ] Large puck logic (make enemies vulnerable, alter enemy logic)
- [ ] Respawn points
