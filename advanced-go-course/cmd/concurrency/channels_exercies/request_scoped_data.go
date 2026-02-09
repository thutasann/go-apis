package main

import (
	"context"
	"fmt"
)

type requestIDKey struct{}

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

func log(ctx context.Context, msg string) {
	id := ctx.Value(requestIDKey{})
	fmt.Printf("[request=%v] %s\n", id, msg)
}

func Request_Scoped_Data() {
	ctx := context.Background()
	ctx = withRequestID(ctx, "req-12345")

	log(ctx, "started request")
	log(ctx, "querying database")
	log(ctx, "sending response")
}
