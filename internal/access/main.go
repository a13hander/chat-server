package access

import (
	"context"

	accessV1 "github.com/a13hander/auth-service-api/pkg/access_v1"
)

type accessChecker struct {
	accessClient accessV1.AccessV1Client
}

func NewAccessChecker(accessClient accessV1.AccessV1Client) *accessChecker {
	return &accessChecker{accessClient: accessClient}
}

func (a *accessChecker) Check(ctx context.Context, endpoint string) (bool, error) {
	_, err := a.accessClient.Check(ctx, &accessV1.CheckRequest{
		Endpoint: endpoint,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
