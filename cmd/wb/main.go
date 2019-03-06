package main

import (
  "fmt"
  "io"
  "log"
  "net/rpc"
  "os"
  "os/exec"
  "time"

  "github.com/vixus0/wb/bw"
  "github.com/vixus0/wb/util"
  "github.com/vixus0/wb/wbd"
)

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
    if _, err := os.Stat(wbd.Sock); err == nil {
      return
    }
  }

  log.Fatal("failed to find", wbd.Sock)
}

func connectWbd() (client *rpc.Client, err error) {
  client, err = rpc.Dial("unix", wbd.Sock)
  return
}

func runCmd(oldSession string, args []string) (session string, out string) {
  session = oldSession

  bwArgs := []string{"--session", session}
  bwArgs = append(bwArgs, args...)

  status := -1

  for ; status != bw.OK; status, out = bw.Cmd(bwArgs...) {
    switch status {
    case bw.LOCKED:
      log.Println("Unlock")
      status, session = bw.Unlock()
    case bw.NOT_LOGGED_IN:
      log.Println("Login")
      status, session = bw.Login()
    }
    if status == bw.ERROR {
      log.Fatalln("bw:", session)
    }
    bwArgs[1] = session
  }

  return
}

func main() {
  bw.LookPath()

  log.SetPrefix("[wb] ")
  log.SetFlags(0)

  // args
  if len(os.Args) == 1 {
    fmt.Println("Usage: wb <bitwarden command>")
    fmt.Println("  See bitwarden cli usage for details")
    os.Exit(1)
  }

  var session, out string
  client, err := connectWbd()

  if err != nil {
    session, out = runCmd("", os.Args[1:])
    spawnWbd(session)
  } else {
    var wbdSession string
    client.Call("Server.GetSession", 0, &wbdSession)
    session, out = runCmd(wbdSession, os.Args[1:])

    if session != wbdSession {
      client.Call("Server.Stop", 0, nil)
      spawnWbd(session)
    }
  }

  fmt.Print(out)
}
