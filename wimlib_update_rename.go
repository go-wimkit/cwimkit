package cwimkit

type updateRenameCommandStruct struct {
	Operation     UpdateOp
	SourcePathPtr Pointer
	TargetPathPtr Pointer
	Flags         int32
}

type UpdateRenameCommand struct {
	SourcePath string
	TargetPath string
	flags      int
}
