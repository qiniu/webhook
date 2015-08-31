package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// --------------------------------------------------------------------------------

type WatchItem struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Script string `json:"script"`
}

type Config struct {
	BindHost string      `json:"bind"`
	Items    []WatchItem `json:"items"`
}

// --------------------------------------------------------------------------------

type Repository struct {
	Url         string `json:"url"` // "https://github.com/qiniu/api"
	AbsoluteUrl string `json:"absolute_url"`
}

type Commit struct {
	Branch string `json:"branch"`
}

type Payload struct {
	Ref      string     `json:"ref"` // "refs/heads/develop"
	Repo     Repository `json:"repository"`
	CanonUrl string     `json:"canon_url"`
	Commits  []Commit   `json:"commits"`
}

// --------------------------------------------------------------------------------

var cfg Config

// --------------------------------------------------------------------------------

func runScript(item *WatchItem) (err error) {
	script := "./" + item.Script
	out, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		log.Printf("Exec command failed: %s\n", err)
	}

	log.Printf("Run %s output: %s\n", script, string(out))
	return
}

func handleGithub(event Payload, cfg *Config) (err error) {
	for _, item := range cfg.Items {
		if event.Repo.Url == item.Repo && strings.Contains(event.Ref, item.Branch) {
			err = runScript(&item)
			if err != nil {
				log.Printf("run script error: %s\n", err)
			}
			break
		}
	}
	return
}

func handleBitbucket(event Payload, cfg *Config) {
	changingBranches := make(map[string]bool)

	for _, commit := range event.Commits {
		changingBranches[commit.Branch] = true
	}

	repo := strings.TrimRight(event.CanonUrl+event.Repo.AbsoluteUrl, "/")

	for _, item := range cfg.Items {
		if strings.TrimRight(item.Repo, "/") == repo && changingBranches[item.Branch] {
			runScript(&item)
		}
	}
	return
}

func handle(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)

	var event Payload
	err := decoder.Decode(&event)
	if err != nil {
		log.Printf("payload json decode failed: %s\n", err)
		return
	}

	if event.CanonUrl == "https://bitbucket.org" {
		handleBitbucket(event, &cfg)
		return
	}

	handleGithub(event, &cfg)
}

// --------------------------------------------------------------------------------

func main() {

	if len(os.Args) < 2 {
		println("Usage: webhook <ConfigFile>\n")
		return
	}

	cfgbuf, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Println("Read config file failed:", err)
		return
	}

	err = json.Unmarshal(cfgbuf, &cfg)
	if err != nil {
		log.Println("Unmarshal config failed:", err)
		return
	}

	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(cfg.BindHost, nil))
}

// --------------------------------------------------------------------------------
