package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	baseUrl = "https://youtrack.dev.linova.de/api"
)

func DeleteIssue(issueId string) {
	requestUrl := fmt.Sprintf("/issues/%s", issueId)
	_, err := DeleteRequest(requestUrl)
	if err != nil {
		log.Fatal(err)
	}
}

func AddVmToIssue(vm string, issue Issue) (Issue, error) {
	issueUpdate := Issue{
		Summary:     strings.ReplaceAll(issue.Summary, `"`, `\"`),
		Description: issue.Description + " " + vm,
	}
	updatedIssue, err := UpdateIssue(issue.ID, []string{"description"}, []string{}, issueUpdate)
	if err != nil {
		return Issue{}, err
	}
	return updatedIssue, nil
}

func CreateIssue(issueToCreate Issue, projectId string) (Issue, error) {
	baseRequestUrl := fmt.Sprintf("/issues/?fields=id,summary,description")

	reqBody := []byte(fmt.Sprintf(`{
		"summary": "%s",
		"description": "%s",
		"project": {
			"id": "%s"
		}
	}`, strings.ReplaceAll(issueToCreate.Summary, `"`, `\"`), issueToCreate.Description, projectId))

	body, err := PostRequest(baseRequestUrl, bytes.NewReader(reqBody))
	if err != nil {
		return Issue{}, fmt.Errorf("failed to make POST request: %w", err)
	}

	return createIssueFromBody(body)
}

func UpdateIssue(issueId string, fieldsToUpdate []string, customFieldsToUpdate []string, updatedIssue Issue) (Issue, error) {
	baseRequestUrl := fmt.Sprintf("/issues/%s?fields=%s", issueId, strings.Join(fieldsToUpdate, ","))
	var requestUrl string
	if len(customFieldsToUpdate) == 0 {
		requestUrl = fmt.Sprintf(baseRequestUrl)
	} else {
		customFieldsToUpdateStr := strings.Join(customFieldsToUpdate, ",")
		requestUrl = baseRequestUrl + fmt.Sprintf(",customFields(%s)", customFieldsToUpdateStr)
	}

	reqBody := []byte(fmt.Sprintf(`{
		"summary": "%s",
		"description": "%s"
	}`, updatedIssue.Summary, updatedIssue.Description))

	body, err := PostRequest(requestUrl, bytes.NewReader(reqBody))
	if err != nil {
		return Issue{}, fmt.Errorf("failed to make POST request: %w", err)
	}

	return createIssueFromBody(body)
}

func GetIssueDetails(issueId string, fields []string) (Issue, error) {
	baseRequestUrl := fmt.Sprintf("/issues/%s?fields=id,summary,description", issueId)
	var requestUrl string
	if len(fields) == 0 {
		requestUrl = fmt.Sprintf(baseRequestUrl)
	} else {
		fieldsStr := strings.Join(fields, ",")
		requestUrl = baseRequestUrl + fmt.Sprintf(",customFields(%s)", fieldsStr)
	}

	body, err := GetRequest(requestUrl)
	if err != nil {
		return Issue{}, err
	}

	return createIssueFromBody(body)
}

func GetIssueList() ([]Issue, error) {
	body, err := GetRequest("/issues")
	if err != nil {
		return nil, err
	}

	var issues []Issue
	err = json.Unmarshal(body, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

func createIssueFromBody(body []byte) (Issue, error) {
	var issue Issue
	err := json.Unmarshal(body, &issue)
	if err != nil {
		return Issue{}, err
	}

	var customFields []map[string]interface{}
	if err := json.Unmarshal(body, &customFields); err == nil {
		issue.CustomFields = customFields
	}

	return issue, nil
}

func DeleteRequest(requestUrl string) ([]byte, error) {
	return Request("Delete", requestUrl, nil)
}

func PostRequest(requestUrl string, reqBody io.Reader) ([]byte, error) {
	return Request("POST", requestUrl, reqBody)
}

func GetRequest(requestUrl string) ([]byte, error) {
	return Request("GET", requestUrl, nil)
}

func Request(method string, requestUrl string, reqBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, baseUrl+requestUrl, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+getApiToken())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("received non-OK HTTP status: %s, response: %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func getApiToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return os.Getenv("YOUTRACK_TOKEN")
}
