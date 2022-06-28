package lib

const (
	contextPassFileSuffix         = "_pass_ctx.go"
	instrumentationPassFileSuffix = "_pass_tracing.go"
)

func ExecutePasses(projectPath string, packagePattern string, rootFunctions []string, funcDecls map[string]bool, backwardCallGraph map[string][]string) {
	GlobalPropagateContext(projectPath, packagePattern, backwardCallGraph, rootFunctions, funcDecls, contextPassFileSuffix)
	GlobalInstrument(projectPath, packagePattern, string("")+contextPassFileSuffix, backwardCallGraph, rootFunctions, instrumentationPassFileSuffix)
}
