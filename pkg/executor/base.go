package executor

type Executor interface {
	Namespace(n string) *Executor
	Image(i string) *Executor
	Selector(s string) *Executor
	Workspace(w string) *Executor

	// Configure executor with the provided argument
	Configure(any) error

	// Execute entrypoint of the execution
	Execute() error
}
