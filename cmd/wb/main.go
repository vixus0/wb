package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/vixus0/wb/bw"
	"github.com/vixus0/wb/util"
	"github.com/vixus0/wb/server"
)

const usage = `wb <bitwarden command>

See bitwarden CLI documentation for commands.

Optional flags for non-interactive use:
  --email --password --method --code
`

type Flags struct {
	Email    *string
	Method   *string
	Code     *string
	Password *string
}

func spawnWbd(input string) {
	cmd := exec.Command("wbd")
	stdin, err := cmd.StdinPipe()
	util.Err("wbd pipe error:", err)

	err = cmd.Start()
	util.Err("wbd spawn error:", err)

	go func() {
		io.WriteString(stdin, input+"\n")
		stdin.Close()
	}()

	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		if _, err := os.Stat(server.Sock); err == nil {
			return
		}
	}

	log.Fatal("failed to find", server.Sock)
}

func getFlagOrInput(thing string, ptr *string, hide bool) {
	if len(*ptr) > 0 {
		return
	}

	if util.IsTTY() {
		log.Printf("Enter %s: ", thing)
		if hide {
			*ptr = util.PasswordInput()
		} else {
			*ptr = util.Input()
		}
		return
	}

	panic(fmt.Sprintf("Need a value for %s", thing))
}

func getSession(fl *Flags, newsession bool) (session string) {
	client, rpcerr := rpc.Dial("unix", server.Sock)

	// if we can't connect to wbd, then assume new session
	if rpcerr != nil {
		newsession = true
	}

	if newsession {
		if rpcerr == nil {
			client.Call("Server.Stop", 0, nil)
		}

		for {
			bwargs := []string{}

			if bw.IsLoggedIn() {
				// need to unlock
				log.Println("Unlock your vault")
				defer func() {
					if r := recover(); r != nil {
						log.Print(r)
						os.Exit(bw.Locked)
					}
				}()
				getFlagOrInput("password", fl.Password, true)
				bwargs = append(bwargs, "unlock", "--raw", *fl.Password)
			} else {
				// need to login
				log.Println("Login to bitwarden")
				defer func() {
					if r := recover(); r != nil {
						log.Print(r)
						os.Exit(bw.NotLoggedIn)
					}
				}()
				getFlagOrInput("email", fl.Email, false)
				getFlagOrInput("password", fl.Password, true)
				getFlagOrInput("code", fl.Code, false)
				bwargs = append(bwargs, "login", "--raw", "--method", *fl.Method, "--code", *fl.Code, *fl.Email, *fl.Password)
			}
			bytes, err := exec.Command("bw", bwargs...).Output()
			session = util.B2S(bytes)
			if err == nil {
				break
			} else {
				msg := fmt.Sprintf("Failed to %s: %s", bwargs[0], session)
				if !util.IsTTY() {
					log.Fatal(msg)
				}
				log.Print(msg)

				// reset strings
				*fl.Email = ""
				*fl.Password = ""
				*fl.Code = ""
			}
		}
		spawnWbd(session)
	} else {
		client.Call("Server.GetSession", 0, &session)
	}

	return
}

func main() {
	bw.LookPath()

	log.SetPrefix("[wb] ")
	log.SetFlags(0)

	fl := &Flags{}

	fl.Email = flag.String("email", "", "Bitwarden email")
	fl.Password = flag.String("password", "", "Bitwarden password")
	fl.Method = flag.String("method", "0", "2fa method (0 = app, 3 = yubikey)")
	fl.Code = flag.String("code", "", "TOTP code for 2fa")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Print(usage)
		os.Exit(1)
	}

	// If this is just a --help request, don't bother with anything else
	for _, f := range flag.Args() {
		if strings.Contains(f, "--help") {
			bytes, _ := exec.Command("bw", flag.Args()...).CombinedOutput()
			fmt.Print(string(bytes))
			return
		}
	}

	var (
		status int
		out    string
	)

	newsession := false

	for status = -1; ; {
		session := getSession(fl, newsession)
		status, out = bw.Cmd(session, flag.Args()...)
		if status == bw.OK || status == bw.Error {
			break
		}
		newsession = true
	}

	trimmed := strings.TrimSpace(out)

	fmt.Printf(trimmed)
	os.Exit(status)
}
