package commands


func getRepoContributors()  {
	MakeRequest(`{}`,"projects/20131402/issues/1","GET")
}

func NewBranch()  {

}

func ExecRepo(cmdArgs map[string]string)  {
}