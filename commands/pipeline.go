package commands

import (
	"fmt"
	"glab/cmd/glab/utils"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/logrusorgru/aurora"
	"github.com/xanzy/go-gitlab"
)

func displayMultiplePipelines(m []*gitlab.PipelineInfo) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing pipelines %d of %d on %s\n\n", len(m), len(m), getRepo())
		for _, pipeline := range m {
			duration := TimeAgo(*pipeline.CreatedAt)

			if pipeline.Status == "success" {
				_, _ = fmt.Fprintln(w, aurora.Green(fmt.Sprint("#", pipeline.ID)), "\t", pipeline.Ref, "\t", aurora.Magenta(duration))
			} else {
				_, _ = fmt.Fprintln(w, aurora.Red(fmt.Sprint("#", pipeline.ID)), "\t", pipeline.Ref, "\t", aurora.Magenta(duration))
			}
		}
	} else {
		fmt.Println("No Pipelines available on " + getRepo())
	}
}

func deletePipeline(cmdArgs map[string]string, arrFlags map[int]string) {
	pipelineID := strings.Trim(arrFlags[1], " ")
	git, repo := InitGitlabClient()

	if CommandArgExists(cmdArgs, pipelineID) {
		arrIds := strings.Split(strings.Trim(pipelineID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Pipeline #" + i2)
			pipeline, _ := git.Pipelines.DeletePipeline(repo, stringToInt(i2))
			if pipeline.StatusCode == 204 {
				fmt.Println(aurora.Green("Pipeline Deleted Successfully"))
			} else if pipeline.StatusCode == 404 {
				fmt.Println(aurora.Red("Pipeline does not exist"))
			} else {
				fmt.Println(aurora.Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(aurora.Red("Invalid command"))
		fmt.Println("Usage: glab pipeline delete <pipeline-id>")
	}
}

func listPipeline(cmdArgs map[string]string, _ map[int]string) {
	git, repo := InitGitlabClient()
	l := &gitlab.ListProjectPipelinesOptions{}
	if CommandArgExists(cmdArgs, "failed") {
		l.Status = gitlab.BuildState("failed")
	} else {
		l.Status = gitlab.BuildState("success")
	}
	if CommandArgExists(cmdArgs, "order_by") {
		l.OrderBy = gitlab.String(cmdArgs["order_by"])
	}
	if CommandArgExists(cmdArgs, "sort") {
		l.Sort = gitlab.String(cmdArgs["sort"])
	}
	mergeRequests, _, err := git.Pipelines.ListProjectPipelines(repo, l)
	if err != nil {
		log.Fatal(err)
	}
	displayMultiplePipelines(mergeRequests)
}

// ExecPipeline is exported
func ExecPipeline(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"list":   listPipeline,
		"delete": deletePipeline,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		if cmdArgs["help"] == "true" {
			pipelineHelpList := map[string]func(){
				"list":   utils.PrintHelpPipelineList,
				"delete": utils.PrintHelpPipelineDelete,
			}
			pipelineHelpList[arrCmd[0]]()
			return
		}
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}
