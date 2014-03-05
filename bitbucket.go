package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

// --------------------------------------------------------------------------------

type WatchItem struct {
	Repo   string `json:"repo"`
	Branch string `json:"branch"`
	Script string `json:"script"`
}

type Commit struct {
	Branch string `json:"branch"`
}

type Config struct {
	BindHost string      `json:"bind"`
	Items    []WatchItem `json:"items"`
}

// --------------------------------------------------------------------------------

type Repository struct {
	Url string `json:"absolute_url"` // "/Wuvist/angelbot"
}

type Payload struct {
	Ref     string     `json:"ref"` // "refs/heads/develop"
	Repo    Repository `json:"repository"`
	Commits []Commit   `json:"commits"`
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

	changingBranches := make(map[string]bool)

	for _, commit := range event.Commits {
		changingBranches[event.Repo.Url+"::"+commit.Branch] = true
	}

	for _, item := range cfg.Items {
		if changingBranches[item.Repo+"::"+item.Branch] {
			runScript(&item)
		}
	}
	return
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
