package rejstry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Registry struct {
	URL string
}

func (r *Registry) Fetch(name string) (*Package, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Request package info from registry.
	url := fmt.Sprintf("%s/%s", r.URL, name)
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("query registry: %v", err)
	}
	defer response.Body.Close()

	// Parse JSON response to usable data.
	pkg := &Package{}
	err = json.NewDecoder(response.Body).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return pkg, nil
}
