package cwimkit

import (
	"context"
	"encoding/binary"
	"io"

	"go.uber.org/multierr"

	"github.com/go-wimkit/cwimkit/internal/csize"
)

const (
	AddNtfs int32 = 1 << iota
	AddDereference
	AddVerbose
	AddBoot
	AddUnixData
	AddNoAcls
	AddStrictAcls
	AddExcludeVerbose
	AddRPFix
	AddNoRPFix
	AddNoUnsupportedExclude
	AddWinConfig
	AddWimBoot
	AddNoReplace
	AddTestFileExclusion
	AddSnapshot
	AddFilePathsUnneeded
)

type updateAddCommandStruct struct {
	Operation     UpdateOp
	SourcePathPtr Pointer
	TargetPathPtr Pointer
	ConfigFilePtr Pointer
	Flags         int32
}

func (c *updateAddCommandStruct) WriteTo(wr io.Writer) (int64, error) {
	return int64(csize.OfMust(c)), binary.Write(wr, binary.LittleEndian, c)
}

func (c *updateAddCommandStruct) Free(ctx context.Context, api *API) error {
	return multierr.Combine(
		api.free(ctx, c.SourcePathPtr),
		api.free(ctx, c.TargetPathPtr),
		api.free(ctx, c.ConfigFilePtr),
	)
}

func newUpdateAddCommandStruct(ctx context.Context, api *API, c UpdateAddCommand) (*updateAddCommandStruct, error) {
	var err error

	s := &updateAddCommandStruct{
		Operation: UpdateOpAdd,
		Flags:     c.Flags,
	}

	s.SourcePathPtr, err = api.dupZString(ctx, &c.SourcePath)
	if err != nil {
		// TODO: use exit stack here to clear allocated pointers.
		return nil, err
	}

	s.TargetPathPtr, err = api.dupZString(ctx, &c.TargetPath)
	if err != nil {
		// TODO: use exit stack here to clear allocated pointers.
		return nil, err
	}

	s.ConfigFilePtr, err = api.dupZString(ctx, c.ConfigFile)
	if err != nil {
		// TODO: use exit stack here to clear allocated pointers.
		return nil, err
	}

	return s, nil
}

type UpdateAddCommand struct {
	SourcePath string
	TargetPath string
	ConfigFile *string
	Flags      int32
}

func (c UpdateAddCommand) WriterTo(ctx context.Context, api *API) (io.WriterTo, error) {
	return newUpdateAddCommandStruct(ctx, api, c)
}
