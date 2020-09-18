# player

Play music from playlist available in `config.ini` or from the inbuilt songs list.
Add (or remove) song(s) from `config.ini` where each entry is key-value pair where *key* is the *song name* and *value* is the *song's duration*. 

## Install
* Install in `GOBIN` or `~/go/bin`:
```
go get github.com/shmsr/x/cmd/player
```
* Install manually:
```
go build
```

## Example
```sh
player  // Ensure config.ini is present on $PWD and GOBIN is present in your $PATH
        // If config.ini is not present, songs are played from the intenal playlist.
```

player expects 5 control strings from STDIN: `play`, `pause`, `prev`, `next`, `exit` with usual meanings
