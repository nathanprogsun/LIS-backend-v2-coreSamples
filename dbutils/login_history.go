package dbutils

import (
	"context"
	"coresamples/ent"
)

func CreateLoginHistory(username string, ip string,
	success bool, failureReason string,
	loginPortal string, token string,
	client *ent.Client, ctx context.Context) error {
	creator := client.LoginHistory.Create().
		SetUsername(username).
		SetLoginIP(ip).
		SetLoginSuccessfully(success)
	if failureReason != "" {
		creator.SetFailureReason(failureReason)
	}
	if loginPortal != "" {
		creator.SetLoginPortal(loginPortal)
	}
	if token != "" {
		creator.SetToken(token)
	}
	return creator.Exec(ctx)
}
