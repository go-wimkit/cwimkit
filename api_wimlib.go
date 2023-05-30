package cwimkit

import (
	"context"
	"os"
)

func (api *API) PrintErrors(ctx context.Context) error {
	errFile, err := api.fdOpen(ctx, int32(os.Stderr.Fd()), "w")
	if err != nil {
		return err
	}

	return api.callWimlibFn(ctx, api.vt.wimlibSetErrorFile, errFile.arg())
}
