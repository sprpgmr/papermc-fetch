package paperapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const jsonResponse = `{
	"version": "v1.20.2",
	"build": 318,
	"channel": "default",
	"downloads": {
	  "application": {
		"name": "1.20.2-318.jar",
		"sha256": "asdf"
	  }
	}
  }`

func TestGetBuildInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, jsonResponse)
	}))

	defer ts.Close()

	buildInfoService := newBuildInfoServiceImpl(ts.URL)

	buildInfo, err := buildInfoService.GetBuildInfo("1.20.2", 318)
	if err != nil {
		t.Error(err)
	}

	if buildInfo.Version != "v1.20.2" {
		t.Errorf("Expected buildInfo.Version '%s' to equal v1.20.2", buildInfo.Version)
	}

	if buildInfo.Channel != "default" {
		t.Errorf("Expected buildInfo.Channel '%s' to equal default", buildInfo.Channel)
	}

	if buildInfo.Build != 318 {
		t.Errorf("Expected buildInfo.Build %d to equal 318", buildInfo.Build)
	}

	if buildInfo.Downloads == nil {
		t.Errorf("build info missing downloads info")
	}

	if buildInfo.Downloads.Application == nil {
		t.Errorf("build info missing downloads.application info")
	}

	if buildInfo.Downloads.Application.Name != "1.20.2-318.jar" {
		t.Errorf("Expected application name to be 1.20.2-318.jar but was %s", buildInfo.Downloads.Application.Name)
	}

	if buildInfo.Downloads.Application.Sha256 != "asdf" {
		t.Errorf("Expected application Sha256 to be asdf but was %s", buildInfo.Downloads.Application.Sha256)
	}
}
