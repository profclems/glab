package commands


func getRepoContributors()  {
	MakeRequest(`{}`,"projects/20131402/issues/1","GET")
}

func ExecRepo(cmdArgs map[string]string)  {
}