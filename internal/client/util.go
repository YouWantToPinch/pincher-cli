package client

import (
	"fmt"
	"net/url"
	"strings"
)

// validateBaseURL ensures that the given URL
// contains a host and contains no path.
func validateBaseURL(newURL string) (u *url.URL, err error) {
	newURL = strings.TrimSuffix(strings.TrimSpace(newURL), "/")

	u, err = url.Parse(newURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	/*
		if u.Scheme != "https" {
			return nil, fmt.Errorf("base URL must use HTTPS")
		}
	*/

	if u.Path != "" {
		return nil, fmt.Errorf("base URL must not have a path (trailing /)")
	}

	if u.Host == "" {
		return nil, fmt.Errorf("base URL must have a host")
	}

	/*
		if strings.Count(u.Hostname(), ".") < 1 {
			return nil, fmt.Errorf("base URL must have a domain and TLD")
		}
	*/

	return
}
