package update

import (
	"encoding/json"
	"glab/internal/request"
	"time"
)

type ReleaseInfo struct {
	Name string`json:"name"`
	PreRelease bool `json:"prerelease"`
	HTMLUrl string `json:"html_url"`
	PublishedAt time.Time `json:"published_at"`
}

// CheckForUpdate checks for latest release
func CheckForUpdate() (ReleaseInfo, error) {
	releasesUrl := "https://api.github.com/repos/profclems/glab/releases/latest"
	resp, err := request.MakeRequest("{}",releasesUrl, "GET")
	var releaseInfo ReleaseInfo
	if err != nil {
		return ReleaseInfo{}, err
	}
	json.Unmarshal([]byte(resp), &releaseInfo)
	return releaseInfo, nil
}