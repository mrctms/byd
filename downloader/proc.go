package downloader

import (
	"os/exec"
)

type processOutput func(output string)

func runDownloadProcess(ydPath string, stdOut processOutput, stdErr processOutput, args ...string) error {
	err := runUpdate(ydPath) // for now should be ok, but maybe it is better to use a worker to check if the update is needed
	if err != nil {
		return err
	}
	cmd := exec.Command(ydPath, args...)
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	for {
		chunk := make([]byte, 1024)
		_, err := outPipe.Read(chunk)
		if err != nil {
			_, err := errPipe.Read(chunk)
			if err != nil {
				break
			}
			stdErr(string(chunk))
		}
		stdOut(string(chunk))
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func runUpdate(ydPath string) error {
	cmd := exec.Command(ydPath, "-U")
	return cmd.Run()
}
