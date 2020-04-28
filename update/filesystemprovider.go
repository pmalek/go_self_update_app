package update

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"github.com/pmalek/proton_task/version"
)

type FileSystemProvider struct {
	binaryDir string

	// cache is a local version number to binary file path cache.
	cache sync.Map
}

func NewFileSystemProvider(binaryDir string) (*FileSystemProvider, error) {
	fo, err := os.Stat(binaryDir)
	if err != nil {
		return nil, errors.Wrapf(err,
			"problem with provided binary updates directory: %s", binaryDir)
	}

	if !fo.IsDir() {
		return nil, errors.Wrapf(err,
			"provided binary updates directory path (%s) is not a directory", binaryDir)
	}

	return &FileSystemProvider{
		binaryDir: binaryDir,
		cache:     sync.Map{},
	}, nil
}

func (fsp *FileSystemProvider) IsUpdateAvailable(v int) (newversion int, err error) {
	fis, err := ioutil.ReadDir(fsp.binaryDir)
	if err != nil {
		return 0, errors.Wrapf(err,
			"failed to walk the binary update directory: %s", fsp.binaryDir)
	}

	var availableVersions []int
	for _, info := range fis {
		if info.IsDir() {
			continue
		}

		v, err := version.GetNumber(info.Name())
		if err != nil {
			continue
		}

		availableVersions = append(availableVersions, v)
		fsp.cache.Store(v, filepath.Join(fsp.binaryDir, info.Name()))
	}

	// No new versions found.
	if len(availableVersions) == 0 {
		return 0, nil
	}

	sort.Slice(availableVersions, func(i, j int) bool {
		return availableVersions[i] > availableVersions[j]
	})

	newestAvailable := availableVersions[0]
	if newestAvailable > v {
		return newestAvailable, nil
	}

	return 0, nil
}

// NOTE:
// No directory tree traversal is done here. Maybe that should get changed?
// Does it make sense to update to a version without checking a particular
// version to be available?
func (fsp *FileSystemProvider) Update(v int) error {
	pathI, ok := fsp.cache.Load(v)
	if !ok {
		return errors.Errorf("requested version %d not available")
	}

	path, ok := pathI.(string)
	if !ok {
		return errors.Errorf("internal cache error: got %T expected string", pathI)
	}

	return errors.Wrapf(
		syscall.Exec(path, nil, os.Environ()),
		"failed to exec version %d from %s", v, path,
	)
}
