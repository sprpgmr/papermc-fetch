package paperapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestGetBuildsList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{ \"version\": \"1.20.2\", \"builds\": [ 5, 4, 6, 1, 3, 2, 7 ] }")
	}))

	defer ts.Close()

	buildsListServiceImpl := newBuildsListServiceImpl(ts.URL)

	buildsList, err := buildsListServiceImpl.GetBuildsList("1.20.2")
	if err != nil {
		t.Error(err)
	}

	expectedBuilds := []int{1, 2, 3, 4, 5, 6, 7}
	if !slices.Equal[[]int](buildsList.Builds, expectedBuilds) {
		t.Errorf("Expected buildsList.Builds %v to equal %v", buildsList.Builds, expectedBuilds)
	}

	if buildsList.Version != "1.20.2" {
		t.Errorf("Expected buildsList.Version '%s' to equal %s", buildsList.Version, "1.20.2")
	}
}
