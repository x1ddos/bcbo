package bo

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/crhym3/bcbo/bc"
	"github.com/gorilla/mux"
)

// Api is the BC Backoffice API.
//
// Currently, it talks RESTfully in JSON over HTTP.
// To create an instance of the Api call NewApi() method.
//
// A typical usage would look like this:
//
// 		api, err = NewApi("/api/", bcClient, bcApiKey)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		http.Handle("/api/", api)
//
type Api struct {
	mux.Router
	bcClient bc.Client
	// TODO: remove the key once authentication is in place
	tempApiKey string
}

// NewApi creates a new instance of Backoffice API.
// Typically, the prefix is something like "/api/".
// apikey is a BC customer API key; it will be removed from the args
// once front-end authentication is in place.
func NewApi(prefix string, s bc.Client, apikey string) (*Api, error) {
	if s == nil {
		return nil, errors.New("BC service cannot be nil")
	}
	if apikey == "" {
		return nil, errors.New("Still need a customer API key, sorry")
	}
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	api := &Api{*r, s, apikey}
	api.HandleFunc("/{customer}/users", asJSON(api.ListUsers)).Methods("GET")
	api.HandleFunc("/{customer}/users/{id}/profile", asJSON(api.GetUserProfile)).Methods("GET")
	api.HandleFunc("/{customer}/users/{id}/activities", asJSON(api.ListActivities)).Methods("GET")
	api.HandleFunc("/{customer}/users/{id}/charts", asJSON(api.GetProfileCharts)).Methods("GET")
	return api, nil
}

// ListUsers fetches all users for a given customer
func (a *Api) ListUsers(r *http.Request) (interface{}, error) {
	customer := mux.Vars(r)["customer"]
	cred := &bc.Cred{customer, a.tempApiKey}
	users, err := a.bcClient.ListUsers(cred)
	if err != nil {
		return nil, err
	}
	return &ListUsersResponse{Items: users}, nil
}

// GetUserProfile returns a profile of a single user identified by an ID.
func (a *Api) GetUserProfile(r *http.Request) (interface{}, error) {
	customer := mux.Vars(r)["customer"]
	userId := mux.Vars(r)["id"]
	cred := &bc.Cred{customer, a.tempApiKey}
	return a.bcClient.GetUserProfile(cred, userId)
}

// ListActivites fetches a chunk of activities of a single user.
// Activities are ordered by date/time in descending order.
// To fetch older activities supply a "page" number > 0.
func (a *Api) ListActivities(r *http.Request) (interface{}, error) {
	customer := mux.Vars(r)["customer"]
	userId := mux.Vars(r)["id"]
	page, _ := strconv.Atoi(r.FormValue("page"))
	cred := &bc.Cred{customer, a.tempApiKey}
	activities, err := a.bcClient.ListActivities(cred, userId, page)
	if err != nil {
		return nil, err
	}
	return ListActivitiesResponse{Items: activities}, nil
}

// GetProfileCharts returns a list of simple chart data based on user activities.
func (a *Api) GetProfileCharts(r *http.Request) (interface{}, error) {
	customer := mux.Vars(r)["customer"]
	userId := mux.Vars(r)["id"]
	cred := &bc.Cred{customer, a.tempApiKey}
	activities, err := a.bcClient.ListActivities(cred, userId, 0)
	if err != nil {
		return nil, err
	}

	actGrouped := make(map[string]int)
	verbGrouped := make(map[string]int)
	for _, a := range activities {
		if _, ok := verbGrouped[a.Verb]; !ok {
			verbGrouped[a.Verb] = 0
		}
		verbGrouped[a.Verb] += 1

		day := a.Timestamp.Format("2006-01-02")
		if _, ok := actGrouped[day]; !ok {
			actGrouped[day] = 0
		}
		actGrouped[day] += 1
	}

	actRows := make([][]interface{}, 0, len(actGrouped))
	for day, counter := range actGrouped {
		actRows = append(actRows, []interface{}{day, counter})
	}

	verbRows := make([][]interface{}, 0, len(verbGrouped))
	for verb, counter := range verbGrouped {
		verbRows = append(verbRows, []interface{}{verb, counter})
	}

	actData := ChartData(actRows)
	verbData := ChartData(verbRows)
	sort.Sort(actData)
	sort.Sort(verbData)

	actChart := &Chart{
		Title: "Activities",
		Columns: []*ChartColumn{
			&ChartColumn{"Day", "date"},
			&ChartColumn{"Count", "number"},
		},
		Rows: actData,
	}

	verbChart := &Chart{
		Title: "Actions",
		Columns: []*ChartColumn{
			&ChartColumn{"Action", "string"},
			&ChartColumn{"Count", "number"},
		},
		Rows: verbData,
	}

	return ChartsResponse{Charts: []*Chart{actChart, verbChart}}, nil
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func (err *apiError) Error() string {
	return err.Message
}

// asJSON is an API method wrapper that knows how to respond with a JSON error.
func asJSON(f func(r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		resp, err := f(r)
		if err == nil {
			err = json.NewEncoder(w).Encode(resp)
		}
		writeJSONError(w, err)
	}
}

func writeJSONError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	log.Printf("ERROR: %v\n", err)
	var apiErr *apiError
	switch err.(type) {
	default:
		apiErr = &apiError{http.StatusInternalServerError, err.Error()}
	case *apiError:
		apiErr = err.(*apiError)
	}
	w.WriteHeader(apiErr.Code)
	json.NewEncoder(w).Encode(apiErr)
}
