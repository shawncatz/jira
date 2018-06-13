package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/andygrunwald/go-jira"
)

var jiraClient *jira.Client

// SprintCompletedDate is for sorting based on completed date
type SprintCompletedDate []jira.Sprint

// Len fulfills the sort interface
func (a SprintCompletedDate) Len() int { return len(a) }

// Swap fulfills the sort interface
func (a SprintCompletedDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less fulfills the sort interface
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

// CreateAnswers stores the answers from the questions of the create command
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

	tp := jira.BasicAuthTransport{
		Username: cfg.Jira.User,
		Password: cfg.Jira.Pass,
	}

	jiraClient, err = jira.NewClient(tp.Client(), cfg.Jira.Base)
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
	if !debug {
		return
	}

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
	return cfg.Jira.Base + "/browse/" + id
}

func getBoards() (list []jira.Board, err error) {
	project := cfg.Jira.Project
	options := &jira.BoardListOptions{ProjectKeyOrID: project}
	br, response, err := jiraClient.Board.GetAllBoards(options)
	if err != nil {
		printErrResponse(response)
		return list, err
	}
	return br.Values, err
}

func getSprints(boardID int, all bool) (list []jira.Sprint, err error) {
	options := &jira.GetAllSprintsOptions{}
	if !all {
		options.State = "active,future"
	}
	return getSprintsWalk(boardID, all)
}

func getSprintsWalk(boardID int, all bool) (list []jira.Sprint, err error) {
	resp := []jira.Sprint{}

	options := &jira.GetAllSprintsOptions{SearchOptions: jira.SearchOptions{StartAt: 0}}
	if !all {
		options.State = "active,future"
	}

	// continue making the call and appending until we get the last response.
	for {
		list, _, err := jiraClient.Board.GetAllSprintsWithOptions(boardID, options)
		if err != nil {
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

func getJiraSprint(boardID, sprintID int) (*jira.Sprint, error) {
	sprints, err := getSprints(boardID, false)
	if err != nil {
		return nil, err
	}

	for _, s := range sprints {
		if s.ID == sprintID {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("sprint (%d) not found on board (%d)", sprintID, boardID)
}

func getIssuesFromSprint(sprintID int) ([]jira.Issue, error) {
	list, response, err := jiraClient.Sprint.GetIssuesForSprint(sprintID)
	if err != nil {
		printErrResponse(response)
		return nil, err
	}

	return list, nil
}

func getPointsField() (*jira.Field, error) {
	fields, response, err := jiraClient.Field.GetList()
	if err != nil {
		printErrResponse(response)
		return nil, err
	}

	for _, f := range fields {
		if f.Name == "Story Points" {
			return &f, nil
		}
	}

	return nil, nil
}
