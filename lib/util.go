package util

import (
  "bufio"
  "fmt"
  "math"
  "os"
  "strings"
  "syscall"

  "golang.org/x/crypto/ssh/terminal"
)

func PasswordInput() (ret string, err error) {
  defer fmt.Print("\n")
  var b []byte

  fmt.Print("Password: ")

  if b, err = terminal.ReadPassword(int(syscall.Stdin)); err == nil {
    return strings.TrimSpace(string(b)), nil
  }

  return
}

func PipeRead() (ret string, err error) {
  reader := bufio.NewReader(os.Stdin)

  if ret, err = reader.ReadString('\n'); err == nil {
    return strings.TrimSpace(ret), nil
  }

  return
}

func BytesToString(b []byte) string {
  return strings.TrimSpace(string(b))
}

func Trunc(s string, l int) string {
  end := int( math.Min(float64(l), float64(len(s))) )
  return s[:end]
}

