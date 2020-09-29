package update

import (
	"encoding/json"
	"time"

	"github.com/profclems/glab/internal/request"
)

type ReleaseInfo struct {
	Name        string    `json:"name"`
	PreRelease  bool      `json:"prerelease"`
	HTMLUrl     string    `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

// GetUpdateInfo checks for latest glab releases
func GetUpdateInfo() (ReleaseInfo, error) {
	releasesUrl := "https://api.github.com/repos/profclems/glab/releases/latest"
	resp, err := request.MakeRequest("{}", releasesUrl, "GET")
	var releaseInfo ReleaseInfo
	if err != nil {
		return ReleaseInfo{}, err
	}
	err = json.Unmarshal([]byte(resp), &releaseInfo)
	if err != nil {
		return ReleaseInfo{}, err
	}
	return releaseInfo, nil
}
