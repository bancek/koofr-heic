package movtomp4

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

func MovToMp4(movPath string, mp4Path string) (err error) {
	cmd := exec.Command("ffmpeg", "-i", movPath, mp4Path)

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
		return fmt.Errorf("MovToMp4 convert error: %s: %s", err, string(stderrBytes))
	}

	return nil
}
