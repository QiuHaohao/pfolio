package editor

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/google/shlex"
)

const (
	tempFilePattern = "pfolio_*.yaml"
)

type MarshalFn func(interface{}) ([]byte, error)
type UnmarshalFn func([]byte, interface{}) error

func EditYaml[T any](editor string, obj T) (res T, err error) {
	return Edit(editor, obj, yaml.Marshal, yaml.Unmarshal)
}

func EditYamlWithRetry[T any](editor string, obj T, checkFn func(T) error) (res T, err error) {
	return EditWithRetry(editor, obj, yaml.Marshal, yaml.Unmarshal, checkFn)
}

func Edit[T any](editor string, obj T, marshalFn MarshalFn, unmarshalFn UnmarshalFn) (res T, err error) {
	initialContent, err := marshalFn(obj)
	if err != nil {
		return
	}

	raw, err := Open(editor, initialContent)
	if err != nil {
		return
	}

	err = unmarshalFn(raw, &res)
	if err != nil {
		return
	}
	return
}

func EditWithRetry[T any](editor string, obj T, marshalFn MarshalFn, unmarshalFn UnmarshalFn, checkFn func(T) error) (res T, err error) {
	initialContent, err := marshalFn(obj)
	if err != nil {
		return
	}

	var raw []byte
	var retry bool
	var resToCheck T

	for {
		raw, err = Open(editor, initialContent)
		if err != nil {
			return
		}

		retry, err = withOptionToRetry(unmarshalFn(raw, &resToCheck))
		if retry {
			initialContent = raw
			continue
		} else if err != nil {
			return
		}

		retry, err = withOptionToRetry(checkFn(resToCheck))
		if retry {
			initialContent = raw
			continue
		} else if err != nil {
			return
		}

		res = resToCheck
		break
	}

	return
}

func withOptionToRetry(err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	fmt.Println("Error: ", err)
	fmt.Println("Press the 'a' key to abort, or any other key to continue editing.")
	char, _, getKeyErr := keyboard.GetSingleKey()
	if getKeyErr != nil {
		log.Fatal(getKeyErr)
	}

	if retry := char != 'a'; retry {
		return true, nil
	}

	return false, err
}

func MustOpen(editor string, initialContent []byte) []byte {
	data, err := Open(editor, initialContent)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func Open(editor string, initialContent []byte) ([]byte, error) {
	f, err := os.CreateTemp("", tempFilePattern)
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())

	tempFileName := f.Name()

	_, err = f.Write(initialContent)
	if err != nil {
		return nil, err
	}

	if err = f.Close(); err != nil {
		return nil, err
	}

	res, err := shlex.Split(editor)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, err
	}

	cmdName := res[0]
	args := append(res[1:], tempFileName)

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err = os.Open(tempFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return data, nil
}
