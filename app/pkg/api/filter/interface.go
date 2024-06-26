package filter

type Options interface {
	Limit() int
	AddField(name, operator, value, dataType string) error
	Fields() []Field
}
