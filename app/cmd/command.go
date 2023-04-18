package cmd

type Command interface {
	Exec(args []string) []byte
}
