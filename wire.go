package tangent_sdk

import (
	"bytes"
	"errors"
	"sync"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/log"
	"github.com/telophasehq/tangent-sdk-go/internal/tangent/logs/mapper"

	"go.bytecodealliance.org/cm"
)

var (
	bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}
)

type Log = log.Logview

type ProcessLogs[T any] func(Log) (T, error)

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

	mapper.Exports.ProcessLogs = func(input cm.List[log.Logview]) (res cm.Result[cm.List[uint8], cm.List[uint8], string]) {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufPool.Put(buf)

		var jw jwriter.Writer

		items := append([]log.Logview(nil), input.Slice()...)
		for _, lv := range items {
			out, err := handler(lv)
			if err != nil {
				res.SetErr(err.Error())
				return
			}

			lv.ResourceDrop()

			if out_marshal, ok := any(out).(easyjson.Marshaler); ok {
				out_marshal.MarshalEasyJSON(&jw)
				jw.RawByte('\n')
			} else {
				res.SetErr("output does not implement easyjson.Marshaler. Did you recompile?")
				return
			}
		}

		if jw.Error != nil {
			res.SetErr(jw.Error.Error())
			return
		}

		jw.DumpTo(buf)
		res.SetOK(cm.ToList(buf.Bytes()))
		return
	}
}
