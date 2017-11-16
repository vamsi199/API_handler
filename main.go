package main

import (
	"encoding/json"
	"github.com/gocarina/gocsv"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type user struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	responseCode int
}
type app struct {
	Name         string `json:"name"`
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"clientsecret"`
	RedirectURL  string `json:"redirecturl"`
	responseCode int
}
type userApp struct {
	user
	Apps []app
}
type userAppOutput struct {
	Username                 string `csv:"username"`
	Email                    string `csv:"email"`
	UserStatus               string `csv:"user status"`
	ApplicationsCreated      string `csv:"applications created"`
	ApplicationsAlreadyExist string `csv:"application already exists"`
}

const (
	StatusCreated       = 201 //TODO: is it 201, or some other code?
	StatusAlreadyExists = 409
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/import", importUserApplicationHandler).Methods("POST")
	http.ListenAndServe(":8089", router)
}

//this does all the job
func importUserApplicationHandler(w http.ResponseWriter, r *http.Request) { //(c *gin.Context) {

	// read the input json
	uas := []userApp{}
	err := json.NewDecoder(r.Body).Decode(&uas)
	if err != nil {
		//TODO
		return
	}
	defer r.Body.Close()

	// process the records
	output := []userAppOutput{}
	for _, ua := range uas { // for each of the users
		//process user data
		// check if user already exists
		out := userAppOutput{}
		out.Username = ua.Username
		out.Email = ua.Email
		resp, err := ua.user.create()
		if err != nil {
			//TODO
			return
		}
		if resp == StatusAlreadyExists {
			out.UserStatus = "user already exists"
		} else if resp == StatusCreated {
			out.UserStatus = "user created"
		}
		// process application data
		applicationsCreated := []string{}
		applicationsAlreadyExist := []string{}
		for _, a := range ua.Apps { // for each of the users' apps
			// check of the app exists, and process data accordingly
			resp, err = a.create(ua.Username)

			if err != nil {
				//TODO
				return
			}
			if resp == StatusAlreadyExists {
				applicationsAlreadyExist = append(applicationsAlreadyExist, a.Name+"("+a.ClientID+")")
			} else if resp == StatusCreated {
				applicationsCreated = append(applicationsCreated, a.Name+"("+a.ClientID+")")
			}
		}
		out.ApplicationsCreated = strings.Join(applicationsCreated, ";")
		out.ApplicationsAlreadyExist = strings.Join(applicationsAlreadyExist, ";")
		output = append(output, out)
	}

	// generate final output file and send output to client
	err = gocsv.Marshal(output, w)
	if err != nil {
		//TODO
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=devices.csv")
	w.Header().Set("Content-Type", "text/csv")

	//
	// data := services.ExportSelectedDevicesAsCSV(token, deviceIDs)

	// mlog.Debug("%s response = %s", METHOD_NAME, string(data.Bytes()))

	// c.Header("Content-Disposition", "attachment; filename=devices.csv")
	// c.Data(http.StatusOK, "text/csv", data.Bytes()) // text/csv

}

func (u user) create() (int, error) {
	//TODO: create user. call the api with required data and send the status code and error back
	if u.Username == "usr1" {
		return StatusAlreadyExists, nil
	} else {
		return StatusCreated, nil
	}
}

func (a app) create(username string) (int, error) {
	//TODO: create application. call the api with required data and send the status code and error back
	if a.Name == "app1" || a.Name == "app5" {
		return StatusAlreadyExists, nil
	}
	return StatusCreated, nil
}
