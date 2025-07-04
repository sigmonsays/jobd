package job

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

// get last N lines of a file
func (me *RunSpec) GetLogs(n int) ([]string, error) {
	logfile := me.GetLogFile()
	// fh, err := os.Open(logfile)
	// if err != nil {
	// 	return nil, err
	// }
	// defer fh.Close()

	ret := make([]string, 0)
	// scanner := bufio.NewScanner(fh)
	// for scanner.Scan() {
	// 	line := scanner.Text() // Get the current line as a string
	// 	ret = append(ret, line)
	// }

	return ReadLastNLines(logfile, n)

	return ret, nil
}

// template friendly version
func (me *RunSpec) GetLastLogs(n int) (string, error) {
	logs, err := me.GetLogs(n)
	if err != nil {
		return "", err
	}
	return strings.Join(logs, "\n"), nil
}

func ReadLastNLines(filePath string, n int) ([]string, error) {
	const readBlockSize = 4096

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var (
		fileSize  = fi.Size()
		offset    = fileSize
		remaining = []byte{}
		lineCount = 0
	)

	for offset > 0 && lineCount <= n {
		blockSize := int64(readBlockSize)
		if offset < blockSize {
			blockSize = offset
		}
		offset -= blockSize

		buf := make([]byte, blockSize)
		_, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}

		remaining = append(buf, remaining...)
		lineCount = bytes.Count(remaining, []byte{'\n'})
	}

	// Use a scanner on the collected buffer to extract lines
	scanner := bufio.NewScanner(bytes.NewReader(remaining))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return only the last N lines
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines, nil
}
