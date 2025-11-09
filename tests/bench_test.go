package tests

import (
	"bytes"
	"testing"

	tangent "github.com/telophasehq/tangent-sdk-go"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"
	"go.bytecodealliance.org/cm"
)

func BenchmarkArenaPipeline(b *testing.B) {
	x := MyStruct{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: MyNested{MyNestedInt: 7},
		MyList:   []string{"a", "b", "c", "d", "e"},
	}

	// arena build only
	b.Run("arena-build", func(b *testing.B) {
		b.ReportAllocs()
		ab := tangent.NewArenaBuilder(2056)
		st := tangent.NewStringTable(16)
		ab.UseStringTable(st)
		for i := 0; i < b.N; i++ {
			ab.Reset()
			st.Reset()
			_ = x.AppendToArena(ab, st)
		}
	})

	// arena → writer (simulate host)
	b.Run("arena->writer", func(b *testing.B) {
		b.ReportAllocs()
		ab := tangent.NewArenaBuilder(2056)
		st := tangent.NewStringTable(16)
		ab.UseStringTable(st)
		root := x.AppendToArena(ab, st)
		bo := mapper.Batchout{
			Strings: cmListStringsCopy(st.Keys()),
			Arena:   cm.ToList(append([]mapper.Node(nil), ab.DebugArena()...)),
			Roots:   cm.ToList([]mapper.Idx{mapper.Idx(root)}),
		}
		var buf bytes.Buffer
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf.Reset()
			_ = writeJSONFromBatch(&buf, bo)
		}
	})
}
