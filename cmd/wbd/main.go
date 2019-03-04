package main

import (
  "log"
  "os/exec"

  "github.com/vixus0/wb/lib"
)

var session_key string

func main() {
  if _, err := exec.LookPath("bw"); err != nil {
    log.Fatal("Couldn't find bw in your PATH")
  }

  if session_key, err := util.PipeRead(); err == nil {
    if len(session_key) == 0 {
      log.Fatal("Empty input!")
    }

    log.Printf("Got session key: %v...", util.Trunc(session_key, 8))
  }
}
