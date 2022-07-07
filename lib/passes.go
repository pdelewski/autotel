package lib

const (
	contextPassFileSuffix         = "_pass_ctx.go"
	instrumentationPassFileSuffix = "_pass_tracing.go"
)

func ExecutePasses(projectPath string,
	packagePattern string,
	rootFunctions []FuncDescriptor,
	funcDecls map[FuncDescriptor]bool,
	backwardCallGraph map[FuncDescriptor][]FuncDescriptor) {

	PropagateContext(projectPath,
		packagePattern,
		backwardCallGraph,
		rootFunctions,
		funcDecls,
		contextPassFileSuffix)

	Instrument(projectPath,
		packagePattern,
		string("")+contextPassFileSuffix,
		backwardCallGraph,
		rootFunctions,
		instrumentationPassFileSuffix)
}
