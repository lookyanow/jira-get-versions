package main

import (
	"encoding/base64"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const jiraURL = "https://jit.ozon.ru/"
const apiUrl = jiraURL + "rest/api/latest/"

func jiraTokenToUserPass(token string) (user string, pass string) {
	sDec, _ := base64.StdEncoding.DecodeString(token)
	s := strings.SplitN(string(sDec), ":", 2)
	return s[0], s[1]
}

type jiraVersions struct {
	self string
	nextPage string
	maxResults int
	startAt int
	total int
	isLast bool
	Values []jira.Version
}

func getVersionsFromProject(project interface{}, orderby string) ([]string, error){
	url :=  "rest/api/latest/project/" + project.(string) + "/version?orderBy=" + orderby
	user, pass := jiraTokenToUserPass(os.Getenv("JIRA_TOKEN"))
	tp := jira.BasicAuthTransport{
		Username: user,
		Password: pass,
	}
	client, err := jira.NewClient(tp.Client(), jiraURL)
	if err != nil {
		return nil, err
	}

	req, err := client.NewRequest("GET", url, nil)
	if err != nil{
		return nil, err
	}

	ver := new(jiraVersions)
	_, err1 := client.Do(req, ver)
	if err1 != nil{
		return nil, err1
	}
	if len(ver.Values) == 0{
		return nil, errors.New("version count is 0")
	}

	var versions []string
	for _, v := range(ver.Values){
		if !v.Released{
			versions = append(versions, v.Name)
		}
	}
	return versions, nil
}
func main(){

	versions,err := getVersionsFromProject("ANDROID","sequence")
	if err != nil{
		panic(err)
	}
	for _, v := range(versions){
		fmt.Println(v)
	}
}

