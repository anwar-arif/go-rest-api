package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status   int     `json:"-"`
	Next     *string `json:"next,omitempty"`
	Previous *string `json:"previous,omitempty"`
	//Results  interface{} `json:"results"`
	//Count    int32       `json:"count" bson:"count"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Username    string   `json:"username"`
	Usertype    string   `json:"user_type,omitempty"`
	Groups      []string `json:"groups"`
	Contact     string   `json:"contact,omitempty" bson:"contact,omitempty"`
	FirstName   string   `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName    string   `json:"last_name,omitempty" bson:"last_name,omitempty"`
	IsStaff     bool     `json:"is_staff,omitempty" bson:"is_staff,omitempty"`
	IsActive    bool     `json:"is_active,omitempty" bson:"is_active,omitempty"`
	IsSuperuser bool     `json:"is_superuser,omitempty"  bson:"is_superuser,omitempty"`
	UserID      string   `json:"userid,omitempty"  bson:"userid,omitempty"`
	Token       string   `json:"token,omitempty"  bson:"token,omitempty"`
	Role        string   `json:"role,omitempty"  bson:"role,omitempty"`
}

type AuthUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    User   `json:"data"`
}

// ServeJSON serves json to http client
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

// ServeJSON a utility func which serves json to http client
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
