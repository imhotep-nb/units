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

// DefineContext registers a new usage context for a unit. It narrows down the domain in 
// which the unit is used and defines what the default symbol is and how to format output.
// The name should be unique and is passed to Ctx(string) for lookup. An empty string is also
// allowed: it will create the Context but not register it for lookup. The caller should keep
// the reference somewhere. 
// The unit string is the default unit symbol and either it already exists or can be calculated. 
// The format string is a normal Go fmt string. Index [1] is the value and index [2] is the unit 
// symbol, e.g. "%[2]s %.2[1]f" to put the unit in front of the value. If both value and unit are 
// referenced in that order in the format string, then the indexes are not necessary, e.g. "%e%s".
func DefineContext(name, unit string, format string) (*Context, error) {
	if name == "" {
		return &Context{"", get(unit), format}, nil
	}
	if _, exists := contexts[name]; exists {
		return nil, errors.New("duplicate context: " + name)
	}
	ctx := &Context{name, get(unit), format}
	contexts[name] = ctx
	return ctx, nil
}

// Ctx looks up a Context by name and returns a reference to it.
// The return value is nil if the name was not registered with DefineContext.
func Ctx(name string) *Context {
	return contexts[name]
}

// DeleteContext unregisters the context.
func DeleteContext(c *Context) {
	delete(contexts, c.Name)
}

// M creates a new Measurement based on the Context. The value is converted to the unit defined
// in the Context.
func (ctx Context) M(value float64, symbol string) Measurement {
	m := M(value, symbol)
	return ctx.Convert(m)
}

// Convert converts a given measurement to the Context's default.
func (ctx Context) Convert(m Measurement) Measurement {
	return Measurement{m.value * m.factor / ctx.unit.factor, ctx.unit}
}

// Format writes a formatted version of the Measurement to the Writer.
func (ctx Context) Format(wr io.Writer, m Measurement) {
	ctxm := ctx.Convert(m)
	fmt.Fprintf(wr, ctx.format, ctxm.Value(), ctxm.Symbol())
}

// String returns a Measurement as string, formatted with the Context format string.
func (ctx Context) String(m Measurement) string {
	ctxm := ctx.Convert(m)
	return fmt.Sprintf(ctx.format, ctxm.Value(), ctxm.Symbol())
}
