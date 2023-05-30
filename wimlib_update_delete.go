package cwimkit

const (
	DeleteFlagForce int32 = 1 << iota
	DeleteFlagRecursive
)

type updateDeleteCommandStruct struct {
	Operation UpdateOp
	PathPtr   Pointer
	Flags     int32
}

type UpdateDeleteCommand struct {
	Path  string
	Flags int32
}
