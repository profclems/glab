package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/logrusorgru/aurora"
)

//use struct instead of interface
//(json.Unmarshal causing big integer which pipline ID value to change)
type pipline struct {
	ID        int64  `json:"id"`
	Sha       string `json:"sha"`
	Ref       string `json:"ref"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	WebURL    string `json:"web_url"`
}

func displayMultiplePiplines(m []pipline) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	if len(m) > 0 {
		fmt.Printf("Showing piplines %d of %d on %s\n\n", len(m), len(m), GetEnv("GITLAB_REPO"))
		for i := 0; i < len(m); i++ {

			duration := TimeAgo(m[i].CreatedAt)

			if m[i].Status == "success" {
				_, _ = fmt.Fprintln(w, aurora.Green(fmt.Sprint("#", m[i].ID)), "\t", m[i].Ref, "\t", aurora.Magenta(duration))
			} else {
				_, _ = fmt.Fprintln(w, aurora.Red(fmt.Sprint("#", m[i].ID)), "\t", m[i].Ref, "\t", aurora.Magenta(duration))
			}
		}
	} else {
		fmt.Println("No Pipelines available on " + GetEnv("GITLAB_REPO"))
	}
}
func deletePipeline(cmdArgs map[string]string, arrFlags map[int]string) {
	issueID := strings.Trim(arrFlags[1], " ")
	if CommandArgExists(cmdArgs, issueID) {
		arrIds := strings.Split(strings.Trim(issueID, "[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Pipeline #" + i2)
			queryStrings := "/" + i2
			resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/pipelines"+queryStrings, "DELETE")
			if resp["responseCode"] == 204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(aurora.Green("Pipeline Deleted Successfully"))
			} else if resp["responseCode"] == 404 {
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
	var queryStrings = "status="
	if CommandArgExists(cmdArgs, "all") {
		queryStrings = ""
	} else if CommandArgExists(cmdArgs, "failed") {
		queryStrings += "failed&"
	} else {
		queryStrings += "success&"
	}
	if CommandArgExists(cmdArgs, "order_by") {
		queryStrings += "order_by=" + cmdArgs["order_by"] + "&"
	}
	if CommandArgExists(cmdArgs, "sort") {
		queryStrings += "sort=" + cmdArgs["sort"] + "&"
	}
	queryStrings = strings.Trim(queryStrings, "& ")
	if len(queryStrings) > 0 {
		queryStrings = "?" + queryStrings
	}

	resp := MakeRequest("{}", "projects/"+GetEnv("GITLAB_PROJECT_ID")+"/pipelines"+queryStrings, "GET")

	if resp["responseCode"] == 200 {
		bodyString := resp["responseMessage"]
		if _, ok := bodyString.(string); ok {
			// fmt.Println(bodyString.(string))
			var m []pipline
			err := json.Unmarshal([]byte(bodyString.(string)), &m)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println()

			displayMultiplePiplines(m)
			fmt.Println()

		}
	} else {
		fmt.Println(resp["responseCode"], resp["responseMessage"])
	}
}

// ExecPipeline is exported
func ExecPipeline(cmdArgs map[string]string, arrCmd map[int]string) {
	commandList := map[interface{}]func(map[string]string, map[int]string){
		"list":   listPipeline,
		"delete": deletePipeline,
	}
	if _, ok := commandList[arrCmd[0]]; ok {
		commandList[arrCmd[0]](cmdArgs, arrCmd)
	} else {
		fmt.Println(arrCmd[0]+":", "Invalid Command")
	}
}
