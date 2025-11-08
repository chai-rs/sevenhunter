package testutil

import "testing"

type ArrangeFn[S, I any] func(t *testing.T, service S, input I)
