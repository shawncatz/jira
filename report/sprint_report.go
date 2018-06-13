package report

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/logrusorgru/aurora"
	"github.com/martinlindhe/imgcat/lib"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type SprintReport struct {
	*Report
	ImageSupported  bool
	TotalPoints     float64
	CompletedPoints float64
	Done            []*SprintReportIssue
	Todo            []*SprintReportIssue
	Open            []*SprintReportIssue

	issues []jira.Issue
	sprint *jira.Sprint
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

// SprintCompletedDate is for sorting based on completed date
type SprintReportIssueSorter []*SprintReportIssue

// Len fulfills the sort interface
func (a SprintReportIssueSorter) Len() int { return len(a) }

// Swap fulfills the sort interface
func (a SprintReportIssueSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less fulfills the sort interface
func (a SprintReportIssueSorter) Less(i, j int) bool { return a[i].Closed.Before(a[j].Closed) }

func NewSprintReport(c *jira.Client, sprint *jira.Sprint) *SprintReport {
	return &SprintReport{
		ImageSupported: true,
		a:              aurora.NewAurora(true),
		sprint:         sprint,
		Report:         &Report{client: c},
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

		r.TotalPoints += points

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
			r.CompletedPoints += points
			r.Done = append(r.Done, issue)
		case "To Do":
			r.Todo = append(r.Todo, issue)
		case "In Progress":
			r.Open = append(r.Open, issue)
		}
	}

	return nil
}

func (r *SprintReport) imageDataSeries() chart.TimeSeries {
	current := r.TotalPoints

	ts := chart.TimeSeries{
		XValues: []time.Time{},
		//	time.Now().AddDate(0, 0, -10),
		//	time.Now().AddDate(0, 0, -9),
		//	time.Now().AddDate(0, 0, -8),
		//	time.Now().AddDate(0, 0, -7),
		//	time.Now().AddDate(0, 0, -6),
		//	time.Now().AddDate(0, 0, -5),
		//	time.Now().AddDate(0, 0, -4),
		//	time.Now().AddDate(0, 0, -3),
		//	time.Now().AddDate(0, 0, -2),
		//	time.Now().AddDate(0, 0, -1),
		//	time.Now(),
		//},
		//YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0},
		YValues: []float64{},
	}

	sort.Sort(SprintReportIssueSorter(r.Done))

	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return ts
	}

	st := r.sprint.StartDate.In(la)
	ts.XValues = append(ts.XValues, time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location()))
	ts.YValues = append(ts.YValues, current)

	for _, i := range r.Done {
		if i.Points > 0 {
			current -= i.Points
			ts.XValues = append(ts.XValues, i.Closed)
			ts.YValues = append(ts.YValues, current)
		}
	}

	ts.XValues = append(ts.XValues, st.AddDate(0, 0, 11))
	ts.YValues = append(ts.YValues, current)

	return ts
}

func (r *SprintReport) imageGuideSeries() chart.TimeSeries {
	s := r.TotalPoints / 10
	c := r.TotalPoints
	ts := chart.TimeSeries{
		XValues: []time.Time{},
		YValues: []float64{},
	}

	la, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return ts
	}

	st := r.sprint.StartDate.In(la)
	t := time.Date(st.Year(), st.Month(), st.Day(), 0, 0, 0, 0, st.Location())
	ts.XValues = append(ts.XValues, t)
	ts.YValues = append(ts.YValues, c)

	fmt.Printf("start: %s, end: %s\n", r.sprint.StartDate, r.sprint.EndDate)
	//t := st
	for i := 0; i < 12; i++ {
		//fmt.Printf("%s %2.0f\n", t.Weekday().String(), c)
		ts.XValues = append(ts.XValues, t)
		ts.YValues = append(ts.YValues, c)

		if t.Weekday() != time.Sunday && t.Weekday() != time.Saturday {
			c -= s
		}
		t = t.AddDate(0, 0, 1)
	}

	ts.XValues = append(ts.XValues, t)
	ts.YValues = append(ts.YValues, c)

	return ts
}

func (r *SprintReport) Image() (string, error) {
	tmp, err := ioutil.TempFile("/tmp", "jira-sprint-report")
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	data := r.imageDataSeries()
	guide := r.imageGuideSeries()

	graph := &chart.Chart{
		Background: chart.Style{
			FillColor: drawing.ColorBlack,
		},
		Canvas: chart.Style{
			FillColor: drawing.ColorBlack,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show:      true,
				FontColor: drawing.ColorWhite,
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show:      true,
				FontColor: drawing.ColorWhite,
			},
		},
		Series: []chart.Series{
			data,
			guide,
		},
	}

	err = graph.Render(chart.PNG, tmp)
	if err != nil {
		return "", err
	}

	return tmp.Name(), nil
}

func (r *SprintReport) PrintImage() error {
	tmpName, err := r.Image()
	if err != nil {
		return err
	}

	tmp, err := os.Open(tmpName)
	if err != nil {
		return err
	}
	defer tmp.Close()

	imgcat.Cat(tmp, os.Stdout)

	return nil
}

func (r *SprintReport) Print() {
	r.printReport()
}

func (r *SprintReport) printReport() {
	fmt.Printf("%s for %s (%s %3.0f / %3.0f)\n",
		r.a.Bold("Sprint Report").Gray(),
		r.a.Bold(r.sprint.Name).Cyan(),
		r.a.Bold("Points:").Gray(),
		r.a.Bold(r.CompletedPoints).Cyan(),
		r.a.Bold(r.TotalPoints).Cyan())

	if r.ImageSupported {
		fmt.Println()
		r.PrintImage()
	}

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
	list, _, err := r.client.Sprint.GetIssuesForSprint(r.sprint.ID)
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
