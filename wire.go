package tangent_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"sync"

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

	mapper.Exports.ProcessLogs = func(input cm.List[cm.Rep]) (res cm.Result[cm.List[uint8], cm.List[uint8], string]) {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer bufPool.Put(buf)

		items := append([]cm.Rep(nil), input.Slice()...)
		for i := range items {

			out, err := handler(Log(items[i]))
			if err != nil {
				res.SetErr(err.Error())
				return
			}

			Log(items[i]).ResourceDrop()

			err = json.NewEncoder(buf).Encode(out)
			if err != nil {
				res.SetErr(err.Error())
				return
			}

		}

		res.SetOK(cm.ToList(buf.Bytes()))
		return
	}
}
