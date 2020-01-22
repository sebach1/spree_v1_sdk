package spree

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Spree struct {
	http.Client

	Site string
	Key  token
}

func (s *Spree) RouteTo(path string, params url.Values, ids ...interface{}) (string, error) {
	base := s.Site + "/api"
	if ids != nil {
		base += fmt.Sprintf(path, ids...)
	} else {
		base += strings.ReplaceAll(path, "%v", "")
	}
	URL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	URL.RawQuery = params.Encode()
	base = URL.String()
	return base, nil
}

func (s *Spree) paramsWithToken() (url.Values, error) {
	if s.Key == "" {
		return nil, ErrNilKey
	}
	params := url.Values{}
	params.Set("token", string(s.Key))
	return params, nil
}

func (s *Spree) SetClient(c http.Client) {
	s.Client = c
}

func (s *Spree) SetCredentials(ctx context.Context, apiKey string) error {
	s.Key = token(apiKey)
	return s.Key.validate()
}

func (s *Spree) SetAndValidateCredentials(ctx context.Context, apiKey string) error {
	err := s.SetCredentials(ctx, apiKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *Spree) Post(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return s.Do(req)
}

func (s *Spree) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return s.Do(req)
}

func (s *Spree) Put(url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return s.Do(req)
}

func (s *Spree) Delete(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return s.Do(req)
}
