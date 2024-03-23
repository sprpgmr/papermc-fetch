package paperapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestSortVersions(t *testing.T) {
	input := []string{"1.20", "1.20.1", "1.19.0", "1.19.3", "1.2.2", "1.20.4"}

	sortVersions(input)

	expected := []string{"1.2.2", "1.19.0", "1.19.3", "1.20", "1.20.1", "1.20.4"}

	if !slices.Equal[[]string](input, expected) {
		t.Errorf("sort didn't sort correctly. Expected %v, got %v", expected, input)
	}
}

func TestCompareVersions(t *testing.T) {
	if compareVersions("1.2", "1.23") != -1 {
		t.Errorf("Expected 1.2 to be less than 1.23")
	}

	if compareVersions("1.20.0", "1.20") != 0 {
		t.Errorf("Expected 1.20.0 to equal 1.20")
	}

	if compareVersions("1.24", "1.23.9") != 1 {
		t.Errorf("Expected 1.24 to be greater than 1.23.9")
	}

	if compareVersions("2", "1.23.9") != 1 {
		t.Errorf("Expected 2 to be greater than 1.23.9")
	}
}

func TestGetVersions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{ \"versions\": [ \"1.2\", \"1.1\", \"1.20.4\", \"1.19.8\", \"1.19.7\", \"1.2.4\", \"1.20.2\", \"1.20.0\" ] }")
	}))

	defer ts.Close()

	versionListService := newVersionsListServiceImpl(ts.URL)

	versionList, err := versionListService.GetVersionsList()
	if err != nil {
		t.Error(err)
	}

	expectedVersions := []string{"1.1", "1.2", "1.2.4", "1.19.7", "1.19.8", "1.20.0", "1.20.2", "1.20.4"}

	if !slices.Equal[[]string](versionList.Versions, expectedVersions) {
		t.Errorf("Expected versionList %v to match %v", versionList.Versions, expectedVersions)
	}
}
