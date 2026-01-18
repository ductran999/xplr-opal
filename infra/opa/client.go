package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type OPAClient interface {
	Allow(ctx context.Context, resource string, input any) (bool, error)
}

type opaClient struct {
	BaseURL string
	Client  *http.Client
}

func NewOPAClient(baseURL string) *opaClient {
	return &opaClient{
		BaseURL: baseURL,
		Client:  http.DefaultClient,
	}
}

func (o *opaClient) Allow(ctx context.Context, resource string, input any) (bool, error) {
	body, _ := json.Marshal(map[string]any{
		"input": input,
	})

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/v1/data/%s/allow", o.BaseURL, resource),
		bytes.NewReader(body),
	)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var out struct {
		Result bool `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, err
	}

	return out.Result, nil
}
