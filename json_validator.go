package go_cake

type JSONValidator interface {
	Validate(item map[string]any) error
}
