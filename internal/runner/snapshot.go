package runner

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pawnkit/pawntest/internal/backend"
)

type snapshotStore struct {
	path    string
	update  bool
	test    string
	values  map[string]string
	loadErr error
}

func newSnapshotStore(source string, update bool) *snapshotStore {
	base := filepath.Base(source)
	path := filepath.Join(filepath.Dir(source), "__snapshots__", base+".snap.json")
	store := &snapshotStore{path: path, update: update, values: map[string]string{}}

	data, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			store.loadErr = err
		}

		return store
	}

	if err := json.Unmarshal(data, &store.values); err != nil {
		store.loadErr = fmt.Errorf("read snapshots %s: %w", path, err)
	}

	return store
}

func registerSnapshotNative(vm backend.VM, state *nativeState, store *snapshotStore) error {
	return vm.RegisterNative("__pawntest_snapshot", func(ctx backend.NativeContext, params []backend.Cell) (backend.Cell, error) {
		if store.loadErr != nil {
			return 0, store.loadErr
		}

		if len(params) < 4 {
			return 0, nil
		}

		name := readStringParam(ctx, params, 0)
		actual := readStringParam(ctx, params, 1)
		key := store.test + "::" + name

		expected, exists := store.values[key]
		if exists && expected == actual {
			return 1, nil
		}

		if store.update {
			store.values[key] = actual
			if err := store.write(); err != nil {
				return 0, err
			}

			return 1, nil
		}

		message := fmt.Sprintf("snapshot %q is missing; rerun with --update-snapshots", name)
		if exists {
			message = fmt.Sprintf("snapshot %q differs\nexpected: %q\nactual:   %q", name, expected, actual)
		}

		setFailure(state, params, 2, message, ctx)

		return 0, nil
	})
}

func (store *snapshotStore) write() error {
	if err := os.MkdirAll(filepath.Dir(store.path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(store.values, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')

	temp, err := os.CreateTemp(filepath.Dir(store.path), ".snap-*")
	if err != nil {
		return err
	}

	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempPath, store.path); err == nil {
		return nil
	}

	if err := os.Remove(store.path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.Rename(tempPath, store.path)
}
