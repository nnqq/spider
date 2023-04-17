package spider

import (
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"net/url"
)

type in interface {
	GetURL() (url string, err error)
}

type out interface {
	WriteResult(url string, statusCode int, message string) error
}

type Spider struct {
	in          in
	out         out
	concurrency int
}

func NewSpider(in in, out out, concurrency int) *Spider {
	return &Spider{
		in:          in,
		out:         out,
		concurrency: concurrency,
	}
}

func (s *Spider) Run() error {
	var eg errgroup.Group
	i := 0
	for {
		if i%s.concurrency == 0 {
			err := eg.Wait()
			if err != nil {
				return fmt.Errorf("eg.Wait: %w", err)
			}
		}

		rawURL, err := s.in.GetURL()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return fmt.Errorf("s.in.GetURL: %w", err)
		}

		i += 1
		eg.Go(func() error {
			statusCode, msg := checkURL(rawURL)
			err = s.out.WriteResult(rawURL, statusCode, msg)
			if err != nil {
				return fmt.Errorf("s.out.WriteResult: %w", err)
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("eg.Wait: %w", err)
	}
	return nil
}

func checkURL(rawURL string) (statusCode int, message string) {
	u := &url.URL{}
	parsed, err := u.Parse(rawURL)
	if err != nil {
		return 0, "url address - not ok (can't parse)"
	}

	res, err := http.Get(parsed.String())
	if err != nil {
		return 0, "url address - not ok (can't make HTTP GET)"
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == http.StatusOK {
		return res.StatusCode, fmt.Sprintf("url address - ok (status %d)", res.StatusCode)
	}
	return res.StatusCode, fmt.Sprintf("url address - not ok (status %d)", res.StatusCode)
}
