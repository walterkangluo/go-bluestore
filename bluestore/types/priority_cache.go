package types

type Priority int

const (
	Pri0 Priority = iota
	Pri1
	Pri2
	Pri3
	Last = Pri3
)

type PriCache interface {
	RequestCacheBytes(pri Priority, chunkBytes uint64) int64

	// if Priority == -1, get all priorities
	GetCacheBytes(pri Priority) int64

	SetCacheBytes(pri Priority, bytes int64)

	AddCacheBytes(pri Priority, bytes int64)

	CommitCacheBytes() int64

	GetCacheRatio() float64

	SetCacheRatio(ratio float64)

	GetCacheName() string
}
