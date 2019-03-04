package main

import (
  "fmt"
  "io"
  "os"
  "os/exec"

  "github.com/vixus0/wb/lib"
)

func main() {
  if _, err := exec.LookPath("bw"); err != nil {
    fmt.Println("Couldn't find bw in your PATH")
    os.Exit(1)
  }

  var password string
  var err error

  if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
    // stdin is pipe
    password, err = util.PipeRead()
  } else {
    // stdin is tty
    password, err = util.PasswordInput()
  }

  if err != nil {
    fmt.Println("password error", err)
    os.Exit(1)
  }

  outBytes, err := exec.Command("bw", "unlock", "--raw", password).Output()
  out := util.BytesToString(outBytes)

  if err != nil {
    fmt.Println("bw error:", err, err.(*exec.ExitError).Stderr, out)
    os.Exit(1)
  }

  // We might have a session key
  fmt.Println("shit goin down")
  cmd := exec.Command("wbd")
  stdin, err := cmd.StdinPipe()

  if err != nil {
    fmt.Println("wbd error:", err)
    os.Exit(1)
  }

  go func() {
    defer stdin.Close()
    io.WriteString(stdin, out+"\n")
  }()
}

