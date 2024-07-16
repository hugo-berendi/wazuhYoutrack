package parser

import (
	"errors"
	"fmt"
	"time"
	"wazuhIssues/internal/api"
	"wazuhIssues/internal/mails"
	"wazuhIssues/internal/regex"
)

func mailToIssue(mail mails.MailContent) (api.Issue, error) {
	return api.Issue{
		Summary:     mail.Rule,
		Description: mail.Date.String() + " " + mail.Vm,
	}, nil
}

func ProcessMails(mailList []string) ([]mails.MailContent, error) {
	var mailContents []mails.MailContent

	for _, mail := range mailList {
		content, err := processMail(mail)
		if err != nil {
			return nil, err
		}
		mailContents = append(mailContents, content)
	}
	return mailContents, nil
}

func processMail(mail string) (mails.MailContent, error) {
	vm, err := getVmFromMail(mail)
	if err != nil {
		return mails.MailContent{}, err
	}

	rule, err := getRuleFromMail(mail)
	if err != nil {
		return mails.MailContent{}, err
	}

	date, err := getDateFromMail(mail)
	if err != nil {
		return mails.MailContent{}, err
	}

	return mails.MailContent{Vm: vm, Rule: rule, Date: date}, nil
}

func getVmFromMail(mailContent string) (string, error) {
	pattern := `Received From: \((.*?)\)`
	match := regex.GetStringSubmatch(mailContent, pattern)

	if match[1] == "" {
		return "", errors.New("no vm found")
	}

	return match[1], nil
}

func getRuleFromMail(mailContent string) (string, error) {
	pattern := `(?m)^Rule: (\d+) fired \(level (\d+)\) -> "(.*?)"$`
	match := regex.GetString(mailContent, pattern)
	if match == "" {
		return "", errors.New("no vm found")
	}

	return match, nil
}

func getDateFromMail(mail string) (time.Time, error) {
	pattern := `(?m)^\d{4} \w{3} \d{2} \d{2}:\d{2}:\d{2}`
	dateStr := regex.GetString(mail, pattern)

	if dateStr == "" {
		return time.Time{}, fmt.Errorf("no date found in mail")
	}

	layout := "2006 Jan 02 15:04:05"
	parsedDate, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedDate, nil
}
