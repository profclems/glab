package gitlab

import (
	"github.com/xanzy/go-gitlab"
)

func CurrentUser(gLab *gitlab.Client) (*gitlab.User, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	u, _, err := gLab.Users.CurrentUser()
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUser(gLab *gitlab.Client, uid int) (*gitlab.User, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	u, _, err := gLab.Users.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func GetUsername(gLab *gitlab.Client, uid int) (string, error) {
	if gLab == nil {
		gLab = gLabClient
	}
	u, err := GetUser(gLab, uid)
	if err != nil {
		return "", err
	}
	return u.Username, nil
}
