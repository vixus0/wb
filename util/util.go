package util

import (
  "bufio"
  "log"
  "math"
  "os"
  "strings"
  "syscall"

  "golang.org/x/crypto/ssh/terminal"
)

func Err(msg string, err error, arg ...string) {
  if err != nil { log.Fatalln(msg, err, arg) }
}

func IsTTY() bool {
  return terminal.IsTerminal(int(os.Stdin.Fd()))
}

func IsPipe() bool {
  stat, _ := os.Stdin.Stat()
  return (stat.Mode() & os.ModeCharDevice) == 0
}

func PasswordInput() string {
  bytes, err := terminal.ReadPassword(int(syscall.Stdin))
  Err("password error:", err)
  return string(bytes)
}

func Input() string {
  reader := bufio.NewReader(os.Stdin)
  input, err := reader.ReadString('\n')
  Err("input error:", err)
  return strings.TrimSpace(input)
}

func Trunc(s string, l int) string {
  end := int( math.Min(float64(l), float64(len(s))) )
  return s[:end]
}

func B2S(bytes []byte) string {
  return strings.TrimSpace(string(bytes))
}
