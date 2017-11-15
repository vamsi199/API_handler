/*Read the input file from the API body to a go-object. You can assume that the input is in body as binary or in the body as form-file.
Prepare distinct list of user data (username, email)
For each of these users, call a stub function to check if the user already exists in the system.
If not exists, call another stub function to create the user

Prepare the list of applications (application, clientid, clientsecret, redirecturl)
For each of these applications, call a stub function to check if the application already exists in the system.
If not exist, call another stub function to create the application

Prepare the final output in the above specified format, one row per user.
In the “user status” field, if the user was already in the system, mention “user already exists”, else “user created”
Under the “applications created” field, list all the application for this user which were not already there in the system, but were created as part of this request.
Under the “applications already exists” field, list all the application for this user which were already in the system
The format for application list in above two is like this “<application> (<clientid>)”
See output snapshot image for the exact layout.

input:
username,email,application, clientid, clientsecret, redirecturl
vamsi,vamsi@gmail.com,app1,client1,clientsecret1,redirecturl1
ram,ram@gmail.com,app1,client2,clientsecret2,redirecturl2
vamsi,vamsi@gmail.com,app2,client3,clientsecret3,redirecturl3
ram,ram@gmail.com,app2,client4,clientsecret4,redirecturl4
sai,sai@gmail.com,app3,client5,clientsecret5,redirecturl5
vamsi,vamsi@gmail.com,app4,client6,clientsecret6,redirecturl6

Output:

username,email,user status,applications created,application already exists
vamsi,vamsi@gmail.com,user already exists,app2(client2);app4(client6),app1(client1)
ram,ram@gmail.com,user created,app2(client4),app1(client2)
sai,sai@gmail.com,user created,app3(client5),



*/

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gocarina/gocsv"
	"net/http"
	"os"
	"strings"
	"encoding/csv"
)

type user struct {
	username string `csv:"username"`
	email    string `csv:"email"`
}
type app struct {
	application, clientid, clientsecret, redirecturl string
}
type userApp struct {
	user
	apps []app
}
type userAppOutput struct {
	user
	userStatus               string `csv:"user status"`
	applicationsCreated      string `csv:"applications created"`
	applicationsAlreadyExist string `csv:"application already exists"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/import", importUserApplicationHandler).Methods("POST")
	http.ListenAndServe(":8089", router)
}


//this does all the job
func importUserApplicationHandler(w http.ResponseWriter, r *http.Request) { //(c *gin.Context) {

	// read the input file

	reader := csv.NewReader(r.Body)
	reader.Comma = ','
	reader.Comment = '#'
	reader.FieldsPerRecord = 6
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("cannot read csv file", err)
		return
	}

	// generate empty output file
	_, err = os.Create("output.txt")
	if err != nil {
		fmt.Println("cannot create the file", err)
		return
	}
	f, err := os.OpenFile("output.txt", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("cannot open the file", err)
		return
	}
	defer f.Close()

	userApps := map[string]userApp{} // map of username to user data
	// process each record
	for _, r := range record {

		var a app
		var ua userApp

		ua.username = r[0]
		ua.email = r[1]
		a.application = r[2]
		a.clientid = r[3]
		a.clientsecret = r[4]
		a.redirecturl = r[5]

		ua.apps = append(userApps[ua.username].apps, a)
		userApps[ua.username] = ua
	}

	// TODO: process the records
	output := []userAppOutput{}
	for _, ua := range userApps { // for each of the users
		//process user data
		// check if user already exists
		out := userAppOutput{}
		out.username = ua.username
		if ua.user.exists() {
			out.userStatus = "user already exists"
		} else {
			out.userStatus = "user created"
		}
		// process application data
		applicationsCreated := []string{}
		applicationsAlreadyExist := []string{}
		for _, a := range ua.apps { // for each of the users' apps
			// check of the app exists, and process data accordingly
			if a.exists(ua.username) {
				applicationsAlreadyExist = append(applicationsAlreadyExist, a.application+"("+a.clientid+")")
			} else {
				err := a.create(ua.username)
				if err != nil {
					//TODO: return error to client using http.Error()
					return
				}
				applicationsCreated = append(applicationsCreated, a.application+"("+a.clientid+")")
			}
		}
		out.applicationsCreated = strings.Join(applicationsCreated, ";")
		out.applicationsAlreadyExist = strings.Join(applicationsAlreadyExist, ";")
		output = append(output, out)
	}

	// generate final output file and send output to client
	err = gocsv.Marshal(output, w)
	if err != nil{
		//TODO
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=devices.csv")
	w.Header().Set("Content-Type","text/csv")

	//
	// data := services.ExportSelectedDevicesAsCSV(token, deviceIDs)

	// mlog.Debug("%s response = %s", METHOD_NAME, string(data.Bytes()))

	// c.Header("Content-Disposition", "attachment; filename=devices.csv")
	// c.Data(http.StatusOK, "text/csv", data.Bytes()) // text/csv

}

func (u user) exists() bool {
	return u.username == "vamsi"
}
func (u user) create() error {
	//TODO: create user
	return nil
}

func (a app) exists(username string) bool {
	return a.application == "app1"
}
func (a app) create(username string) error { // NOTE: dont need to pass all the fields again, as this is a method, all the values will be available using the  receiver object, i.e. 'a'a in this case.
	//TODO: create application
	return nil
}

//// NOTES
// user duplicate will no longer exist as we are using maps. maps will overwrite if same key is passed again. and we are appending the apps, so we are good here
// app duplicates will not exist. we can go with this assumption
