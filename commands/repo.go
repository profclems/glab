package commands

import (
	"bytes"
	"fmt"
	"glab/cmd/glab/utils"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/xanzy/go-gitlab"
)

func getRepoContributors() {
	MakeRequest(`{}`, "projects/20131402/issues/1", "GET")
}

func newBranch() {

}

//CloneWriter w
type CloneWriter struct {
	Total uint64
}

func (wc *CloneWriter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.progress()
	return n, nil
}

func (wc CloneWriter) progress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rCloning... %s complete", humanize.Bytes(wc.Total))
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func archiveRepo(repository string, format string, name string) {

	// tar.gz, tar.bz2, tbz, tbz2, tb2, bz2, tar, and zip
	extensions := []string{"tar.gz", "tar.bz2", "tbz", "tbz2", "tb2", "bz2", "tar", "zip"}
	if b := contains(extensions, format); !b {

		fmt.Println("fatal: --format must be one of " + strings.Join(extensions[:], ","))

		return
	}

	git, _ := InitGitlabClient()
	l := &gitlab.ArchiveOptions{}
	l.Format = gitlab.String(format)
	ext := *l.Format
	archiveName := strings.Replace(GetEnv("GITLAB_REPO"), "/", "-", -1)
	if len(strings.TrimSpace(name)) != 0 {
		archiveName = name + "." + ext
	}

	bt, _, err := git.Repositories.Archive(GetEnv("GITLAB_PROJECT_ID"), l)
	if err != nil {
		log.Fatalf("Failed to clone repos: %v", err)
	}

	r := bytes.NewReader(bt)
	out, err := os.Create(archiveName + ".tmp")
	if err != nil {

		log.Fatalf("Failed to create archive repos: %v", err)
	}

	counter := &CloneWriter{}
	if _, err = io.Copy(out, io.TeeReader(r, counter)); err != nil {
		out.Close()
		log.Fatalf("Failed to write repos: %v", err)
	}

	fmt.Print("\n")
	out.Close()

	if err = os.Rename(archiveName+".tmp", archiveName); err != nil {
		log.Fatalf("Failed to rename tmp repos: %v", err)
	}
	fmt.Println("finish ", archiveName)
}

func cloneRepo(cmdOptions map[string]string, cmdArgs map[int]string) {
	repo := cmdArgs[1]
	dir := cmdArgs[2]

	if len(strings.TrimSpace(repo)) == 0 || !strings.Contains(repo, "/") {

		fmt.Println("fatal: You must specify a owner/repository to clone.")

		return
	}
	if CommandArgExists(cmdOptions, "format") {
		archiveRepo(repo, cmdOptions["format"], dir)

		return

	}

	repos := strings.Split(repo, "/")
	u, _ := url.Parse(GetEnv("GITLAB_URI"))
	url := GetEnv("GITLAB_URI") + ":" + GetEnv("GITLAB_TOKEN") + "@" + u.Host + "/" + repos[0] + "/" + repos[1]
	// git clone https://gitlab.com:<personal_access_token>@gitlab.com/user/repo.git'
	cmd := exec.Command("git", "clone", url, dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to clone repos: %v", err)
	}
	fmt.Printf("%s\n", out)

}

// ExecRepo is ...
func ExecRepo(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"clone": cloneRepo,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		if cmdArgs["help"] == "true" {
			repoHelpList := map[string]func(){
				"clone": utils.PrintHelpRepoClone,
			}
			repoHelpList[arrCmd[0]]()
			return
		}
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}
