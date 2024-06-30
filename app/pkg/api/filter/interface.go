package filter

type Options interface {
	Limit() int
	AddField(name, operator string, values []string, dataType string) error
	Fields() []Field
}
