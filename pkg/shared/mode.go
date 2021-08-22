package shared

// PerformanceMode indicates which mode the code is operating under. If Normal,
// then the code is NOT using special hardware instructions and instead relying
// on portable Go code. If Fast, then the code IS using special hardware instructions
// that is platform dependent. Each package exports a func that can be used to debug
// or inspect the configuration at runtime.
type PerformanceMode int

const (
	Normal PerformanceMode = iota
	Fast
)

type CheckMode func() PerformanceMode