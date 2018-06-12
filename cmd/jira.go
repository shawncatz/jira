package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"

	"github.com/shawncatz/go-jira"
)

var jiraClient *jira.Client

type SprintCompletedDate []jira.Sprint

func (a SprintCompletedDate) Len() int           { return len(a) }
func (a SprintCompletedDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SprintCompletedDate) Less(i, j int) bool { return a[i].ID < a[j].ID }

// order by EndDate
//func (a SprintCompletedDate) Less(i, j int) bool {
//	it := a[i].EndDate
//	jt := a[j].EndDate
//	if it == nil {
//		return false
//	}
//	if jt == nil {
//		return true
//	}
//	return it.Before(*jt)
//}

type CreateAnswers struct {
	Project     string
	Title       string
	Description string
	Type        string
	Sprint      string
}

// initClient sets up the JIRA client
func initClient() {
	var err error
	base := viper.GetString("jira.base")
	user := viper.GetString("jira.user")
	pass := viper.GetString("jira.pass")

	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}

	jiraClient, err = jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}
}

func jiraCreate(answers *CreateAnswers) (*jira.Issue, error) {
	f := &jira.IssueFields{
		Project:     jira.Project{Key: answers.Project},
		Type:        jira.IssueType{Name: answers.Type},
		Summary:     answers.Title,
		Description: answers.Description,
		Labels:      []string{"from-cli"},
	}

	if answers.Sprint != "Backlog" {
		f.Sprint = &jira.Sprint{Name: answers.Sprint}
	}

	i := jira.Issue{
		Fields: f,
	}

	if debug {
		fmt.Printf("%#v\n", i)
		fmt.Printf("%#v\n", i.Fields)
	}

	issue, response, err := jiraClient.Issue.Create(&i)
	if err != nil {
		//printErr(err.Error())
		b, _ := ioutil.ReadAll(response.Response.Body)
		//fmt.Printf("response:\n%s\n", string(b))
		return nil, fmt.Errorf("error creating issue: %s\nResponse: %s", err.Error(), string(b))
	}

	return issue, nil
}

func printErrResponse(response *jira.Response) {
	r := response.Response
	printErr("Jira error: %d : %s\n", r.StatusCode, r.Status)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if len(b) > 0 {
		printErr("%s\n", b)
	}
}

func issueURL(id string) string {
	return viper.GetString("jira.base") + "/browse/" + id
}

func getBoards() (list []jira.Board, err error) {
	project := viper.GetString("jira.project")
	options := &jira.BoardListOptions{ProjectKeyOrID: project}
	br, response, err := jiraClient.Board.GetAllBoards(options)
	if err != nil {
		printErr("Error: %s", err)
		if debug {
			printErrResponse(response)
		}
		return list, err
	}
	return br.Values, err
}

func getSprints(id int, all bool) (list []jira.Sprint, err error) {
	options := &jira.GetAllSprintsOptions{}
	if !all {
		options.State = "active,future"
	}
	return getSprintsWalk(id, all)
}

func getSprintsWalk(boardID int, all bool) (list []jira.Sprint, err error) {
	resp := []jira.Sprint{}

	options := &jira.GetAllSprintsOptions{SearchOptions: jira.SearchOptions{StartAt: 0}}
	if !all {
		options.State = "active,future"
	}

	// continue making the call and appending until we get the last response.
	for {
		list, response, err := jiraClient.Board.GetAllSprintsWithOptions(boardID, options)
		if err != nil {
			if debug {
				printErrResponse(response)
			}
			return nil, err
		}

		resp = append(resp, list.Values...)

		if list.IsLast {
			break
		}

		options.StartAt = list.StartAt + len(list.Values)
	}

	return resp, nil
}
