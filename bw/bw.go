package bw

import (
  "log"
  "os/exec"
  "strings"

  "github.com/vixus0/wb/util"
)

const (
  OK = 0
  ERROR = 1
  LOCKED = 2
  NOT_LOGGED_IN = 3
)

func LookPath() {
  if _, err := exec.LookPath("bw"); err != nil {
    log.Fatal("Couldn't find bw in your PATH")
  }
}

func Cmd(args ...string) (status int, out string) {
  bytes, err := exec.Command("bw", args...).CombinedOutput()
  out = string(bytes)
  status = OK

  if err != nil {
    if strings.HasPrefix(out, "Session key is invalid") {
      status = LOCKED
    } else if strings.HasPrefix(out, "Vault is locked") {
      status = LOCKED
    } else if strings.HasPrefix(out, "You are not logged in") {
      status = NOT_LOGGED_IN
    } else {
      status = ERROR
    }
  }

  return
}

func Lock() {
  exec.Command("bw", "lock").Run()
}

func Unlock() (status int, out string) {
  password := readPassword()
  status, out = Cmd("unlock", "--raw", password)
  return
}

func Login() (status int, out string) {
  email := util.Input("Email: ")
  password := readPassword()
  totp := util.Input("TOTP: ")
  status, out = Cmd("login", "--raw", "--code", totp, email, password)
  return
}

func readPassword() (password string) {
  if util.IsPipe() {
    password = util.Input("Password: ")
  } else {
    password = util.PasswordInput()
  }
  return
}
