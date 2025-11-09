package tangent_sdk

import (
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"
	"go.bytecodealliance.org/cm"
)

type ArenaBuilder struct {
	// Flat node storage reused across messages.
	arena []mapper.Node

	// Object assembly state (single reusable buffers).
	fields   []mapper.Field // flat buffer of fields for all open objects
	objStart []int          // stack of start indexes into fields

	// Array assembly state (single reusable buffers).
	arrIdx   []mapper.Idx // flat buffer of element indexes for all open arrays
	arrStart []int        // stack of start indexes into arrIdx

	// frozen storage backing published cm.List views
	frozenFields []mapper.Field
	frozenIdx    []mapper.Idx

	// Per-call string table (not owned).
	st *StringTable
}

func NewArenaBuilder(capacity int) *ArenaBuilder {
	return &ArenaBuilder{
		arena:        make([]mapper.Node, 0, capacity),
		fields:       make([]mapper.Field, 0, 32),
		objStart:     make([]int, 0, 8),
		arrIdx:       make([]mapper.Idx, 0, 32),
		arrStart:     make([]int, 0, 8),
		frozenFields: make([]mapper.Field, 0, 64),
		frozenIdx:    make([]mapper.Idx, 0, 64),
	}
}

func (b *ArenaBuilder) UseStringTable(st *StringTable) { b.st = st }

func (b *ArenaBuilder) Reset() {
	b.arena = b.arena[:0]
	b.fields = b.fields[:0]
	b.objStart = b.objStart[:0]
	b.arrIdx = b.arrIdx[:0]
	b.arrStart = b.arrStart[:0]
	b.frozenFields = b.frozenFields[:0]
	b.frozenIdx = b.frozenIdx[:0]
	// NOTE: don't touch b.st; caller manages its lifetime and Reset.
}

func (b *ArenaBuilder) push(n mapper.Node) uint32 {
	b.arena = append(b.arena, n)
	return uint32(len(b.arena) - 1)
}

func (b *ArenaBuilder) Null() uint32       { return b.push(mapper.NodeNull()) }
func (b *ArenaBuilder) Bool(v bool) uint32 { return b.push(mapper.NodeBoolean(v)) }
func (b *ArenaBuilder) Int(v int64) uint32 { return b.push(mapper.NodeInteger(v)) }
func (b *ArenaBuilder) Float(v float64) uint32 {
	return b.push(mapper.NodeFloat(v))
}
func (b *ArenaBuilder) String(v string) uint32 {
	return b.push(mapper.NodeString_(v))
}

func (b *ArenaBuilder) Bytes(v []byte) uint32 {
	return b.push(mapper.NodeBytes(cm.ToList(v)))
}

// -------------------- Field helpers (key_id-based) --------------------

func (b *ArenaBuilder) FieldID(keyID uint32, val uint32) mapper.Field {
	return mapper.Field{Keyid: keyID, Val: mapper.Idx(val)}
}

func (b *ArenaBuilder) FieldKey(key string, val uint32) mapper.Field {
	if b.st == nil {
		panic("ArenaBuilder.FieldKey: no StringTable set; call UseStringTable first")
	}
	id := b.st.ID(key)
	return mapper.Field{Keyid: id, Val: mapper.Idx(val)}
}

// For codegen convenience:
func (b *ArenaBuilder) ObjectAddFieldID(keyID uint32, val uint32) {
	b.fields = append(b.fields, mapper.Field{Keyid: keyID, Val: mapper.Idx(val)})
}
func (b *ArenaBuilder) ObjectAddFieldKey(key string, val uint32) {
	if b.st == nil {
		panic("ArenaBuilder.ObjectAddFieldKey: no StringTable set; call UseStringTable first")
	}
	b.fields = append(b.fields, b.FieldKey(key, val))
}

// -------------------- Object builders --------------------

func (b *ArenaBuilder) ObjectStart() { b.objStart = append(b.objStart, len(b.fields)) }

func (b *ArenaBuilder) ObjectAdd(f mapper.Field) { b.fields = append(b.fields, f) }

// ObjectStartReserve is like ObjectStart but ensures capacity to append n fields
// without additional allocations (given the current slice state).
func (b *ArenaBuilder) ObjectStartReserve(n int) {
	b.objStart = append(b.objStart, len(b.fields))
	if n <= 0 {
		return
	}
	needed := n - (cap(b.fields) - len(b.fields))
	if needed > 0 {
		// Allocate exactly enough for the upcoming fields plus existing length.
		newCap := len(b.fields) + n
		newBuf := make([]mapper.Field, len(b.fields), newCap)
		copy(newBuf, b.fields)
		b.fields = newBuf
	}
}

func (b *ArenaBuilder) ObjectEnd() uint32 {
	top := len(b.objStart) - 1
	start := b.objStart[top]
	b.objStart = b.objStart[:top]

	// Copy child segment into frozen store (one memcpy, amortized).
	child := b.fields[start:len(b.fields)]
	base := len(b.frozenFields)
	b.frozenFields = append(b.frozenFields, child...)
	view := b.frozenFields[base : base+len(child)] // stable view

	// Pop child from scratch, but KEEP capacity so parent doesn't re-alloc.
	b.fields = b.fields[:start]

	return b.push(mapper.NodeObject(cm.ToList(view)))
}

// -------------------- Array builders --------------------

func (b *ArenaBuilder) ArrayStart() { b.arrStart = append(b.arrStart, len(b.arrIdx)) }

func (b *ArenaBuilder) ArrayAdd(idx uint32) {
	b.arrIdx = append(b.arrIdx, mapper.Idx(idx))
}

// ArrayStartReserve is like ArrayStart but ensures capacity to append n elements
// without additional allocations (given the current slice state).
func (b *ArenaBuilder) ArrayStartReserve(n int) {
	b.arrStart = append(b.arrStart, len(b.arrIdx))
	if n <= 0 {
		return
	}
	needed := n - (cap(b.arrIdx) - len(b.arrIdx))
	if needed > 0 {
		newCap := len(b.arrIdx) + n
		newBuf := make([]mapper.Idx, len(b.arrIdx), newCap)
		copy(newBuf, b.arrIdx)
		b.arrIdx = newBuf
	}
}

func (b *ArenaBuilder) ArrayEnd() uint32 {
	top := len(b.arrStart) - 1
	start := b.arrStart[top]
	b.arrStart = b.arrStart[:top]

	elems := b.arrIdx[start:len(b.arrIdx)]
	base := len(b.frozenIdx)
	b.frozenIdx = append(b.frozenIdx, elems...)
	view := b.frozenIdx[base : base+len(elems)]

	b.arrIdx = b.arrIdx[:start]
	return b.push(mapper.NodeArray(cm.ToList(view)))
}

// -------------------- Finalize (single arena + many roots) --------------------

// BuildBatch returns a deep-copied batch (safe; no aliasing).
func (b *ArenaBuilder) BuildBatch(st *StringTable, roots []uint32) mapper.Batchout {
	idxs := make([]mapper.Idx, len(roots))
	for i := range roots {
		idxs[i] = mapper.Idx(roots[i])
	}
	return mapper.Batchout{
		Strings: cm.ToList(st.Keys()),
		Arena:   cm.ToList(b.arena),
		Roots:   cm.ToList(idxs),
	}
}

// Test-only helper (keep behind a test build tag in real code).
func (b *ArenaBuilder) DebugArena() []mapper.Node { return b.arena }
