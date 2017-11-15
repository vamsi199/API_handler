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

Output:

username,email,user status,applications created;application already exists
vamsi,vamsi@gmail.com,user already exists,app2(client2);app1(client1)
ram,ram@gmail.com,user created,app2(client4);app1(client2)
sai,sai@gmail.com,user created,app3(client5);



*/

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"net/http"
	//"os"
	"encoding/csv"
	//"os"
	"os"
)

var appstatus string
var userstatus string

type user struct {
	username, email string
}
type users struct {
	usernames, emails []string
}
type app struct {
	application, clientid, clientsecret, redirecturl string
}
type apps struct {
	applications, clientids, clientsecrets, redirecturls []string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/userdata", userhandler).Methods("POST")
	http.ListenAndServe(":8080", router)

}
func (u user) userexists(user string) (userstatus string) {
	if u.username == "vamsi" {
		userstatus = "user already exists"
		//fmt.Print(u.username, u.email, userstatus)
	} else {
		userstatus = "user not exits,create an account"
		//fmt.Print(u.username, u.email, userstatus)
	}
	return

}
func (a app) appexists(app string) (appstatus string) {

	if a.application == "app1" {
		appstatus = "application already exists"
		//fmt.Print(a.application, a.clientid, a.clientsecret, a.redirecturl, appstatus)
	} else {
		appstatus = "application not exits,create an application"
		//fmt.Print(a.application, a.clientid, a.clientsecret, a.redirecturl, appstatus)
	}
	return
}

func userhandler(w http.ResponseWriter, r *http.Request) {

	reader := csv.NewReader(r.Body)
	reader.Comma = ','
	reader.Comment = '#'
	reader.FieldsPerRecord = 6
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("cannot read csv file", err)
		return

	}
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

	for i, usr := range record {

		var u user
		var a app
		var us users
		var as apps
		var out string

		u.username = usr[0]
		u.email = usr[1]
		a.application = usr[2]
		a.clientid = usr[3]
		a.clientsecret = usr[4]
		a.redirecturl = usr[5]
		if i == 0 {
			out=fmt.Sprintln(u.username, ",", u.email, ",", "user status", ",", a.application, ",", "application status")
		} else {
			out=fmt.Sprintln(u.username, ",", u.email, ",", fmt.Sprint(u.userexists(u.username)), ",", a.application, ",", fmt.Sprint(a.appexists(a.application)))

		}
		us.usernames =append(us.usernames,u.username)
		us.emails=append(us.emails,u.email)
		as.applications=append(as.applications,a.application)
		as.clientids=append(as.clientids,a.clientid)
		as.clientsecrets=append(as.clientsecrets,a.clientsecret)
		as.redirecturls=append(as.redirecturls,a.redirecturl)
		fmt.Println(out,as,us)
		_, err = f.WriteString(out)
		if err != nil {
			fmt.Println("cannot write the content to the file", err)
			return
		}
	}
	fmt.Fprint(w, "file received")

}
