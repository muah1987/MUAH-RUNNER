package adapters

type GenericAdapter struct{ name string }

func NewGenericAdapter(name string) *GenericAdapter { return &GenericAdapter{name: name} }
func (a *GenericAdapter) Name() string              { return a.name }
func (a *GenericAdapter) Platform() string          { return "standalone" }
func (a *GenericAdapter) IsAvailable() bool         { return true }
func (a *GenericAdapter) GetToken() string          { return "" }
