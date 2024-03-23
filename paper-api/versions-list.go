package paperapi

import (
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// VersionsList contains a list of available paper versions
type VersionsList struct {
	Versions []string `json:"versions"`
}

// VersionsListService provides methods for getting a list of versions
type VersionsListService interface {
	GetVersionsList() (*VersionsList, error)
}

type versionsListServiceImpl struct {
	baseURL string
}

func newVersionsListServiceImpl(baseURL string) *versionsListServiceImpl {
	return &versionsListServiceImpl{
		baseURL: baseURL,
	}
}

// GetVersionsList will query the paper website to get a list of versions available
func (v *versionsListServiceImpl) GetVersionsList() (*VersionsList, error) {
	resp, err := http.Get(v.baseURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	dec := *json.NewDecoder(resp.Body)

	versionList := &VersionsList{}
	err = dec.Decode(versionList)

	sortVersions(versionList.Versions)

	return versionList, err
}

func sortVersions(versions []string) {
	slices.SortFunc[[]string](versions, compareVersions)
}

func compareVersions(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < max(len(aParts), len(bParts)); i++ {
		aPiece := "0"
		bPiece := "0"

		if len(aParts) >= i+1 {
			aPiece = aParts[i]
		}

		if len(bParts) >= i+1 {
			bPiece = bParts[i]
		}

		aInt, err := strconv.Atoi(aPiece)
		if err != nil {
			aInt = 0
		}

		bInt, err := strconv.Atoi(bPiece)
		if err != nil {
			bInt = 0
		}

		if aInt > bInt {
			return 1
		}

		if aInt < bInt {
			return -1
		}
	}

	return 0
}
