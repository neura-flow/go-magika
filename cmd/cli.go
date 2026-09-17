package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/neura-flow/go-magika/pkg/magika"
	"github.com/neura-flow/go-magika/pkg/util"
)

// cli is a basic CLI that infers the content type of the files listed on the command line.
func cli(ctx context.Context, w io.Writer, args ...string) error {
	s, err := magika.NewScanner()
	if err != nil {
		return err
	}
	defer s.Close()

	absPaths := make([]string, 0, len(args))
	for _, p := range args {
		if absPath, _ := filepath.Abs(path.Clean(util.ExpandPath(p))); absPath != "" {
			if fi, err := os.Stat(absPath); err != nil && !os.IsNotExist(err) {
				absPaths = append(absPaths, absPath)
			} else if !fi.IsDir() {
				absPaths = append(absPaths, absPath)
			} else {
				absPaths = append(absPaths, filepath.Join(absPath, "**"))
			}
		}
	}

	globs, _ := util.Globs(util.OS(), absPaths)

	// For each filename given as argument, read the file and scan its content.
	for _, a := range globs {
		_, _ = fmt.Fprintf(w, "%s: ", a)
		b, err := os.ReadFile(a)
		if err != nil {
			_, _ = fmt.Fprintf(w, "%v\n", err)
			continue
		}
		ct, err := s.Scan(ctx, bytes.NewReader(b), len(b))
		if err != nil {
			_, _ = fmt.Fprintf(w, "scan: %v\n", err)
			continue
		}
		_, _ = fmt.Fprintf(w, "%s\n", ct.Label)
	}
	return nil
}
