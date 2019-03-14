# Witbarden

Witbarden (`wb`) is a wrapper for the [Bitwarden CLI][bitwarden-cli] (`bw`).

## Usage

```
go get github.com/vixus0/wb
wb <some bw command>
```

## Why?

- I wanted to learn some Go
- I didn't want to have to deal with my Bitwarden session key manually
- It plays better with external tools by allowing you to pass login credentials via the command-line
- There's only one shared session wherever you call `wb` from

## How it works

`wb` spawns a background process `wbd` to hold on to the Bitwarden session key in memory.
Subsequent calls to `wb` that require a session key fetch it from `wbd` using Go's _net/rpc_ library.
`wbd` dies on its own after ten minutes or is asked to die if the session key is invalid for whatever reason.

[bitwarden-cli]: https://github.com/bitwarden/cli
