package report

import (
	"fmt"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/logrusorgru/aurora"
)

type SprintReport struct {
	*Report
	SprintID   int
	SprintName string
	Done       []*SprintReportIssue
	Todo       []*SprintReportIssue
	Open       []*SprintReportIssue

	issues []jira.Issue
	a      aurora.Aurora
}

type SprintReportIssue struct {
	Key     string
	Points  float64
	Status  string
	Type    string
	Closed  time.Time
	Summary string
}

func NewSprintReport(c *jira.Client, sprintID int, sprintName string) *SprintReport {
	return &SprintReport{
		SprintID:   sprintID,
		SprintName: sprintName,
		a:          aurora.NewAurora(true),
		Report:     &Report{client: c},
	}
}

func (r *SprintReport) Build() error {
	err := r.getIssues()
	if err != nil {
		return fmt.Errorf("error getting issues: %s", err)
	}

	field, err := r.getPointsField()
	if err != nil {
		return fmt.Errorf("error getting points field: %s", err)
	}

	for _, i := range r.issues {
		if i.Fields.Type.Name == "Sub-task" {
			continue
		}

		var points float64
		if i.Fields.Unknowns[field.Key] != nil {
			points = i.Fields.Unknowns[field.Key].(float64)
		}

		issue := &SprintReportIssue{
			Key:     i.Key,
			Points:  points,
			Status:  i.Fields.Status.Name,
			Type:    i.Fields.Type.Name,
			Closed:  time.Time(i.Fields.Resolutiondate),
			Summary: i.Fields.Summary,
		}

		switch i.Fields.Status.StatusCategory.Name {
		case "Done":
			r.Done = append(r.Done, issue)
		case "To Do":
			r.Todo = append(r.Todo, issue)
		case "In Progress":
			r.Open = append(r.Open, issue)
		}
	}

	return nil
}

func (r *SprintReport) Print() {
	r.printReport()
}

func (r *SprintReport) printReport() {
	fmt.Printf("%s %s %s\n", r.a.Bold("Sprint Report").Gray(), "for", r.a.Bold(r.SprintName).Cyan())

	fmt.Printf("\n%s\n", r.a.Bold("To Do").Gray())
	if len(r.Todo) > 0 {
		for _, i := range r.Todo {
			r.printReportIssue(i)
		}
	} else {
		fmt.Printf(" no issues\n")
	}
	fmt.Printf("\n%s\n", r.a.Bold("Open").Gray())
	if len(r.Open) > 0 {
		for _, i := range r.Open {
			r.printReportIssue(i)
		}
	} else {
		fmt.Printf("  no issues\n")
	}
	fmt.Printf("\n%s\n", r.a.Bold("Done").Gray())
	if len(r.Done) > 0 {
		for _, i := range r.Done {
			r.printReportIssue(i)
		}
	} else {
		fmt.Printf("  no issues\n")
	}
}

func (r *SprintReport) printReportIssue(issue *SprintReportIssue) {
	fmt.Printf("%10.10s %3.0f %-15.15s %-10.10s %-50.50s\n", r.a.Bold(issue.Key).Cyan(), issue.Points, issue.Status, issue.Type, issue.Summary)
}

func (r *SprintReport) getIssues() error {
	list, _, err := r.client.Sprint.GetIssuesForSprint(r.SprintID)
	if err != nil {
		return err
	}

	r.issues = list
	return nil
}

func (r *SprintReport) getPointsField() (*jira.Field, error) {
	fields, _, err := r.client.Field.GetList()
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		if f.Name == "Story Points" {
			return &f, nil
		}
	}

	return nil, nil
}
