package appctx

import (
	"context"
	"sync/atomic"
)

type WailsCtxProvider interface {
	WailsCtx() context.Context
}

var ctxValue atomic.Value

func RegisterCtcProvider(cp WailsCtxProvider) {
	ctxValue.CompareAndSwap(nil, cp)
}

func GetAppCtx() context.Context {
	return ctxValue.Load().(WailsCtxProvider).WailsCtx()
}
