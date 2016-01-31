package unit

import (
	"errors"
	"fmt"
	"io"
)

// Context is a usage domain for Measurement values, it qualifies a unit,
// allowing it to be formatted differenty.
type Context struct {
	Name string
	*unit
	format string
}

var contexts = make(map[string]*Context)

func DefineContext(name, unit string, format string) (*Context, error) {
	if _, exists := contexts[name]; exists {
		return nil, errors.New("duplicate context: " + name)
	}
	ctx := &Context{name, get(unit), format}
	contexts[name] = ctx
	return ctx, nil
}

func Ctx(name string) *Context {
	return contexts[name]
}

func DeleteContext(c *Context) {
	delete(contexts, c.Name)
}

func (ctx Context) M(value float64, symbol string) Measurement {
	m := M(value, symbol)
	return ctx.Convert(m)
}

func (ctx Context) Convert(m Measurement) Measurement {
	return Measurement{m.value * m.factor / ctx.unit.factor, ctx.unit}
}

func (ctx Context) Format(wr io.Writer, m Measurement) {
	ctxm := ctx.Convert(m)
	fmt.Fprintf(wr, ctx.format, ctxm.Value(), ctxm.Symbol())
}

func (ctx Context) String(m Measurement) string {
	ctxm := ctx.Convert(m)
	return fmt.Sprintf(ctx.format, ctxm.Value(), ctxm.Symbol())
}
