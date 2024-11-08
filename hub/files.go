package hub

import (
	"github.com/pkg/errors"
	"iter"
	"path"
	"strings"
)

// IterFileNames iterate over the file names stored in the repo.
// It doesn't trigger the downloading of the repo, only of the repo info.
func (r *Repo) IterFileNames() iter.Seq2[string, error] {
	// Download info and files.
	err := r.DownloadInfo(false)
	if err != nil {
		// Error downloading: yield error only.
		return func(yield func(string, error) bool) {
			yield("", err)
			return
		}
	}
	return func(yield func(string, error) bool) {
		for _, si := range r.info.Siblings {
			fileName := si.Name
			if path.IsAbs(fileName) || strings.Index(fileName, "..") != -1 {
				yield("", errors.Errorf("model %q contains illegal file name %q -- it cannot be an absolute path, nor contain \"..\"",
					r.ID, fileName))
				return
			}
			if !yield(fileName, nil) {
				return
			}
		}
		return
	}
}

// DownloadFiles downloads the repository files, and return the path to the downloaded files in the cache structure.
// The returned downloadPaths can be read, but shouldn't be modified, since there may be other programs using the same
// files.
func (r *Repo) DownloadFiles(files ...string) (downloadedPaths []string, err error) {
	if len(files) == 0 {
		return
	}

	// Get/create cacheDir.
	var cacheDir string
	cacheDir, err = r.repoCache()
	if err != nil {
		return nil, err
	}
	_ = cacheDir

	// Get commitHash for current revision.
	var commitHash string
	commitHash, err = r.readCommitHashForRevision()
	if err != nil {
		return nil, err
	}

	return
}

// DownloadFile is a shortcut to DownloadFiles with only one file.
func (r *Repo) DownloadFile(file string) (downloadedPath string, err error) {
	res, err := r.DownloadFiles(file)
	if err != nil {
		return "", err
	}
	return res[0], nil
}
