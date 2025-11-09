package tangent_sdk

import (
	"errors"
	"sync"

	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/log"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"

	"go.bytecodealliance.org/cm"
)

type Log = log.Logview

// Handler processes a batch and writes NDJSON output to emitter.
type ProcessLogs[T any] func(Log) (T, error)

var arenaPool = sync.Pool{
	New: func() any {
		return NewArenaBuilder(1024)
	},
}

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

	mapper.Exports.ProcessLogs = func(input cm.List[cm.Rep]) (res cm.Result[cm.List[mapper.Outval], cm.List[mapper.Outval], string]) {
		batch := append([]cm.Rep(nil), input.Slice()...)
		logs := make([]Log, len(batch))
		outvals := make([]mapper.Outval, 0, len(batch))
		ab := arenaPool.Get().(*ArenaBuilder)
		defer func() {
			ab.Reset()
			arenaPool.Put(ab)
		}()
		for i := range batch {
			logs[i] = Log(batch[i])

			out, err := handler(logs[i])
			logs[i].ResourceDrop()
			if err != nil {
				res.SetErr(err.Error())
				return
			}
			ab.Reset()
			if m, ok := any(out).(OutvalMarshaler); ok {
				root := m.AppendToArena(ab)
				outvals = append(outvals, ab.Build(root))
				continue
			}
		}

		res.SetOK(cm.ToList(outvals))
		return
	}
}
