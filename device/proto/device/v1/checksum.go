package devicev1

import (
	"context"

	"google.golang.org/protobuf/encoding/protojson"
)

const DefaultInvalidChecksum = "invalid-checksum"

type generator interface {
	GenerateChecksum(ctx context.Context, data []byte) (string, error)
}

func (r *DiagnosticsResponse) GenerateChecksum(ctx context.Context, gen generator) string {
	r.Checksum = ""
	data, err := protojson.Marshal(r)
	if err != nil {
		return DefaultInvalidChecksum
	}
	checksum, err := gen.GenerateChecksum(ctx, data)
	if err != nil {
		return DefaultInvalidChecksum
	}
	return checksum
}
