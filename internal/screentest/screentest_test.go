// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package screentest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

func TestReadTests(t *testing.T) {
	type args struct {
		filename string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				filename: "testdata/readtests.txt",
			},
			want: []*testcase{
				{
					name:           "go.dev homepage",
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/go.dev/",
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.localhost-6060.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage.diff.png"),
					viewportWidth:  1536,
					viewportHeight: 960,
					screenshotType: fullScreenshot,
				},
				{
					name:           "go.dev homepage 540x1080",
					urlA:           "https://go.dev/",
					urlB:           "http://localhost:6060/go.dev/",
					outImgA:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.localhost-6060.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "go-dev-homepage-540x1080.diff.png"),
					viewportWidth:  540,
					viewportHeight: 1080,
					screenshotType: fullScreenshot,
				},
				{
					name:           "about page",
					urlA:           "https://go.dev/about",
					urlB:           "http://localhost:6060/go.dev/about",
					outImgA:        filepath.Join(cache, "readtests-txt", "about-page.go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "about-page.localhost-6060.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "about-page.diff.png"),
					screenshotType: fullScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
				{
					name:              "pkg.go.dev homepage .go-Carousel",
					urlA:              "https://pkg.go.dev/",
					urlB:              "https://beta.pkg.go.dev/",
					outImgA:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.pkg-go-dev.png"),
					outImgB:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.beta-pkg-go-dev.png"),
					outDiff:           filepath.Join(cache, "readtests-txt", "pkg-go-dev-homepage--go-Carousel.diff.png"),
					screenshotType:    elementScreenshot,
					screenshotElement: ".go-Carousel",
					viewportWidth:     1536,
					viewportHeight:    960,
					tasks: chromedp.Tasks{
						chromedp.Click(".go-Carousel-dot"),
					},
				},
				{
					name:           "net package doc",
					urlA:           "https://pkg.go.dev/net",
					urlB:           "https://beta.pkg.go.dev/net",
					outImgA:        filepath.Join(cache, "readtests-txt", "net-package-doc.pkg-go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "net-package-doc.beta-pkg-go-dev.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "net-package-doc.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
				},
				{
					name:           "net package doc 540x1080",
					urlA:           "https://pkg.go.dev/net",
					urlB:           "https://beta.pkg.go.dev/net",
					outImgA:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.pkg-go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.beta-pkg-go-dev.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "net-package-doc-540x1080.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  540,
					viewportHeight: 1080,
					tasks: chromedp.Tasks{
						chromedp.WaitReady(`[role="treeitem"][aria-expanded="true"]`),
					},
				},
				{
					name:           "about",
					urlA:           "https://pkg.go.dev/about",
					cacheA:         true,
					urlB:           "http://localhost:8080/about",
					outImgA:        filepath.Join(cache, "readtests-txt", "about.pkg-go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "about.localhost-8080.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "about.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
				},
				{
					name:           "eval",
					urlA:           "https://pkg.go.dev/eval",
					cacheA:         true,
					urlB:           "http://localhost:8080/eval",
					outImgA:        filepath.Join(cache, "readtests-txt", "eval.pkg-go-dev.png"),
					outImgB:        filepath.Join(cache, "readtests-txt", "eval.localhost-8080.png"),
					outDiff:        filepath.Join(cache, "readtests-txt", "eval.diff.png"),
					screenshotType: viewportScreenshot,
					viewportWidth:  1536,
					viewportHeight: 960,
					tasks: chromedp.Tasks{
						chromedp.Evaluate("console.log('Hello, world!')", nil),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTests(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("readTests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got,
				cmp.AllowUnexported(testcase{}),
				cmp.Comparer(func(a, b chromedp.ActionFunc) bool {
					return fmt.Sprint(a) == fmt.Sprint(b)
				}),
				cmp.Comparer(func(a, b chromedp.Selector) bool {
					return fmt.Sprint(a) == fmt.Sprint(b)
				}),
			); diff != "" {
				t.Errorf("readTests() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCheckHandler(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	type args struct {
		glob   string
		output string
	}
	d, err := os.UserCacheDir()
	if err != nil {
		t.Errorf("os.UserCacheDir(): %v", err)
	}
	cache := filepath.Join(d, "screentest")
	var tests = []struct {
		name      string
		args      args
		wantErr   bool
		wantFiles []string
	}{
		{
			name: "pass",
			args: args{
				glob: "testdata/pass.txt",
			},
			wantErr: false,
		},
		{
			name: "fail",
			args: args{
				output: filepath.Join(cache, "fail-txt"),
				glob:   "testdata/fail.txt",
			},
			wantErr: true,
			wantFiles: []string{
				filepath.Join(cache, "fail-txt", "homepage.diff.png"),
				filepath.Join(cache, "fail-txt", "homepage.go-dev.png"),
				filepath.Join(cache, "fail-txt", "homepage.pkg-go-dev.png"),
			},
		},
		{
			name: "cached",
			args: args{
				output: "testdata/screenshots/cached",
				glob:   "testdata/cached.txt",
			},
			wantFiles: []string{
				filepath.Join("testdata", "screenshots", "cached", "homepage.go-dev.png"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckHandler(tt.args.glob, false, nil); (err != nil) != tt.wantErr {
				t.Fatalf("CheckHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.wantFiles) != 0 {
				files, err := filepath.Glob(
					filepath.Join(tt.args.output, "*.png"))
				if err != nil {
					t.Fatal("error reading diff output")
				}
				if diff := cmp.Diff(tt.wantFiles, files); diff != "" {
					t.Errorf("readTests() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestTestHandler(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	TestHandler(t, "testdata/pass.txt", false, nil)
}

func TestHeaders(t *testing.T) {
	// Skip this test if Google Chrome is not installed.
	_, err := exec.LookPath("google-chrome")
	if err != nil {
		t.Skip()
	}
	go headerServer()
	if err := runDiff(context.Background(), &testcase{
		name:              "go.dev homepage",
		urlA:              "http://localhost:6061",
		cacheA:            true,
		urlB:              "http://localhost:6061",
		outImgA:           filepath.Join("testdata", "screenshots", "headers", "headers-test.localhost-6061.png"),
		outImgB:           filepath.Join("testdata", "screenshots", "headers", "headers-test.localhost-6061.png"),
		outDiff:           filepath.Join("testdata", "screenshots", "headers", "headers-test.diff.png"),
		viewportWidth:     1536,
		viewportHeight:    960,
		screenshotType:    elementScreenshot,
		screenshotElement: "#result",
	}, false, map[string]interface{}{"Authorization": "Bearer token"}); err != nil {
		t.Fatal(err)
	}
}

func headerServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, `<!doctype html>
		<html>
		<body>
		  <span id="result">%s</span>
		</body>
		</html>`, req.Header.Get("Authorization"))
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", 6061), mux)
}