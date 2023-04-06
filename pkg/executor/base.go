package executor

type Executor interface {
	// Configure executor with the provided argument
	Configure(any) error

	// Execute entrypoint of the execution
	Execute(workspacePath string) error
}
