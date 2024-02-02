package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
	gomegaTypes "github.com/onsi/gomega/types"
	"github.com/pkg/errors"
)

type matchGoldenFile struct {
	paths []string
}

func (c matchGoldenFile) Match(actual interface{}) (success bool, err error) {

	value, ok := actual.(string)
	if !ok {
		return false, errors.New("should receive a string")
	}

	locationParts := []string{"testdata", "resources"}
	locationParts = append(locationParts, c.paths[:len(c.paths)-1]...)
	directory := filepath.Join(locationParts...)
	path := filepath.Join(directory, c.paths[len(c.paths)-1])

	updateTestData := os.Getenv("UPDATE_TEST_DATA")
	isUpdating := updateTestData == "1" || updateTestData == "true"
	if isUpdating {
		if err := os.MkdirAll(directory, 0744); err != nil {
			return false, errors.Wrapf(err, "creating directory %s", directory)
		}

		if err := os.WriteFile(path, []byte(value), 0600); err != nil {
			return false, errors.Wrapf(err, "writing file: %s", path)
		}

		return true, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return false, errors.Wrapf(err, "reading file: %s", path)
	}

	if diff := cmp.Diff(string(data), value); diff != "" {
		msg := fmt.Sprintf("Expected content for resource %s not matching", path)
		msg += "\n" + diff
		return false, errors.New(msg)
	}

	return true, nil
}

func (c matchGoldenFile) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("resource does not match resource at %s", filepath.Join(c.paths...))
}

func (c matchGoldenFile) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("resource should not match resource at %s", filepath.Join(c.paths...))
}

func MatchGoldenFile(paths ...string) gomegaTypes.GomegaMatcher {
	return &matchGoldenFile{
		paths: paths,
	}
}
