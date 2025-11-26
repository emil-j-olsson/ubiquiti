package devicev1

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
)

type generator interface {
	GenerateChecksum(ctx context.Context, data []byte) (string, error)
}

func (r *DiagnosticsResponse) GenerateChecksum(ctx context.Context, gen generator) string {
	invalid := "invalid-checksum"
	r.Checksum = ""
	data, err := protojson.Marshal(r)
	if err != nil {
		return invalid
	}
	checksum, err := gen.GenerateChecksum(ctx, data)
	if err != nil {
		return invalid
	}
	return checksum
}
