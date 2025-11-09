package tangent_sdk

// in package pluginout (importable by user plugins)
type OutvalMarshaler interface {
	// Append nodes into ab and return the root node index.
	AppendToArena(ab *ArenaBuilder) (root uint32)
}

// Scalars
func AppendString(ab *ArenaBuilder, s string) uint32 { return ab.String(s) }
func AppendInt(ab *ArenaBuilder, v int64) uint32     { return ab.Int(v) }
func AppendFloat(ab *ArenaBuilder, v float64) uint32 { return ab.Float(v) }
func AppendBool(ab *ArenaBuilder, v bool) uint32     { return ab.Bool(v) }
func AppendBytes(ab *ArenaBuilder, b []byte) uint32  { return ab.Bytes(b) }

// Slices of things that already implement AppendToArena
func AppendSlice[T interface{ AppendToArena(*ArenaBuilder) uint32 }](ab *ArenaBuilder, xs []T) uint32 {
	ab.ArrayStart()
	for i := range xs {
		ab.ArrayAdd(xs[i].AppendToArena(ab))
	}
	return ab.ArrayEnd()
}

// Map[string]T
func AppendStringMap[T interface{ AppendToArena(*ArenaBuilder) uint32 }](ab *ArenaBuilder, m map[string]T) uint32 {
	ab.ObjectStart()
	for k, v := range m {
		ab.ObjectAdd(ab.Field(k, v.AppendToArena(ab)))
	}
	return ab.ObjectEnd()
}
