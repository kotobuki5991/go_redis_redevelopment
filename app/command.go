package main

type Command interface {
	Exec(args []string) []byte
}
