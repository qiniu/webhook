package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func runScript(item *WatchItem) (err error) {

	script := "./" + item.Script
	out, err := exec.Command("bash", "-c", script).Output()
	if err != nil {
		log.Println("Exec command failed:", err)
	}

	log.Print("Run ", script, " output:\n", string(out))
	return
}

func handleGithub(event Payload, cfg *Config) (err error) {
	for _, item := range cfg.Items {
		if event.Repo.Url == item.Repo && event.Ref == "refs/heads/"+item.Branch {
			err = runScript(&item)
			break
		}
	}
	return
}

func handleBitbucket(event Payload, cfg *Config) (err error) {
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

func handleQuery(query url.Values, cfg *Config) (err error) {

	payload := query.Get("payload")
	b := []byte(payload)

	var payloadObj map[string]interface{}
	err = json.Unmarshal(b, &payloadObj)
	if err != nil {
		log.Println("json.Unmarshal payload failed:", err)
		return
	}

	b, _ = json.MarshalIndent(payloadObj, "", "    ")
	text := string(b)
	log.Println(text)

	var event Payload
	err = json.Unmarshal(b, &event)
	if err != nil {
		log.Println("json.Unmarshal payload failed:", err)
		return
	}

	if event.CanonUrl == "https://bitbucket.org" {
		return handleBitbucket(event, cfg)
	}

	return handleGithub(event, cfg)
}

// --------------------------------------------------------------------------------

var cfg Config

func handle(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	handleQuery(req.Form, &cfg)
}

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
