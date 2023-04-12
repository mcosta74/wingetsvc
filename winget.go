package wingetsvc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/exp/slog"
)

type ServiceInfo struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type WingetController interface {
	Search(ctx context.Context, term string) ([]ServiceInfo, error)
	Versions(ctx context.Context, packageId string) ([]string, error)
}

func NewWingetController(logger *slog.Logger) WingetController {
	var controller WingetController
	if runtime.GOOS == "windows" {
		controller = &windowsController{}
	} else {
		controller = &noopController{}
	}
	return loggingControllerMiddleware(logger, controller)
}

type noopController struct{}

func (c *noopController) Search(ctx context.Context, term string) ([]ServiceInfo, error) {
	return nil, nil
}

func (c *noopController) Versions(ctx context.Context, packageId string) ([]string, error) {
	return nil, nil
}

type windowsController struct{}

func (c *windowsController) Search(ctx context.Context, term string) ([]ServiceInfo, error) {
	output, err := c.call(ctx, "search", term)
	if err != nil {
		return nil, err
	}

	return parseSearchOutput(output)
}

func (c *windowsController) call(ctx context.Context, args ...string) ([]byte, error) {
	wingetPath := filepath.Join(
		os.Getenv("LOCALAPPDATA"),
		"\\Microsoft\\WindowsApps\\winget.exe",
	)

	r, w, _ := os.Pipe()

	attrs := &os.ProcAttr{
		Files: []*os.File{nil, w, os.Stderr},
	}

	proc, err := os.StartProcess(wingetPath, append([]string{wingetPath}, args...), attrs)
	if err != nil {
		return nil, fmt.Errorf("StartProcess() failed: %w", err)
	}

	state, err := proc.Wait()
	if err != nil {
		return nil, fmt.Errorf("Wait() failed: %w", err)
	}

	w.Close()

	bb := &bytes.Buffer{}
	io.Copy(bb, r)

	if state.ExitCode() != 0 {
		return nil, fmt.Errorf(strings.TrimPrefix(strings.TrimSpace(bb.String()), "\b-\b \r"))
	}
	return bb.Bytes(), nil
}

func (c *windowsController) Versions(ctx context.Context, packageId string) ([]string, error) {
	output, err := c.call(ctx, "show", packageId, "--versions")
	if err != nil {
		return nil, err
	}

	return parseVersionsOutput(output)
}

func loggingControllerMiddleware(logger *slog.Logger, next WingetController) WingetController {
	return &loggingControllerMw{
		logger: logger,
		next:   next,
	}
}

type loggingControllerMw struct {
	logger *slog.Logger
	next   WingetController
}

func (c *loggingControllerMw) Search(ctx context.Context, term string) (records []ServiceInfo, err error) {
	defer func() {
		c.logger.Info("Search called", "term", term, "records", len(records), "failure", bool(err != nil))
	}()
	records, err = c.next.Search(ctx, term)
	return
}

func (c *loggingControllerMw) Versions(ctx context.Context, packageId string) (records []string, err error) {
	defer func() {
		c.logger.Info("Search called", "packageId", packageId, "records", len(records), "failure", bool(err != nil))
	}()
	records, err = c.next.Versions(ctx, packageId)
	return
}
