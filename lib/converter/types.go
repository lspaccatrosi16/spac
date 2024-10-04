package converter

type Convertable interface {
	Convert(from string, to string) string
	InputFormat() []string
	OutputFormat() []string
	Name() string
}
