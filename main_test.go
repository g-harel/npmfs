package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockRegistry "github.com/g-harel/npmfs/internal/registry/mock"
)

func TestRoutes(t *testing.T) {
	tt := map[string]struct {
		Path     string
		Status   int
		Redirect string
	}{
		"home": {
			Path:   "/",
			Status: http.StatusOK,
		},
		"static favicon": {
			Path:     "/favicon.ico",
			Status:   http.StatusOK,
			Redirect: "/assets/favicon.ico",
		},
		"static robots.txt": {
			Path:     "/robots.txt",
			Status:   http.StatusOK,
			Redirect: "/assets/robots.txt",
		},
		"static icon": {
			Path:   "/assets/icon-package.svg",
			Status: http.StatusOK,
		},
		"package versions shortcut": {
			Path:     "/test",
			Status:   http.StatusOK,
			Redirect: "/package/test/",
		},
		"namespaced package versions shortcut": {
			Path:     "/@test/test",
			Status:   http.StatusOK,
			Redirect: "/package/@test/test/",
		},
		"package versions": {
			Path:     "/package/test",
			Status:   http.StatusOK,
			Redirect: "/package/test/",
		},
		"namespaced package versions": {
			Path:     "/package/@test/test",
			Status:   http.StatusOK,
			Redirect: "/package/@test/test/",
		},
		"package versions from compare": {
			Path:     "/compare/test",
			Status:   http.StatusOK,
			Redirect: "/package/test/",
		},
		"namespaced package versions from compare": {
			Path:     "/compare/@test/test",
			Status:   http.StatusOK,
			Redirect: "/package/@test/test/",
		},
		"package contents": {
			Path:     "/package/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/package/test/0.0.0/",
		},
		"package contents v redirect": {
			Path:     "/package/test/v/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/package/test/0.0.0/",
		},
		"namespaced package contents": {
			Path:     "/package/@test/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/package/@test/test/0.0.0/",
		},
		"namespaced package contents v redirect": {
			Path:     "/package/@test/test/v/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/package/@test/test/v/0.0.0/",
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
			Path:     "/compare/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/compare/test/0.0.0/",
		},
		"namespaced package compare picker": {
			Path:     "/compare/@test/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/compare/@test/test/0.0.0/",
		},
		"package compare": {
			Path:     "/compare/test/0.0.0/1.1.1",
			Status:   http.StatusOK,
			Redirect: "/compare/test/0.0.0/1.1.1/",
		},
		"namespaced package compare": {
			Path:     "/compare/@test/test/0.0.0/1.1.1",
			Status:   http.StatusOK,
			Redirect: "/compare/@test/test/0.0.0/1.1.1/",
		},
		"package download": {
			Path:     "/download/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/download/test/0.0.0/",
		},
		"namespaced package download": {
			Path:     "/download/@test/test/0.0.0",
			Status:   http.StatusOK,
			Redirect: "/download/@test/test/0.0.0/",
		},
		"file download": {
			Path:   "/download/test/0.0.0/test.js",
			Status: http.StatusOK,
		},
		"namespaced file download": {
			Path:   "/download/@test/test/0.0.0/test.js",
			Status: http.StatusOK,
		},
	}

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

	srv := httptest.NewServer(routes(client))
	defer srv.Close()

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%v%v", srv.URL, tc.Path))
			if err != nil {
				t.Fatalf("send GET request: %v", err)
			}

			if res.StatusCode != tc.Status {
				t.Fatalf("expected/received status codes do not match\n%v\n%v", tc.Status, res.StatusCode)
			}

			path := res.Request.URL.EscapedPath()
			if tc.Path != path {
				if tc.Redirect == "" {
					t.Fatalf("unexpected redirection\n%v\n%v", tc.Path, path)
				} else {
					if tc.Redirect != path {
						t.Fatalf("expected/received redirected paths do not match\n%v\n%v", tc.Redirect, path)
					}
				}
			}
		})
	}
}
