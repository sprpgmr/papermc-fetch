package paperapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// BuildInfo contains information about a specific paper build.
type BuildInfo struct {
	Version   string        `json:"version"`
	Channel   string        `json:"channel"`
	Downloads *DownloadInfo `json:"downloads"`
	Build     int           `json:"build"`
}

// DownloadInfo contains information about the available downloads for a Build
type DownloadInfo struct {
	Application *ApplicationInfo `json:"application"`
}

// ApplicationInfo contains information about the available download, including file name, and sha256 hash
type ApplicationInfo struct {
	Name   string `json:"name"`
	Sha256 string `json:"sha256"`
}

// BuildInfoService provides methods for getting build info
type BuildInfoService interface {
	GetBuildInfo(version string, build int) (*BuildInfo, error)
}

type buildInfoServiceImpl struct {
	baseURL string
}

func newBuildInfoServiceImpl(baseURL string) *buildInfoServiceImpl {
	return &buildInfoServiceImpl{
		baseURL: baseURL,
	}
}

// GetBuildInfo will get the build info for a specific version and build number
func (s *buildInfoServiceImpl) GetBuildInfo(version string, build int) (*BuildInfo, error) {
	if len(version) == 0 {
		return nil, errors.New("version must be specified to get build info")
	}

	url := fmt.Sprint(s.baseURL, "/versions/", version, "/builds/", build)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	buildInfo := &BuildInfo{}
	err = dec.Decode(buildInfo)

	return buildInfo, err
}
