package plain

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/Palladium-blockchain/go-logger/internal/syncx"
	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type Encoder struct {
	pool *syncx.Pool[*bytes.Buffer]
}

func NewEncoder() *Encoder {
	return &Encoder{
		pool: syncx.NewPool(func() *bytes.Buffer {
			return new(bytes.Buffer)
		}),
	}
}

func (e *Encoder) Encode(w io.Writer, record core.Record) error {
	b := e.pool.Get()
	defer e.pool.Put(b)

	b.WriteString(string(record.Level))
	b.WriteString(" | ")
	b.WriteString(record.Message)

	for _, field := range record.Fields {
		b.WriteString(" ")
		b.WriteString(field.Key)
		b.WriteString(": ")
		if err := e.encodeAny(b, field.Value); err != nil {
			b.WriteString("<not-encodable>")
		}
		b.WriteString(" ")
	}

	b.WriteString("\n")

	_, err := w.Write(b.Bytes())

	b.Reset()

	return err
}

func (e *Encoder) encodeAny(b *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case nil:
		b.WriteString("<nil>")
	case string:
		b.WriteString(x)
	case int:
		b.Write(strconv.AppendInt(nil, int64(x), 10))
	case int64:
		b.Write(strconv.AppendInt(nil, x, 10))
	case bool:
		b.Write(strconv.AppendBool(nil, x))
	case float64:
		b.Write(strconv.AppendFloat(nil, x, 'f', -1, 64))
	default:
		_, err := fmt.Fprint(b, x)
		return err
	}

	return nil
}
