package compat

import "github.com/pawnkit/pawntest/internal/backend"

type Case struct {
	Name            string
	AMX             []byte
	ReferenceEngine string
	ReferenceTier   string
	PublicArgs      map[string][]backend.Cell
	Natives         map[string]backend.NativeFunc
}
