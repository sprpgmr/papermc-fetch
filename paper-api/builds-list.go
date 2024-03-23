package paperapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
)

// BuildsList contains information about the available builds made for the version
type BuildsList struct {
	Version string `json:"version"`
	Builds  []int  `json:"builds"`
}

// BuildsListService provides methods for getting a list of builds
type BuildsListService interface {
	GetBuildsList(version string) (*BuildsList, error)
}

type buildsListServiceImpl struct {
	baseURL string
}

func newBuildsListServiceImpl(baseURL string) *buildsListServiceImpl {
	return &buildsListServiceImpl{
		baseURL: baseURL,
	}
}

// GetBuildsList gets a list of builds for the version provided
func (s *buildsListServiceImpl) GetBuildsList(version string) (*BuildsList, error) {
	if len(version) == 0 {
		return nil, errors.New("must specify version to get builds for")
	}

	buildsURL := s.baseURL + "/versions/" + version

	resp, err := http.Get(buildsURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	buildsList := &BuildsList{}

	err = dec.Decode(buildsList)

	slices.Sort[[]int](buildsList.Builds)

	return buildsList, err
}
