// tangent_sdk/string_table.go
package tangent_sdk

type StringTable struct {
	ids  map[string]uint32 // key -> id
	keys []string          // id -> key
}

func NewStringTable(cap int) *StringTable {
	return &StringTable{
		ids:  make(map[string]uint32, cap),
		keys: make([]string, 0, cap),
	}
}

func (st *StringTable) Reset() {
	clear(st.ids)
	st.keys = st.keys[:0]
}

func (st *StringTable) ID(key string) uint32 {
	if id, ok := st.ids[key]; ok {
		return id
	}
	id := uint32(len(st.keys))
	st.ids[key] = id
	st.keys = append(st.keys, key)
	return id
}

func (st *StringTable) Keys() []string { return st.keys }
