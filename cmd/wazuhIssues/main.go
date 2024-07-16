package main

import (
	"log"
	"strings"
	"wazuhIssues/internal/api"
	"wazuhIssues/internal/mails"
	"wazuhIssues/internal/parser"
)

func main() {
	mailFolder := "./wazuhmails"

	files, err := mails.FindFiles(mailFolder, "wazuhmails/Wazuh notification - \\([^)]*\\) [A-Za-z]+ - Alert level [0-9]+-?[0-9]+?\\.eml")
	if err != nil {
		log.Fatal(err)
	}

	mailList, err := mails.ReadMails(files)
	if err != nil {
		log.Fatal(err)
	}

	mailContentList, err := parser.ProcessMails(mailList)

	for _, mailContent := range mailContentList {
		issueList, err := newIssueList()
		if err != nil {
			log.Fatal(err)
		}

		vmAddedToIssue := false

		for _, issue := range issueList {
			if mailContent.Rule == issue.Summary {
				if strings.Contains(issue.Description, mailContent.Vm) {
					vmAddedToIssue = true
					break
				}
				_, err := api.AddVmToIssue(mailContent.Vm, issue)
				if err != nil {
					log.Fatal(err)
				}
				vmAddedToIssue = true
				break
			}
		}

		if vmAddedToIssue {
			continue
		}

		_, err = newIssueWithMailContent(mailContent)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func newIssueWithMailContent(mailContent mails.MailContent) (api.Issue, error) {
	newIssue := api.Issue{
		Summary:     mailContent.Rule,
		Description: mailContent.Vm,
	}
	return api.CreateIssue(newIssue, "0-2")
}

func newIssueList() ([]api.Issue, error) {
	issueList, err := api.GetIssueList()
	if err != nil {
		return nil, err
	}

	var fullIssueList []api.Issue
	for _, issue := range issueList {
		fullIssue, err := api.GetIssueDetails(issue.ID, []string{})
		if err != nil {
			return nil, err
		}
		fullIssueList = append(fullIssueList, fullIssue)
	}
	return fullIssueList, nil
}
