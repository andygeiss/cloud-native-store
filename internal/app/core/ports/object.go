package ports

import "context"

type ObjectPort[K comparable, V any] interface {
	Delete(ctx context.Context, key K) (err error)
	Get(ctx context.Context, key K) (value V, err error)
	Put(ctx context.Context, key K, value V) (err error)
}
