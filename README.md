# mattermost-game-engine
ASCII game engine whose board state will be rendered through calls to the Mattermost API

TODO:
- [ ] Make every board state exportable as a single string, such that it would be rendered the same in Notepad.exe as it would in a terminal.
  - [ ] Remove all escape sequences
- Move wall handling out from engine to game logic - not every ascii game will have the need for walls/impassable/passable tile separation
- [ ] Mattermost support
- [ ] Finish pac-man clone
  - [ ] Movement should happen automatically, arrow keys should only be used for changing directions
  - [ ] Slow down player movement
  - [ ] Enemy movement logic
  - [ ] Large puck logic (make enemies vulnerable, alter enemy logic)
  - [ ] Respawn points