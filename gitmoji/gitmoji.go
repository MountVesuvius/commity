package gitmoji 

import (
	"fmt"
	"io"
	"log"
	"net/http"
    "encoding/json"
	"os"
)

const gitmojifile = "https://raw.githubusercontent.com/carloscuesta/gitmoji/master/packages/gitmojis/src/gitmojis.json"
const emojiFile = "./gitmojis.json"

// the json part at the end is for pulling from the .json file
type Gitmoji struct {
    Emoji string `json:"emoji"`
    Entity string `json:"entity"` // unicode HexNCR for the emoji
    Code string `json:"code"`
    Description string `json:"description"`
    Name string `json:"name"`
}

type Gitmojis struct {
    List []Gitmoji `json:"gitmojis"`
}

// error checking might be a bit verbose, not sure yet
func download(path string) ([]byte, error) {
    response, err := http.Get(path)
    if err != nil {
        fmt.Println("oops")
        return nil, fmt.Errorf("Unable to access file %s: %v", path, err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP Error: %v", err)
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, fmt.Errorf("Error reading file: %v", err)
    }

    return body, nil
}

// might be worth baking this into the LoadGitmojis func
func readGitmojis(content []byte) (*Gitmojis, error) {
    var gitmojis Gitmojis
    if err := json.Unmarshal(content, &gitmojis); err != nil {
        return nil, fmt.Errorf("Unable to unmarshal gitmojis: %v", err)
    }
    return &gitmojis, nil
}

// still a lot of error checking which feels verbose
func LoadGitmojis() *Gitmojis {
    // if the file doesn't exists
    if _, err := os.Stat(emojiFile); err != nil {

        content, err := download(gitmojifile)
        if err != nil {
            log.Fatalln("Download failed:", err)
        }

        os.WriteFile(emojiFile, content, 0644)
    } 
    content, err := os.ReadFile(emojiFile)
    if err != nil {
        panic(err)
    }

    gitmojis, err := readGitmojis(content)
    if err != nil {
        log.Fatalln("Parse failed:", err)
    }
    return gitmojis
}
