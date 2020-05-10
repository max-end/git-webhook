package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type config struct {
	Password string
	Path     string
	LogName  string `yaml:"log-name"`
	Bind     string
}

var cfg = initConfig()

func initConfig() config {
	yamlIo, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Panic(err)
	}
	cfg := config{}
	err = yaml.Unmarshal(yamlIo, &cfg)
	if err != nil {
		log.Panic(err)
	}
	return cfg
}
func gitWebHook(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &data)
	file, err := os.OpenFile(cfg.LogName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger := log.New(file, "", log.Llongfile)
	logger.SetFlags(log.LstdFlags)
	if data["password"] != cfg.Password {
		logger.Print("password error")
		return
	}
	cmd := exec.Command("git", "pull")
	cmd.Dir = cfg.Path
	output, err := cmd.Output()
	if err != nil {
		logger.Print(err)
		return
	}
	logger.Println(string(output))
	fmt.Fprint(w, "success")
}

func main() {
	http.HandleFunc("/", gitWebHook)
	err := http.ListenAndServe(cfg.Bind, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
