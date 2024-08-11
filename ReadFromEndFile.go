package fileBackReader

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func ReadFromEndFile(filePath string, countLines int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}

	var args []string
	var name string = ""
	if runtime.GOOS == "windows" {
		name = "powershell"
  		args = []string{"-Command", fmt.Sprintf("Get-Content '%s' | Select-Object -Last %d", file.Name(), countLines)}
	} else if runtime.GOOS == "linux" {
		name = "tail"
		args = []string{"-n", fmt.Sprint(countLines), file.Name() }
	}
	c := exec.Command(name, args...)
    c.Stderr = os.Stderr
	data, err := c.Output()
	if err != nil {
		return []string{}, err
	}
	value := strings.Split(string(data), " ")
	return value, nil
}


func LineCount(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}

    for {
        c, err := r.Read(buf)
        count += bytes.Count(buf[:c], lineSep)

        switch {
        case err == io.EOF:
            return count, nil

        case err != nil:
            return count, err
        }
    }
}