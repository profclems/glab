package commands

import (
	"log"

	"github.com/xanzy/go-gitlab"
	"glab/internal/git"
)

func currentUser(token string) (string, error) {
	gLab, _ := git.InitGitlabClient()
	u, _, err := gLab.Users.CurrentUser()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

func getUser(uid int) (*gitlab.User, error) {
	gLab, _ := git.InitGitlabClient()
	u, _, err := gLab.Users.GetUser(uid)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func getUsername(uid int) string {
	u, err := getUser(uid)
	if err != nil {
		log.Fatal(err)
	}
	return u.Username
}

func getUserActivities() ([]*gitlab.UserActivity, error)  {
	gLab, _ := git.InitGitlabClient()
	l := &gitlab.GetUserActivitiesOptions{}
	ua, _, err := gLab.Users.GetUserActivities(l)
	if err != nil {
		return nil, err
	}
	return ua, nil
}