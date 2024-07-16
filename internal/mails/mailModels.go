package mails

import "time"

type MailContent struct {
	Vm   string
	Rule string
	Date time.Time
}
