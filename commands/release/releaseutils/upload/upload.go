package upload

import (
	"fmt"
	"io"

	"github.com/profclems/glab/internal/glinstance"
	"github.com/profclems/glab/pkg/iostreams"
	"github.com/xanzy/go-gitlab"
)

type ReleaseAsset struct {
	Name     *string               `json:"name,omitempty"`
	URL      *string               `json:"url,omitempty"`
	FilePath *string               `json:"filepath,omitempty"`
	LinkType *gitlab.LinkTypeValue `json:"link_type,omitempty"`
}

type ReleaseFile struct {
	Open  func() (io.ReadCloser, error)
	Name  string
	Label string
	Path  string
	Type  *gitlab.LinkTypeValue
}

func CreateLink(c *gitlab.Client, projectID, tagName string, asset *ReleaseAsset) (*gitlab.ReleaseLink, error) {
	releaseLink, _, err := c.ReleaseLinks.CreateReleaseLink(projectID, tagName, &gitlab.CreateReleaseLinkOptions{
		Name:     asset.Name,
		URL:      asset.URL,
		FilePath: asset.FilePath,
		LinkType: asset.LinkType,
	})
	if err != nil {
		return nil, err
	}
	return releaseLink, nil
}

type Context struct {
	Client      *gitlab.Client
	IO          *iostreams.IOStreams
	AssetFiles  []*ReleaseFile
	AssetsLinks []*ReleaseAsset
}

// UploadFiles uploads a file into a release repository.
func (c *Context) UploadFiles(projectID, tagName string) error {
	if c.AssetFiles == nil {
		return nil
	}
	color := c.IO.Color()
	for _, file := range c.AssetFiles {
		fmt.Fprintf(c.IO.StdErr, "%s Uploading to release\t%s=%s %s=%s\n",
			color.ProgressIcon(), color.Blue("file"), file.Path,
			color.Blue("name"), file.Name)

		projectFile, _, err := c.Client.Projects.UploadFile(
			projectID,
			file.Path,
			nil,
		)
		if err != nil {
			return err
		}

		gitlabBaseURL := fmt.Sprintf("%s://%s/", glinstance.OverridableDefaultProtocol(), glinstance.OverridableDefault())
		// projectFile.URL from upload: /uploads/<hash>/filename.txt
		linkURL := gitlabBaseURL + projectID + projectFile.URL
		filename := "/" + file.Name

		_, err = CreateLink(c.Client, projectID, tagName, &ReleaseAsset{
			Name:     &file.Label,
			URL:      &linkURL,
			FilePath: &filename,
			LinkType: file.Type,
		})

		if err != nil {
			return err
		}
	}
	c.AssetFiles = nil

	return nil
}

func (c *Context) CreateReleaseAssetLinks(projectID string, tagName string) error {
	if c.AssetsLinks == nil {
		return nil
	}
	color := c.IO.Color()
	for _, asset := range c.AssetsLinks {
		releaseLink, err := CreateLink(c.Client, projectID, tagName, asset)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.IO.StdErr, "%s Added release asset\t%s=%s %s=%s\n",
			color.GreenCheck(), color.Blue("name"), *asset.Name,
			color.Blue("url"), releaseLink.DirectAssetURL)
	}
	c.AssetsLinks = nil

	return nil
}
