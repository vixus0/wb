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
- It plays better with external tools (you can optionally pass credentials via the command-line)
- One shared session wherever you call `wbd` from

## How it works

`wb` spawns a background process `wbd` to hold on to the Bitwarden session key in memory.
Subsequent calls to `wb` that require a session key fetch it from `wbd` using Go's _net/rpc_ library.
`wbd` dies on its own after ten minutes or is asked to die if the session key is invalid for whatever reason.

[bitwarden-cli]: https://github.com/bitwarden/cli
