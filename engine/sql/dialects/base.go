package dialects

import "strings"

type baseDialect struct {
	Dialect
	uri URI
	// quoter Quoter
}

func (bd *baseDialect) Init(dialect Dialect, uri *URI) error {
	bd.Dialect, bd.uri = dialect, *uri
	return nil
}

func (bd *baseDialect) QuoteExpr(str string) string {
	return strings.ReplaceAll(str, "\"", "\\\"")
}
