package commands

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"strings"
)

func DeleteMergeRequest(cmdArgs map[string]string, arrFlags map[int]string)  {
	mergeId := strings.Trim(arrFlags[1]," ")
	if CommandArgExists(cmdArgs, mergeId) {
		arrIds := strings.Split(strings.Trim(mergeId,"[] "), ",")
		for _, i2 := range arrIds {
			fmt.Println("Deleting Merge Request #"+i2)
			queryStrings := "/"+i2
			resp := MakeRequest("{}","projects/"+GetEnv("GITLAB_PROJECT_ID")+"/merge_requests"+queryStrings,"DELETE")
			if resp["responseCode"]==204 {
				bodyString := resp["responseMessage"]
				fmt.Println(bodyString)
				fmt.Println(Green("Merge Request Deleted Successfully"))
			} else if resp["responseCode"]==404 {
				fmt.Println(Red("Merge Request does not exist"))
			} else {
				fmt.Println(Red("Could not complete request."))
			}
			fmt.Println()
		}
	} else {
		fmt.Println(Red("Invalid command"))
		fmt.Println("Usage: glab merge delete <merge-id>")
	}
}

func ExecMergeRequest(cmdArgs map[string]string, arrCmd map[int]string)  {
	commandList := map[interface{}]func(map[string]string,map[int]string) {
		"delete" : DeleteMergeRequest,
	}
	commandList[arrCmd[0]](cmdArgs, arrCmd)
}
