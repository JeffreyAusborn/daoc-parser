# daoc-parser
Parsing Dark Age of Camelot chat.log file in order to provide user and enemy based stats 

# How to build and run
- For linux
    - ```go build```
    - ```./parser --file chat.log --stream```
- For windows
    - ```GOOS=windows GOARCH=amd64 go build -o daocParserWindows.exe```
    - ```./daocParserWindows.exe --file chat.log --stream```

# Parameters
- --file is the path to the chat.log
- --stream is a boolean flag to psuedo stream the chat logs into the parser
    - If enabling stream, you will need to disable daoc logs to flush the buffer.
    - Will read the log file every 3 seconds

# TODO Features
- Create more regex stats
    - Heals
    - Pets
    - PvE?
- Create stats on the values
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
- Window overlay
    - Create window
    - Set opacity and make sure it stays on top
    - Parse logs
    - Update values in window
- Log enable / disable
    - Enable and disable logs for the user in order to flush the log buffer
