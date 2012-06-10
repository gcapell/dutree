package main

import (
	"fmt"
)

// Copied from http://golang.org/doc/effective_go.html

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	precision := 1
	switch {
	case b >= YB:
		return fmt.Sprintf("%.*fYB", precision, b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.*fZB", precision, b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.*fEB", precision, b/EB)
	case b >= PB:
		return fmt.Sprintf("%.*fPB", precision, b/PB)
	case b >= TB:
		return fmt.Sprintf("%.*fTB", precision, b/TB)
	case b >= GB:
		return fmt.Sprintf("%.*fGB", precision, b/GB)
	case b >= MB:
		return fmt.Sprintf("%.*fMB", precision, b/MB)
	case b >= KB:
		return fmt.Sprintf("%.*fKB", precision, b/KB)
	}
	return fmt.Sprintf("%.*fB", precision, b)
}
