package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockRegistry "github.com/g-harel/npmfs/internal/registry/mock"
)

func TestRoutes(t *testing.T) {
	client := &mockRegistry.Client{
		Latest: "1.1.1",
		Contents: map[string]map[string]string{
			"0.0.0": {
				"README.md": "",
			},
			"1.1.1": {
				"README.md": "",
			},
		},
	}

	tt := map[string]struct {
		Path   string
		Status int
	}{
		"home": {
			Path:   "/",
			Status: http.StatusOK,
		},
		"static favicon": {
			Path:   "/favicon.ico",
			Status: http.StatusOK,
		},
		"static robots.txt": {
			Path:   "/robots.txt",
			Status: http.StatusOK,
		},
		"static icon": {
			Path:   "/assets/icon-package.svg",
			Status: http.StatusOK,
		},
		"package versions shortcut": {
			Path:   "/test",
			Status: http.StatusOK,
		},
		"namespaced package versions shortcut": {
			Path:   "/@test/test",
			Status: http.StatusOK,
		},
		"package versions": {
			Path:   "/package/test",
			Status: http.StatusOK,
		},
		"namespaced package versions": {
			Path:   "/package/@test/test",
			Status: http.StatusOK,
		},
		"package versions from compare": {
			Path:   "/compare/test",
			Status: http.StatusOK,
		},
		"namespaced package versions from compare": {
			Path:   "/compare/@test/test",
			Status: http.StatusOK,
		},
		"package contents": {
			Path:   "/package/test/0.0.0",
			Status: http.StatusOK,
		},
		"package contents v redirect": {
			Path:   "/package/test/v/0.0.0",
			Status: http.StatusOK,
		},
		"namespaced package contents": {
			Path:   "/package/@test/test/0.0.0",
			Status: http.StatusOK,
		},
		"namespaced package contents v redirect": {
			Path:   "/package/@test/test/v/0.0.0",
			Status: http.StatusOK,
		},
		"package file": {
			Path:   "/package/test/0.0.0/README.md",
			Status: http.StatusOK,
		},
		"namespaced package file": {
			Path:   "/package/@test/test/0.0.0/README.md",
			Status: http.StatusOK,
		},
		"package compare picker": {
			Path:   "/compare/test/0.0.0",
			Status: http.StatusOK,
		},
		"namespaced package compare picker": {
			Path:   "/compare/@test/test/0.0.0",
			Status: http.StatusOK,
		},
		"package compare": {
			Path:   "/compare/test/0.0.0/1.1.1",
			Status: http.StatusOK,
		},
		"namespaced package compare": {
			Path:   "/compare/@test/test/0.0.0/1.1.1",
			Status: http.StatusOK,
		},
	}

	srv := httptest.NewServer(routes(client))
	defer srv.Close()

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%v%v", srv.URL, tc.Path))
			if err != nil {
				t.Fatalf("send GET request: %v", err)
			}

			if res.StatusCode != tc.Status {
				t.Fatalf("expected/received do not match\n%v\n%v", tc.Status, res.StatusCode)
			}
		})
	}
}
