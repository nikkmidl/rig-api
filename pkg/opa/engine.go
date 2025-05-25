package opa

import (
	"context"
	_ "embed"

	"github.com/nikkmidl/rig-api/pkg/config"
	"github.com/open-policy-agent/opa/v1/rego"
)

//go:embed policies/access.rego
var policy string

type Evaluator struct {
	query rego.PreparedEvalQuery
	ctx   context.Context
}
type blockedUsersInput struct {
	user          string
	blocked_users []string
}

func NewEvaluator(ctx context.Context) (*Evaluator, error) {
	r := rego.New(
		rego.Query("data.github.access"),
		rego.Module("access.rego", policy),
	)
	q, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}
	return &Evaluator{query: q, ctx: ctx}, nil
}

func (e *Evaluator) IsBlocked(user string) (err error, blocked bool) {
	// Prepare the input for the query
	input := blockedUsersInput{user: user, blocked_users: config.Config.BlockedUsers}

	// Evaluate the query with the input
	rs, err := e.query.Eval(e.ctx, rego.EvalInput(input))
	if err != nil || len(rs) == 0 {
		return err, false
	}

	// Check if the result contains a "deny" reason
	reasons, ok := rs[0].Expressions[0].Value.(map[string]any)
	if ok && len(reasons) > 0 {
		return err, reasons["deny"].(bool)
	}
	return err, false
}
