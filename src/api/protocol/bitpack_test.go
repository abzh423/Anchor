package protocol_test

import (
	"log"
	"math/rand"
	"testing"

	"github.com/golangminecraft/minecraft-server/src/api/protocol"
)

func TestPackInt64Array(t *testing.T) {
	value := protocol.PackInt64Array([]int64{
		rand.Int63(),
		rand.Int63(),
		rand.Int63(),
	}, 12)

	log.Println(value)
}
