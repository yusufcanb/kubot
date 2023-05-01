package app

type RuntimeArgs struct {
	TopLevelSuiteName string
	Namespace         string
	Image             string
	Selector          string
	WorkspacePath     string
	BatchSize         int
}
