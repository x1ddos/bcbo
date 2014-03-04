package bo

import "github.com/crhym3/bcbo/bc"

type ListUsersResponse struct {
	Items []*bc.User `json:"items"`
	More  bool       `json:"more"`
}

type ListActivitiesResponse struct {
	Items []*bc.Activity `json:"items"`
	More  bool           `json:"more"`
}

type ChartsResponse struct {
	Charts []*Chart `json:"charts"`
}

type Chart struct {
	Title   string         `json:"title"`
	Columns []*ChartColumn `json:"columns"`
	Rows    ChartData      `json:"rows"`
}

type ChartColumn struct {
	Name string `json:"name"`
	Typ  string `json:"typ"`
}

type ChartData [][]interface{}

func (d ChartData) Len() int           { return len(d) }
func (d ChartData) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d ChartData) Less(i, j int) bool { return d[i][0].(string) < d[j][0].(string) }
