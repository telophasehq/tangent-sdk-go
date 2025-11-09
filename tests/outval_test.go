package tests

import (
	"testing"

	tangent "github.com/telophasehq/tangent-sdk-go"
)

func TestMyStructAnon(t *testing.T) {
	x := MyStructAnon{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: struct {
			MyNestedInt int `json:"nested_int"`
		}{MyNestedInt: 7},
		MyList: []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 6 {
		t.Fatalf("expected 6 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested == nil {
		t.Fatalf("MyNested is not an object")
	}
	nfields := nested.Slice()
	if len(nfields) != 1 || nfields[0].Key != "MyNestedInt" {
		t.Fatalf("MyNested fields = %+v, want single field MyNestedInt", nfields)
	}
	if got := arena[int(nfields[0].Val)].Integer(); got == nil || *got != int64(7) {
		t.Fatalf("MyNested.MyNestedInt = %v, want 7", got)
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}

func TestMyStructAnonPtr(t *testing.T) {
	x := MyStructAnonPtr{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: &struct {
			MyNestedInt int `json:"nested_int"`
		}{MyNestedInt: 7},
		MyList: []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 6 {
		t.Fatalf("expected 6 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested == nil {
		t.Fatalf("MyNested is not an object")
	}
	nfields := nested.Slice()
	if len(nfields) != 1 || nfields[0].Key != "MyNestedInt" {
		t.Fatalf("MyNested fields = %+v, want single field MyNestedInt", nfields)
	}
	if got := arena[int(nfields[0].Val)].Integer(); got == nil || *got != int64(7) {
		t.Fatalf("MyNested.MyNestedInt = %v, want 7", got)
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}

func TestMyStructAnonPtrNil(t *testing.T) {
	x := MyStructAnonPtr{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyList:   []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 4 {
		t.Fatalf("expected 4 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested != nil {
		t.Fatalf("MyNested is not nil")
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}

func TestMyStruct(t *testing.T) {
	y := MyNested{
		MyNestedInt: 7,
	}
	x := MyStruct{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: y,
		MyList:   []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 6 {
		t.Fatalf("expected 6 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested == nil {
		t.Fatalf("MyNested is not an object")
	}
	nfields := nested.Slice()
	if len(nfields) != 1 || nfields[0].Key != "MyNestedInt" {
		t.Fatalf("MyNested fields = %+v, want single field MyNestedInt", nfields)
	}
	if got := arena[int(nfields[0].Val)].Integer(); got == nil || *got != int64(7) {
		t.Fatalf("MyNested.MyNestedInt = %v, want 7", got)
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}

func TestMyStructPtr(t *testing.T) {
	y := &MyNested{
		MyNestedInt: 7,
	}
	x := MyStructPtr{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyNested: y,
		MyList:   []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 6 {
		t.Fatalf("expected 6 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested == nil {
		t.Fatalf("MyNested is not an object")
	}
	nfields := nested.Slice()
	if len(nfields) != 1 || nfields[0].Key != "MyNestedInt" {
		t.Fatalf("MyNested fields = %+v, want single field MyNestedInt", nfields)
	}
	if got := arena[int(nfields[0].Val)].Integer(); got == nil || *got != int64(7) {
		t.Fatalf("MyNested.MyNestedInt = %v, want 7", got)
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}

func TestMyStructPtrNil(t *testing.T) {
	x := MyStructPtr{
		MyInt:    42,
		MyFloat:  3.5,
		MyString: "hello",
		MyList:   []string{"a", "b"},
	}

	ab := tangent.NewArenaBuilder(0)
	root := x.AppendToArena(ab)
	out := ab.Build(root)

	arena := out.Arena.Slice()
	rootNode := arena[int(out.Root)]
	obj := rootNode.Object()
	if obj == nil {
		t.Fatalf("root is not an object node")
	}
	fields := obj.Slice()
	if len(fields) != 4 {
		t.Fatalf("expected 4 root fields, got %d", len(fields))
	}

	// Build a map of field name -> idx for easier assertions.
	fmap := make(map[string]uint32, len(fields))
	for _, f := range fields {
		fmap[f.Key] = uint32(f.Val)
	}

	// MyInt
	if got := arena[int(fmap["MyInt"])].Integer(); got == nil || *got != int64(42) {
		t.Fatalf("MyInt = %v, want 42", got)
	}

	// MyFloat
	if got := arena[int(fmap["MyFloat"])].Float(); got == nil || *got != 3.5 {
		t.Fatalf("MyFloat = %v, want 3.5", got)
	}

	// MyString
	if got := arena[int(fmap["MyString"])].String_(); got == nil || *got != "hello" {
		t.Fatalf("MyString = %v, want %q", got, "hello")
	}

	// MyNested
	nested := arena[int(fmap["MyNested"])].Object()
	if nested != nil {
		t.Fatalf("MyNested is not nil")
	}

	// MyList
	arr := arena[int(fmap["MyList"])].Array()
	if arr == nil {
		t.Fatalf("MyList is not an array")
	}
	idxs := arr.Slice()
	if len(idxs) != 2 {
		t.Fatalf("MyList length = %d, want 2", len(idxs))
	}
	if s := arena[int(idxs[0])].String_(); s == nil || *s != "a" {
		t.Fatalf("MyList[0] = %v, want %q", s, "a")
	}
	if s := arena[int(idxs[1])].String_(); s == nil || *s != "b" {
		t.Fatalf("MyList[1] = %v, want %q", s, "b")
	}
}
