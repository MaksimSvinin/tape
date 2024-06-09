package tape

import (
	"io"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func mtCmd(mtcmd, dev string, args ...string) ([]byte, error) {
	cmdargs := append([]string{"-f", dev}, args...)
	cmd := exec.Command(mtcmd, cmdargs...)
	stdout, _ := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		err = errors.Wrap(err, "mt command setup stderr pipe")
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		err = errors.Wrap(err, "mt start command")
		return nil, err
	}

	cmdout, err := io.ReadAll(stdout)
	if err != nil {
		err = errors.Wrap(err, "mt read stdout output")
		return []byte{}, err
	}

	cmderr, err := io.ReadAll(stderr)
	if err != nil {
		err = errors.Wrap(err, "mt read stderr output")
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		err = errors.Wrap(err, "mt wait command")
		err = errors.Wrap(err, strings.TrimSuffix(string(cmderr), "\n"))
		return nil, err
	}
	return cmdout, nil
}
