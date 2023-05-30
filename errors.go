package cwimkit

type ModuleError string

func (e ModuleError) Error() string { return "cwimlib: " + string(e) }

var (
	ErrNoMemory   = ModuleError("out of memory")
	ErrOOBRead    = ModuleError("out of bounds read")
	ErrOOBWrite   = ModuleError("out of bounds write")
	ErrNoFunction = ModuleError("missing function")
)
