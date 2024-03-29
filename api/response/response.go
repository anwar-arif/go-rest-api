package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status   int     `json:"-"`
	Next     *string `json:"next,omitempty"`     // for pagination, link for next page
	Previous *string `json:"previous,omitempty"` // for pagination, link for previous page
	//Results  interface{} `json:"results"`
	//Count    int32       `json:"count" bson:"count"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Serve serves json to http client
func (r *Response) ServeJSON(w http.ResponseWriter) error {
	resp := &Response{
		Next:     r.Next,
		Previous: r.Previous,
		Data:     r.Data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

// Serve a utility func which serves json to http client
func ServeJSON(w http.ResponseWriter, status int, previous, next *string, message string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &Response{
		Status:   status,
		Next:     next,
		Previous: previous,
		//Results:  data,
		Data:    data,
		Message: message,
		//Count:   count,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}

func Serve(w http.ResponseWriter, status int, message string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &Response{
		Status:   status,
		Next:     nil,
		Previous: nil,
		Message:  message,
		Data:     data,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}
	return nil
}

func ServeJSONData(w http.ResponseWriter, data interface{}, status int) {
	res := responseData{
		status: status,
		Data:   data,
	}
	res.serveJSON(w)
}

type responseData struct {
	status int
	Data   interface{} `json:"data,omitempty"`
}

func (res *responseData) serveJSON(w http.ResponseWriter) {
	if res.status == 0 {
		res.status = http.StatusOK
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

// http api response messages
const (
	CannotProcessToken   = "can't process the token"
	CannotProcessRequest = "can't process the request"
	Successful           = "successful"
	UserNotFound         = "user not found"
	InvalidCredential    = "invalid credential"
	DeletedSuccessfully  = "deleted successfully"
	CreatedSuccessfully  = "created successfully"
)
