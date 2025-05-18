package events

import (
	"crypto/sha256"
)

var Discriminators = make(map[[8]byte]string)

func init() {
	for name := range Registry {
		hash := sha256.Sum256([]byte("event:" + name))
		var disc [8]byte
		copy(disc[:], hash[:8])
		Discriminators[disc] = name
	}
}
