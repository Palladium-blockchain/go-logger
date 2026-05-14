package plain

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/Palladium-blockchain/go-logger/pkg/core"
)

type Formatter struct{}

func (f *Formatter) Format(record core.Record) ([]byte, error) {
	var b bytes.Buffer

	b.WriteString(string(record.Level))
	b.WriteString(" | ")
	b.WriteString(record.Message)

	for _, field := range record.Fields {
		b.WriteString(" ")
		b.WriteString(field.Key)
		b.WriteString(": ")
		if err := f.formatAny(&b, field.Value); err != nil {
			b.WriteString("<not-formattable>")
		}
		b.WriteString(" ")
	}

	b.WriteString("\n")

	return b.Bytes(), nil
}

func (f *Formatter) formatAny(b *bytes.Buffer, v any) error {
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
