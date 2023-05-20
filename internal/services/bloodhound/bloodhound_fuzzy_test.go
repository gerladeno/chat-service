//go:build go1.18

package bloodhound_test

import (
	"testing"

	"github.com/gerladeno/chat-service/internal/services/bloodhound"
)

// go test --fuzz=FuzzSearch -fuzztime=10s ./...

func FuzzSearch(f *testing.F) {
	f.Add("Мой логин такой же как на карте 5179-4279-0625-0126 (IVAN KUZYAKIN) – ivankuzyakin.")
	f.Add("Hello, my dear friend!\n🙌")

	f.Fuzz(func(t *testing.T, input string) {
		_, err := bloodhound.Search(input)
		if err != nil {
			t.Fatal(err)
		}
	})
}
