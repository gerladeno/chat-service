package utils

import (
	"context"
	"time"
)

func SlicesCollide[T comparable](s1 []T, s2 ...T) bool {
	m := make(map[T]struct{})
	for _, s := range s1 {
		m[s] = struct{}{}
	}
	var ok bool
	for _, s := range s2 {
		if _, ok = m[s]; ok {
			return true
		}
	}
	return false
}

func Sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
