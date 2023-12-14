package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrInvalidFileName = errors.New("invalid file name")

func readFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(f)
	reader := bufio.NewReader(f)
	line, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	v := bytes.ReplaceAll(line, []byte("/x00"), []byte("\n"))
	v = bytes.TrimRight(v, " \t\n")
	return string(v), nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}
	for _, info := range fileInfos {
		v, _ := readFile(path.Join(dir, info.Name()))
		fileInfo, err := info.Info()
		if err != nil {
			return env, err
		}
		if strings.Contains(fileInfo.Name(), "=") {
			return env, ErrInvalidFileName
		}
		env[info.Name()] = EnvValue{
			Value:      v,
			NeedRemove: fileInfo.Size() == 0,
		}
	}
	return env, nil
}
