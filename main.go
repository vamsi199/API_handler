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
Vamsi,vamsi@gmail.com,app1,client1,clientsecret1,redirecturl1
ram,ram@gmail.com,app1,client1,clientsecret1,redirecturl1
Vamsi,vamsi@gmail.com,app2,client2,clientsecret2,redirecturl2
prasad,prasad@gmail.com,app2,client2,clientsecret2,redirecturl2
sai,sai@gmail.com,app3,client3,clientsecret3,redirecturl3



*/

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"net/http"
	//"os"
	"encoding/csv"
)

type user struct {
	username, email []string
}
type app struct {
	application, clientid, clientsecret, redirecturl []string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/userdata", userhandler).Methods("POST")
	http.ListenAndServe(":8080", router)

}

func userhandler(w http.ResponseWriter, r *http.Request) {
	/*b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("cannot read the body", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	fmt.Println(b)
	s := string(b)
	fmt.Println(s)

	_, err = os.Create("input.txt")
	if err != nil {
		fmt.Println("cannot create the file", err)
		return

	}
	f, err := os.OpenFile("input.txt", os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("cannot open the file", err)
		return

	}
	defer f.Close()
	_, err = f.WriteString(s)
	if err != nil {
		fmt.Println("cannot write the content to the file", err)
		return

	}*/
	reader := csv.NewReader(r.Body)
	reader.Comma = ','
	reader.Comment = '#'
	reader.FieldsPerRecord = 6
	record, err := reader.ReadAll()
	if err != nil {
		fmt.Println("cannot read csv file", err)
		return

	}
	fmt.Println(record)

	for i, usr := range record {
		var u user
		var a app
		//var us []user
		//var ap []app

		u.email[i] = usr[0]
		u.username[i] = usr[1]
		a.application[i] = usr[2]
		a.clientid[i] = usr[3]
		a.clientsecret[i] = usr[4]
		a.redirecturl[i] = usr[5]
		fmt.Println(u, a)

	}
	fmt.Println()
}
