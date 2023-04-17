package list

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Reader struct {
	ch    <-chan string
	errCh <-chan error
}

func NewReader(path string) *Reader {
	ch, errCh := lineByLine(path)
	return &Reader{
		ch:    ch,
		errCh: errCh,
	}
}

func (r *Reader) GetURL() (string, error) {
	select {
	case line := <-r.ch:
		return line, nil
	case err := <-r.errCh:
		return "", err
	}
}

func lineByLine(path string) (<-chan string, <-chan error) {
	ch, errCh := make(chan string), make(chan error)

	file, err := os.Open(path)
	if err != nil {
		errCh <- fmt.Errorf("os.Open: %w", err)
		return ch, errCh
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	go func() {
		defer func() {
			_ = file.Close()
			close(ch)
			close(errCh)
		}()

		for scanner.Scan() {
			ch <- scanner.Text()
		}
		errCh <- io.EOF
	}()
	return ch, errCh
}
