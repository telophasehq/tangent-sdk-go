package tangent_sdk

import (
	"github.com/telophasehq/tangent_sdk/internal/tangent/logs/mapper"
	"go.bytecodealliance.org/cm"
)

type ArenaBuilder struct {
	// Nodes for the whole arena (reused across messages)
	arena []mapper.Node

	// --- Object assembly state (single reusable buffers) ---
	fields   []mapper.Field // flat buffer of fields for all open objects
	objStart []int          // stack of start indexes into fields

	// --- Array assembly state (single reusable buffers) ---
	arrIdx   []mapper.Idx // flat buffer of element indexes for all open arrays
	arrStart []int        // stack of start indexes into arrIdx
}

func NewArenaBuilder(capacity int) *ArenaBuilder {
	return &ArenaBuilder{
		arena:    make([]mapper.Node, 0, capacity),
		fields:   make([]mapper.Field, 0, 32),
		objStart: make([]int, 0, 8),
		arrIdx:   make([]mapper.Idx, 0, 32),
		arrStart: make([]int, 0, 8),
	}
}

func (b *ArenaBuilder) Reset() {
	b.arena = b.arena[:0]
	b.fields = b.fields[:0]
	b.objStart = b.objStart[:0]
	b.arrIdx = b.arrIdx[:0]
	b.arrStart = b.arrStart[:0]
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

// Safe: copies the bytes so caller can reuse/modify original slice.
func (b *ArenaBuilder) Bytes(v []byte) uint32 {
	cp := append([]byte(nil), v...)
	return b.push(mapper.NodeBytes(cm.ToList(cp)))
}

// Unsafe/zero-copy variant if you control ownership & lifetime of v.
// func (b *ArenaBuilder) BytesNoCopy(v []byte) uint32 {
// 	return b.push(mapper.NodeBytes(cm.ToList(v)))
// }

func (b *ArenaBuilder) Field(key string, val uint32) mapper.Field {
	return mapper.Field{Key: key, Val: mapper.Idx(val)}
}

// ---------- Object builders (no per-object copies) ----------

func (b *ArenaBuilder) ObjectStart() {
	b.objStart = append(b.objStart, len(b.fields))
}

func (b *ArenaBuilder) ObjectAdd(f mapper.Field) {
	b.fields = append(b.fields, f)
}

func (b *ArenaBuilder) ObjectEnd() uint32 {
	// Take a view (slice) over the shared fields buffer.
	top := len(b.objStart) - 1
	start := b.objStart[top]
	b.objStart = b.objStart[:top]

	fs := b.fields[start:len(b.fields)]
	// IMPORTANT: We do NOT shrink b.fields here; we keep the data so the view remains valid.
	// All later appends go to the end of b.fields; we never mutate prior slots.

	return b.push(mapper.NodeObject(cm.ToList(fs)))
}

// ---------- Array builders (no per-array copies) ----------

func (b *ArenaBuilder) ArrayStart() {
	b.arrStart = append(b.arrStart, len(b.arrIdx))
}

func (b *ArenaBuilder) ArrayAdd(idx uint32) {
	b.arrIdx = append(b.arrIdx, mapper.Idx(idx))
}

func (b *ArenaBuilder) ArrayEnd() uint32 {
	top := len(b.arrStart) - 1
	start := b.arrStart[top]
	b.arrStart = b.arrStart[:top]

	elems := b.arrIdx[start:len(b.arrIdx)]
	// Same idea: keep elems alive by not truncating arrIdx.

	return b.push(mapper.NodeArray(cm.ToList(elems)))
}

// ---------- Finalize ----------

// BuildView returns an Outval that references the builder's internal slices (no copy).
// Call Reset() only AFTER the host has consumed the Outval.
func (b *ArenaBuilder) BuildView(root uint32) mapper.Outval {
	return mapper.Outval{
		Root:  mapper.Idx(root),
		Arena: mapper.Arena(cm.ToList(b.arena)),
	}
}

// Build returns a deep-copied Outval (safe but allocates).
func (b *ArenaBuilder) Build(root uint32) mapper.Outval {
	arenaCopy := make([]mapper.Node, len(b.arena))
	copy(arenaCopy, b.arena)
	return mapper.Outval{
		Root:  mapper.Idx(root),
		Arena: mapper.Arena(cm.ToList(arenaCopy)),
	}
}
