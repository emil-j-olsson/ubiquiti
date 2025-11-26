package checksum

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

const DefaultChecksumTimeout = 5 * time.Second

type Generator struct {
	binaryPath string
	timeout    time.Duration
}

func NewGenerator(path string) *Generator {
	return &Generator{
		binaryPath: path,
		timeout:    DefaultChecksumTimeout,
	}
}

func (g *Generator) GenerateChecksum(ctx context.Context, data []byte) (string, error) {
	ectx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()
	cmd := exec.CommandContext(ectx, g.binaryPath)
	cmd.Stdin = bytes.NewReader(data)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("checksum binary failed: %w (stderr: %s)", err, stderr.String())
	}
	return stdout.String(), nil
}
