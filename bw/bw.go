package bw

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/vixus0/wb/util"
)

const (
	OK          = 0
	Error       = 1
	Locked      = 2
	NotLoggedIn = 3
)

var bwFileLocations = []string{
	fmt.Sprintf("%s/.config/Bitwarden CLI/data.json", os.Getenv("HOME")),
	fmt.Sprintf("%s/Library/Application Support/Bitwarden CLI/data.json", os.Getenv("HOME")),
}

func Cmd(session string, args ...string) (status int, out string) {
	args = append(args, "--session", session)
	bytes, err := exec.Command("bw", args...).CombinedOutput()
	out = string(bytes)
	status = OK

	if err != nil {
		if strings.Contains(out, "Vault is locked") {
			status = Locked
		} else if strings.Contains(out, "You are not logged in") {
			status = NotLoggedIn
		} else {
			status = Error
		}

		if util.IsTTY() == false {
			log.Print(out)
			os.Exit(status)
		}
	}

	return
}

func LookPath() {
	if _, err := exec.LookPath("bw"); err != nil {
		log.Fatal("Couldn't find bw in your PATH")
	}
}

func IsLoggedIn() bool {
	var data map[string]interface{}
	var err error

	for _, bwFile := range bwFileLocations {
		if bytes, err := ioutil.ReadFile(bwFile); err == nil {
			if err = json.Unmarshal(bytes, &data); err == nil {
				for k := range data {
					if k == "userId" {
						return true
					}
				}
			}
		}
	}

	log.Print(err)
	return false
}
