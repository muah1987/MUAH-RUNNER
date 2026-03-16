package adapters

type Adapter interface {
Name() string
Platform() string
IsAvailable() bool
GetToken() string
}
