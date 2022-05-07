package cody_dir

import (
	"sync"
	"testing"
	"time"
)

var mu sync.Mutex
var m map[string]interface{}

func RWMaps() {
	mu.Lock()
	defer mu.Unlock()
	if m == nil {
		m = make(map[string]interface{})
	}
	m["haha"] = time.Now().UnixNano()
}

func Test_compare(t *testing.T) {
	RWMaps()
}

//go test -bench=. -benchmem
func BenchmarkConcatStringBytesBuffer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RWMaps()
	}
}
