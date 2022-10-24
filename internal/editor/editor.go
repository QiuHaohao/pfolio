package editor

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"

	"gopkg.in/yaml.v3"

	"github.com/google/shlex"
)

const (
	tempFilePattern = "pfolio_*.yaml"
)

type MarshalFn func(interface{}) ([]byte, error)
type UnmarshalFn func([]byte, interface{}) error

type config struct {
	editor string
	stdin  *os.File
	stdout *os.File
	stderr *os.File
}

func newDefaultConfig() *config {
	return &config{
		editor: os.Getenv("EDITOR"),
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

type Config func(*config)

func WithEditor(editor string) func(*config) {
	return func(c *config) {
		c.editor = editor
	}
}

func WithStdin(f *os.File) func(*config) {
	return func(c *config) {
		c.stdin = f
	}
}

func WithStdout(f *os.File) func(*config) {
	return func(c *config) {
		c.stdout = f
	}
}

func WithStderr(f *os.File) func(*config) {
	return func(c *config) {
		c.stderr = f
	}
}

func apply(cfg *config, fns ...Config) {
	for _, fn := range fns {
		fn(cfg)
	}
}

type Editor interface {
	Edit(interface{}) (interface{}, error)
}

type RetryEditor interface {
	EditWithRetry(interface{}, func(interface{}) error) (interface{}, error)
}

func NewObjectEditor(marshalFn MarshalFn, unmarshalFn UnmarshalFn, configs ...Config) Editor {
	return objEditor{
		marshalFn:   marshalFn,
		unmarshalFn: unmarshalFn,
		configs:     configs,
	}
}

type objEditor struct {
	marshalFn   MarshalFn
	unmarshalFn UnmarshalFn
	configs     []Config
}

func (e objEditor) Edit(obj interface{}) (res interface{}, err error) {
	return Edit(obj, e.marshalFn, e.unmarshalFn, e.configs...)
}

type objEditorWithRetry struct {
	marshalFn   MarshalFn
	unmarshalFn UnmarshalFn
	checkFn     func(interface{}) error
	configs     []Config
}

func NewObjEditorWithRetry(marshalFn MarshalFn, unmarshalFn UnmarshalFn, configs ...Config) RetryEditor {
	return &objEditorWithRetry{
		marshalFn:   marshalFn,
		unmarshalFn: unmarshalFn,
		configs:     configs,
	}
}

func NewYamlObjEditorWithRetry(configs ...Config) RetryEditor {
	return NewObjEditorWithRetry(yaml.Marshal, yaml.Unmarshal, configs...)
}

func (e objEditorWithRetry) EditWithRetry(obj interface{}, checkFn func(interface{}) error) (interface{}, error) {
	return EditWithRetry(obj, e.marshalFn, e.unmarshalFn, checkFn, e.configs...)
}

func EditYaml(obj interface{}, configs ...Config) (res interface{}, err error) {
	return Edit(obj, yaml.Marshal, yaml.Unmarshal)
}

func EditYamlWithRetry(obj interface{}, checkFn func(interface{}) error, configs ...Config) (res interface{}, err error) {
	return EditWithRetry(obj, yaml.Marshal, yaml.Unmarshal, checkFn)
}

func Edit(obj interface{}, marshalFn MarshalFn, unmarshalFn UnmarshalFn, configs ...Config) (res interface{}, err error) {
	initialContent, err := marshalFn(obj)
	if err != nil {
		return
	}

	raw, err := Open(initialContent, configs...)
	if err != nil {
		return
	}

	err = unmarshalFn(raw, &res)
	if err != nil {
		return
	}
	return
}

func EditWithRetry(obj interface{}, marshalFn MarshalFn, unmarshalFn UnmarshalFn, checkFn func(interface{}) error, configs ...Config) (res interface{}, err error) {
	cfg := newDefaultConfig()
	apply(cfg, configs...)

	initialContent, err := marshalFn(obj)
	if err != nil {
		return
	}

	var raw []byte
	var retry bool
	objType := reflect.ValueOf(obj).Type()
	resToCheckPtr := reflect.New(objType).Interface()

	for {
		raw, err = Open(initialContent, configs...)
		if err != nil {
			return
		}

		retry, err = withOptionToRetry(unmarshalFn(raw, resToCheckPtr), cfg.stdin, cfg.stdout)
		if retry {
			initialContent = raw
			continue
		} else if err != nil {
			return
		}

		resToCheck := reflect.ValueOf(resToCheckPtr).Elem().Interface()
		retry, err = withOptionToRetry(checkFn(resToCheck), cfg.stdin, cfg.stdout)
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

func withOptionToRetry(err error, stdin *os.File, stdout *os.File) (bool, error) {
	if err == nil {
		return false, nil
	}

	fmt.Fprintln(stdout, "Error: ", err)
	fmt.Fprintln(stdout, "Input 'a' to abort, or anything else to continue editing:")

	var input string

	if _, scanErr := fmt.Fscanln(stdin, &input); scanErr != nil {
		log.Fatal(scanErr)
	}

	if retry := input != "a"; retry {
		return true, nil
	}

	return false, err
}

func MustOpen(initialContent []byte, configs ...Config) []byte {
	data, err := Open(initialContent, configs...)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func Open(initialContent []byte, configs ...Config) ([]byte, error) {
	cfg := newDefaultConfig()
	apply(cfg, configs...)

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

	res, err := shlex.Split(cfg.editor)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, err
	}

	cmdName := res[0]
	args := append(res[1:], tempFileName)

	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = cfg.stdin
	cmd.Stdout = cfg.stdout
	cmd.Stderr = cfg.stderr

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
