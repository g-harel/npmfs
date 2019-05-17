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

type Package struct {
	Tags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]struct {
		Dist struct {
			Tarball string `json:"tarball"`
		} `json:"dist"`
	} `json:"versions"`
}

func (r *Registry) Fetch(name string) (*Package, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("%s/%s", r.URL, name)
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("query registry: %v", err)
	}
	defer response.Body.Close()

	pkg := &Package{}
	err = json.NewDecoder(response.Body).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	return pkg, nil
}
