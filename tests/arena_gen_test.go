package tests

import (
	"bytes"
	"encoding/json"
	"testing"

	tangent "github.com/telophasehq/tangent-sdk-go"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"
	"go.bytecodealliance.org/cm"
)

type arenaMarshaler interface {
	AppendToArena(ab *tangent.ArenaBuilder, st *tangent.StringTable) uint32
}

func makeFixtures() []any {
	return []any{
		MyStructAnon{
			MyInt: 42, MyFloat: 3.5, MyString: "hello",
			MyNested: struct {
				MyNestedInt int `json:"nested_int"`
			}{7},
			MyList: []string{"a", "b", "c"},
		},
		MyStructAnonPtr{
			MyInt: 42, MyFloat: 3.5, MyString: "hello",
			MyNested: &struct {
				MyNestedInt int `json:"nested_int"`
			}{7},
			MyList: []string{"a", "b", "c"},
		},
		MyStruct{
			MyInt: 42, MyFloat: 3.5, MyString: "hello",
			MyNested: MyNested{MyNestedInt: 7},
			MyList:   []string{"a", "b", "c"},
		},
		MyStructPtr{
			MyInt: 42, MyFloat: 3.5, MyString: "hello",
			MyNested: &MyNested{MyNestedInt: 7},
			MyList:   []string{"a", "b", "c"},
		},
	}
}

func TestArenaEqualsStdJSON(t *testing.T) {
	ab := tangent.NewArenaBuilder(128)
	st := tangent.NewStringTable(16)
	ab.UseStringTable(st)

	for _, v := range makeFixtures() {
		ab.Reset()
		st.Reset()

		m, ok := v.(arenaMarshaler)
		if !ok {
			t.Fatalf("%T does not implement AppendToArena", v)
		}

		root := m.AppendToArena(ab, st)

		// Build one batch_out: strings, arena, single root
		bo := mapper.Batchout{
			Strings: cmListStringsCopy(st.Keys()),
			Arena:   cmListNodesCopy(ab),
			Roots:   cmListIdxOne(root),
		}

		// Arena → JSON
		var got bytes.Buffer
		if err := writeJSONFromBatch(&got, bo); err != nil {
			t.Fatalf("writeJSONFromBatch: %v", err)
		}

		// std json for the source value
		wantOne, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("json.Marshal: %v", err)
		}
		want := append(wantOne, '\n') // NDJSON

		if !bytes.Equal(got.Bytes(), want) {
			t.Errorf("mismatch for %T\n got: %s\nwant: %s", v, got.Bytes(), want)
		}
	}
}

// small helpers to avoid importing cm in tests
func cmListStringsCopy(xs []string) cm.List[string] {
	cp := append([]string(nil), xs...)
	return cm.ToList(cp)
}
func cmListNodesCopy(ab *tangent.ArenaBuilder) cm.List[mapper.Node] {
	cp := append([]mapper.Node(nil), ab.DebugArena()...) // add DebugArena() that returns []mapper.Node
	return cm.ToList(cp)
}
func cmListIdxOne(root uint32) cm.List[mapper.Idx] {
	return cm.ToList([]mapper.Idx{mapper.Idx(root)})
}
