package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/showflight", showflight)
	fmt.Println("listening...")
	err := http.ListenAndServe(GetPort(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// Get the Port from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "4747"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, rootForm)
}

const rootForm = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Get Flight Plan</title>
<link rel="stylesheet" href="/stylesheets/goview.css">
</head>
<body>
<h1><img style="margin-left: 120px;" src="images/gsv.png" alt="Flight Plan" /></h1>
<p>Please enter your destination:</p>
<form style="margin-left: 120px;" action="/showflight" method="post" accept-charset="utf-8">
<input type="text" name="str" value="Destination" id="str" />
<input type="submit" value="Get flight plan." />
</form>
</body>
</html>
`

var upperTemplate = template.Must(template.New("showflight").Parse(upperTemplateHTML))

func showflight(w http.ResponseWriter, r *http.Request) {
	dest := r.FormValue("str")
	safeDest := url.QueryEscape(dest)
	fullUrl := fmt.Sprintf(
		"https://www.google.com/flights/#search;f=%s;t=JFK,EWR,LGA;d=2015-04-08;r=2015-04-12",
		safeDest)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
	}
	defer resp.Body.Close()
	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
	}
	res := make(map[string][]map[string]map[string]map[string]interface{}, 0)
	json.Unmarshal(body, &res)
	// %.13f is used to convert float64 to a string
	queryUrl :=
		fmt.Sprintf(
			"https://www.google.com/flights/#search;f=%s;t=JFK,EWR,LGA;d=2015-04-08;r=2015-04-12")
	tempErr := upperTemplate.Execute(w, queryUrl)
	if tempErr != nil {
		http.Error(w, tempErr.Error(), http.StatusInternalServerError)
	}
}

const upperTemplateHTML = `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>Display Flight</title>
<link rel="stylesheet" href="/stylesheets/goview.css">
</head>
<body>
<h1>Your Flight Plan.</h1>
<a href="https://www.google.com/flights/#search;f=%s;t=JFK,EWR,LGA;d=2015-04-08;r=2015-04-12">Your flight.<a>
</body>
</html>
`
