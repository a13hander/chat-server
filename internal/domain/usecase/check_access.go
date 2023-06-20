package usecase

import (
	"context"
)

type AccessChecker interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type CheckAccessUseCase interface {
	Run(ctx context.Context, endpoint string) (bool, error)
}

type checkAccessUseCase struct {
	checker AccessChecker
}

func NewCheckAccessUseCase(checker AccessChecker) *checkAccessUseCase {
	return &checkAccessUseCase{checker: checker}
}

func (c *checkAccessUseCase) Run(ctx context.Context, endpoint string) (bool, error) {
	return c.checker.Check(ctx, endpoint)
}
