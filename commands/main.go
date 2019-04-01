package commands

// Command represents the cli command structure
type Command struct {
	Name        string
	Aliases     []string
	Description string
	Action      func(args ...string)
}
