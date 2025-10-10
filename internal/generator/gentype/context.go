package gentype

import (
	"context"
)

type ctxKey string

const (
	ctxKeyRecursionDepth ctxKey = "recursion_depth"

	maxRecursionDepth = 50
)

func ContextIncRecursionDepth(ctx context.Context) context.Context {
	depth := 0

	if v := ctx.Value(ctxKeyRecursionDepth); v != nil {
		depth = v.(int)
	}

	return context.WithValue(ctx, ctxKeyRecursionDepth, depth+1)
}

func ContextGetRecursionDepth(ctx context.Context) int {
	if v := ctx.Value(ctxKeyRecursionDepth); v != nil {
		return v.(int)
	}

	return 0
}

func ContextMustValidateRecursionDepth(ctx context.Context, issuer string) {
	depth := ContextGetRecursionDepth(ctx)

	if depth > maxRecursionDepth {
		panic(issuer + ": max recursion depth reached")
	}
}
