package tests

import (
	"testing"

	tangent "github.com/telophasehq/tangent-sdk-go"
)

func BenchmarkArena_vs_JSON(b *testing.B) {
	x := MyStruct{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: MyNested{MyNestedInt: 7},
		MyList:   []string{"a", "b", "c", "d", "e"},
	}

	b.Run("arena-reuse", func(b *testing.B) {
		ab := tangent.NewArenaBuilder(1024)
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			ab.Reset()
			root := x.AppendToArena(ab)
			_ = ab.BuildView(root)
		}
	})
}
