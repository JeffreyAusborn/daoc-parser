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