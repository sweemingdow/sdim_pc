package unet

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"sdim_pc/backend/config"
	"sdim_pc/backend/utils/parser/json"
)

type HttpSender struct {
	hCli *http.Client
	cfg  config.HttpReqConfig
}

func NewHttpSender(cfg config.HttpReqConfig) *HttpSender {
	tp := &http.Transport{
		MaxConnsPerHost:       cfg.MaxConn,
		MaxIdleConnsPerHost:   cfg.MaxIdleConn,
		IdleConnTimeout:       cfg.ConnIdleTimeout,
		ResponseHeaderTimeout: cfg.ResponseTimeout,
	}

	return &HttpSender{
		hCli: &http.Client{
			Transport: tp,
		},
		cfg: cfg,
	}
}

func (hs *HttpSender) JsonPost(host string, body any, qp map[string]string) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), hs.cfg.RequestTimeout)
	defer cancel()

	var (
		req *http.Request
		err error
	)

	u, err := hs.encodeUrlWithParam(host, qp)
	if err != nil {
		return 0, nil, err
	}

	if body != nil {
		switch v := body.(type) {
		case []byte:
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(v))
		case string:
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader([]byte(v)))
		default:
			bodies, fe := json.Fmt(v)
			if fe != nil {
				err = fe
				return 0, nil, err
			}

			req, err = http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(bodies))
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	}

	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := hs.hCli.Do(req)
	if err != nil {
		return 0, nil, err
	}

	status := resp.StatusCode
	defer resp.Body.Close()

	respBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return status, respBuf, nil
}

func (hs *HttpSender) JsonGet(host string, qp map[string]string) (int, []byte, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), hs.cfg.RequestTimeout)
	//defer cancel()

	u, err := hs.encodeUrlWithParam(host, qp)
	if err != nil {
		return 0, nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, u, nil)

	req.Header.Set("Content-Type", "application/json")

	resp, err := hs.hCli.Do(req)
	if err != nil {
		return 0, nil, err
	}

	status := resp.StatusCode
	defer resp.Body.Close()

	respBuf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return status, respBuf, nil
}

func (hs *HttpSender) encodeUrlWithParam(host string, qp map[string]string) (string, error) {
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	if len(qp) > 0 {
		q := u.Query()

		for k, v := range qp {
			q.Set(k, v)
		}

		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
