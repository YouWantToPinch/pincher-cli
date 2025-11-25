package client

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) doRequest(method, url, token string, payload, out any) (*http.Response, error) {
	req, err := c.MakeRequest(method, url, token, payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	if out != nil {
		err := json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func marshal[T any](item T) ([]byte, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func decodeResponse[T any](r *http.Response) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

func unmarshal[T any](r *http.Response) (T, error) {
	var v T
	var err error

	jsonData, err := io.ReadAll(r.Body)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(jsonData, &v)
	if err != nil {
		return v, err
	}

	return v, nil
}
