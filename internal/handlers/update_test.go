package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"went/internal/utils"
)

func TestCheckLatestVersionStable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tag_name":"v2.1.0"}`))
	}))
	defer ts.Close()

	oldLatestReleaseURL := latestReleaseURL
	oldReleasesURL := releasesURL
	oldHTTPClient := httpClient
	latestReleaseURL = ts.URL
	releasesURL = ts.URL
	httpClient = ts.Client()
	defer func() {
		latestReleaseURL = oldLatestReleaseURL
		releasesURL = oldReleasesURL
		httpClient = oldHTTPClient
	}()

	result, err := CheckLatestVersion(UpdateOptions{Channel: utils.ChannelStable})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Latest != "v2.1.0" {
		t.Fatalf("expected latest stable v2.1.0, got %q", result.Latest)
	}
	if !result.UpdateAvailable {
		t.Fatal("expected update to be available")
	}
}

func TestCheckLatestVersionBeta(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
		  {"tag_name":"v2.1.0", "prerelease": false},
		  {"tag_name":"v2.1.1-beta.1", "prerelease": true},
		  {"tag_name":"v2.1.1-beta.2", "prerelease": true}
		]`))
	}))
	defer ts.Close()

	oldLatestReleaseURL := latestReleaseURL
	oldReleasesURL := releasesURL
	oldHTTPClient := httpClient
	latestReleaseURL = ts.URL
	releasesURL = ts.URL
	httpClient = ts.Client()
	defer func() {
		latestReleaseURL = oldLatestReleaseURL
		releasesURL = oldReleasesURL
		httpClient = oldHTTPClient
	}()

	result, err := CheckLatestVersion(UpdateOptions{Channel: utils.ChannelBeta})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Latest != "v2.1.1-beta.2" {
		t.Fatalf("expected latest beta v2.1.1-beta.2, got %q", result.Latest)
	}
	if !result.UpdateAvailable {
		t.Fatal("expected update to be available for beta channel")
	}
}
