package heictojpg

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func HeicToJpg(heicPath string, jpgPath string) (err error) {
	cmd := exec.Command("convert", heicPath, jpgPath)

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stderrBytes, err := ioutil.ReadAll(stderr)

	if err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("HeicToJpg convert error: %s: %s", err, string(stderrBytes))
	}

	return nil
}
