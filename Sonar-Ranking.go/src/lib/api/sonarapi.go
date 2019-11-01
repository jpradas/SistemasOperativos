package api

import (
	b64 "encoding/base64"
	"encoding/json"
	"io"
	"lib/structs"
	"log"
	"net/http"
	"strings"
	"time"
)

type ApiClient struct {
	URLbase     string
	Credentials string
	Client      *http.Client
}

func New(URLbase string, Credentials string) *ApiClient {
	if !pingServer(URLbase) {
		log.Fatalln("Couldn't reach server. Did you type URL correctly?")
		log.Fatalln(URLbase)
	}

	c := &ApiClient{
		URLbase,
		base64Encoding(Credentials),
		&http.Client{Timeout: time.Duration(10 * time.Second)},
	}
	return c
}

func (c *ApiClient) RetrieveProjects() *structs.Projects {
	add := "/api/projects/search"
	body := makeGetRequest(add, c)

	defer body.Close()

	var projects structs.Projects
	if err := json.NewDecoder(body).Decode(&projects); err != nil {
		log.Fatalln(err)
	}

	return &projects
}

func (c *ApiClient) RetrieveIssues() *structs.Issues {
	add := "/api/issues/search?s=FILE_LINE&resolved=false&ps=500&facets=severities,types,rules&additionalFields=_all"
	body := makeGetRequest(add, c)

	defer body.Close()

	var issues structs.Issues
	if err := json.NewDecoder(body).Decode(&issues); err != nil {
		log.Fatalln(err)
	}

	return &issues
}

func (c *ApiClient) RetrieveIssuesByProjects(projects string) *structs.Issues {
	add := strings.Join([]string{"/api/issues/search?s=FILE_LINE&resolved=false&ps=500&facets=severities,types,rules&projects=", projects, "&additionalFields=_all"}, "")
	body := makeGetRequest(add, c)

	defer body.Close()

	var issues structs.Issues
	if err := json.NewDecoder(body).Decode(&issues); err != nil {
		log.Fatalln(err)
	}

	return &issues
}

func (c *ApiClient) RetrieveSecurityReport(project string) *structs.SecurityReport {
	add := strings.Join([]string{"/api/security_reports/show?project=", project, "&standard=owaspTop10&includeDistribution=false"}, "")
	body := makeGetRequest(add, c)

	defer body.Close()

	var report structs.SecurityReport
	if err := json.NewDecoder(body).Decode(&report); err != nil {
		log.Fatalln(err)
	}

	return &report
}

func pingServer(URLbase string) bool {
	request, err := http.Get(URLbase)
	if err != nil || request.StatusCode != 200 {
		return false
	}
	return true
}

func base64Encoding(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func makeGetRequest(add string, c *ApiClient) io.ReadCloser {
	request, err := http.NewRequest("GET", c.URLbase+add, nil)
	request.Header.Set("Authorization", "Basic "+c.Credentials)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := c.Client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	return resp.Body
}
