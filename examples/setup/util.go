package setup

import (
	"math/rand"
	"time"

	"github.com/dank/rlapi"
)

const RandomPlayerID = rlapi.PlayerID("Steam|76561198085817112|0")

func RandString(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
