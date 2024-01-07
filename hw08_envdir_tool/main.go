package main

import (
	"log"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatalf("Error in reading envs from dir: %s", err)
	}
	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
