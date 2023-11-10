# daoc-parser
Parsing Dark Age of Camelot chat.log file in order to provide user and enemy based stats 

# How to build and run
- For linux/mac
    - Install GO
        - https://go.dev/doc/install
    - ```go build```
    - copy chatgator to the folder where chat.log exists
    - ```./chatgator```
- For windows
    - Install GO
        - https://go.dev/doc/install
    - Might need to install tdm-gcc
        - https://jmeubank.github.io/tdm-gcc/
    - Set CGO_ENABLED=1
        - ```go env -w CGO_ENABLED=1```
    - Build
    - ```go build```
    - copy chatgator.exe to the folder where chat.log exists
    - run the executable
    - ```./chatgator.exe```
    - you can run as admin on the executbale itself, youll see two windows pop up

# TODO Features
- Create more regex stats
    - Pets
    - PvE?
- Create rates on the values
    - Hit rate
    - Resist rate
    - Parry rate
    - Block rate
    - Evade rate
    - Min damage
    - Max damage
    - Average damage
- Style and Spells breakdown
    - Create similar overall stats but by style and spell
- Support abilities?
    - mez, root, snare, etc
    - logs don't show this but we could parse daoc char planner, rip out all abilities that is a mez, root, snare.
    - Use a map of abilities to know ability/style types
- Window overlay
    - Convert the window to an overlay?
- Log enable / disable
    - Enable and disable logs for the user in order to flush the log buffer

- Refactor the UI and the app itself - electrode maybe