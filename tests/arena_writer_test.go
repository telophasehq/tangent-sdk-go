package tests

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"
)

func writeJSONString(dst *bytes.Buffer, s string) {
	dst.WriteByte('"')
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			if i > start {
				dst.WriteString(s[start:i])
			}
			dst.WriteByte('\\')
			dst.WriteByte(s[i])
			start = i + 1
		case '\n':
			if i > start {
				dst.WriteString(s[start:i])
			}
			dst.WriteString(`\n`)
			start = i + 1
		case '\r':
			if i > start {
				dst.WriteString(s[start:i])
			}
			dst.WriteString(`\r`)
			start = i + 1
		case '\t':
			if i > start {
				dst.WriteString(s[start:i])
			}
			dst.WriteString(`\t`)
			start = i + 1
		default:
			if s[i] < 0x20 {
				if i > start {
					dst.WriteString(s[start:i])
				}
				dst.WriteString(fmt.Sprintf(`\u%04X`, s[i]))
				start = i + 1
			}
		}
	}
	if start < len(s) {
		dst.WriteString(s[start:])
	}
	dst.WriteByte('"')
}

func writeJSONFromBatch(dst *bytes.Buffer, bo mapper.Batchout) error {
	strs := bo.Strings.Slice()
	arena := bo.Arena.Slice()
	roots := bo.Roots.Slice()

	var writeNode func(idx mapper.Idx)
	writeNode = func(idx mapper.Idx) {
		n := arena[int(idx)]
		switch {
		case n.Null():
			dst.WriteString("null")
		case n.Boolean() != nil:
			if *n.Boolean() {
				dst.WriteString("true")
			} else {
				dst.WriteString("false")
			}
		case n.Integer() != nil:
			dst.WriteString(fmt.Sprint(*n.Integer()))
		case n.Float() != nil:
			f := n.Float()
			if math.IsInf(*f, 0) || math.IsNaN(*f) {
				dst.WriteString("null")
			} else {
				dst.WriteString(fmt.Sprintf("%g", *f))
			}
		case n.String_() != nil:
			writeJSONString(dst, *n.String_())
		case n.Bytes() != nil:
			dst.WriteByte('"')
			dst.WriteString(base64.StdEncoding.EncodeToString(n.Bytes().Slice()))
			dst.WriteByte('"')
		case n.Array() != nil:
			dst.WriteByte('[')
			first := true
			for _, child := range n.Array().Slice() {
				if !first {
					dst.WriteByte(',')
				}
				first = false
				writeNode(child)
			}
			dst.WriteByte(']')
		case n.Object() != nil:
			dst.WriteByte('{')
			first := true
			for _, f := range n.Object().Slice() {
				if !first {
					dst.WriteByte(',')
				}
				first = false
				writeJSONString(dst, strs[int(f.Keyid)])
				dst.WriteByte(':')
				writeNode(f.Val)
			}
			dst.WriteByte('}')
		default:
			panic("unknown node")
		}
	}

	for _, r := range roots {
		writeNode(r)
		dst.WriteByte('\n') // NDJSON
	}
	return nil
}
