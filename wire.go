package tangent_sdk

import (
	"errors"
	"sync"

	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/log"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"

	"go.bytecodealliance.org/cm"
)

type Log = log.Logview

type ProcessLogs[T any] func(Log) (T, error)

type OutvalMarshaler interface {
	AppendToArena(ab *ArenaBuilder, st *StringTable) (root uint32)
}

var arenaPool = sync.Pool{New: func() any { return NewArenaBuilder(1024) }}
var stringsPool = sync.Pool{New: func() any { return NewStringTable(32) }}

// Wire connects metadata, probe selectors, and a handler to Tangent's ABI.
func Wire[T any](meta Metadata, selectors []Selector, handler ProcessLogs[T]) {
	if handler == nil {
		panic(errors.New("handler must not be nil"))
	}

	mapper.Exports.Metadata = func() mapper.Meta {
		return meta.ToMapper()
	}

	mapper.Exports.Probe = func() cm.List[mapper.Selector] {
		mapped := make([]mapper.Selector, len(selectors))
		for i := range selectors {
			mapped[i] = selectors[i].toMapper()
		}
		return cm.ToList(mapped)
	}

	mapper.Exports.ProcessLogs = func(input cm.List[cm.Rep]) (res cm.Result[mapper.BatchoutShape, mapper.Batchout, string]) {
		batch := append([]cm.Rep(nil), input.Slice()...)
		logs := make([]Log, len(batch))

		ab := arenaPool.Get().(*ArenaBuilder)
		st := stringsPool.Get().(*StringTable)
		ab.Reset()
		st.Reset()
		ab.UseStringTable(st)
		roots := make([]mapper.Idx, 0, len(batch))

		defer func() {
			ab.Reset()
			arenaPool.Put(ab)
			st.Reset()
			stringsPool.Put(st)
		}()
		for i := range batch {
			logs[i] = Log(batch[i])

			out, err := handler(logs[i])
			logs[i].ResourceDrop()
			if err != nil {
				res.SetErr(err.Error())
				return
			}

			if m, ok := any(out).(OutvalMarshaler); ok {
				root := m.AppendToArena(ab, st)
				roots = append(roots, mapper.Idx(root))
				continue
			}

		}

		bo := mapper.Batchout{
			Strings: cm.ToList(st.Keys()),
			Arena:   cm.ToList(ab.arena),
			Roots:   cm.ToList(roots),
		}

		res.SetOK(bo)
		return
	}
}
