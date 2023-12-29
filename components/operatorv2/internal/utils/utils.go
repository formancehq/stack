package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/formancehq/operator/v2/internal/reconcilers"
	"io/fs"
	"path/filepath"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"

	pkgError "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrNotFound               = pkgError.New("no configuration found")
	ErrMultipleInstancesFound = pkgError.New("multiple resources found")
)

func HashFromConfigMap(configMap *corev1.ConfigMap) string {
	digest := sha256.New()
	if err := json.NewEncoder(digest).Encode(configMap.Data); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(digest.Sum(nil))
}

func GetSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx reconcilers.Context, stackName string) (OBJECT, error) {

	var t OBJECT
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(OBJECT)

	var list LIST
	list = reflect.New(reflect.TypeOf(list).Elem()).Interface().(LIST)

	err := ctx.GetClient().List(ctx, list, client.MatchingFields{
		".spec.stack": stackName,
	})
	if err != nil {
		return t, err
	}

	items := reflect.ValueOf(list).
		Elem().
		FieldByName("Items")

	switch items.Len() {
	case 0:
		return t, nil
	case 1:
		return items.Index(0).Addr().Interface().(OBJECT), nil
	default:
		return t, ErrMultipleInstancesFound
	}
}

func HasSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx reconcilers.Context, stackName string) (bool, error) {
	ret, err := GetSingleStackDependencyObject[LIST, OBJECT](ctx, stackName)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		return false, err
	}
	if reflect.ValueOf(ret).Elem().IsZero() {
		return false, nil
	}
	return true, nil
}
func RequireSingleStackDependencyObject[LIST client.ObjectList, OBJECT client.Object](ctx reconcilers.Context, stackName string) (OBJECT, error) {
	var ret OBJECT
	ret, err := GetSingleStackDependencyObject[LIST, OBJECT](ctx, stackName)
	if err != nil && !errors.Is(err, ErrMultipleInstancesFound) {
		return ret, err
	}
	if reflect.ValueOf(ret).Elem().IsZero() {
		return ret, ErrNotFound
	}
	return ret, nil
}

func CopyDir(f fs.FS, root, path string, ret *map[string]string) {
	dirEntries, err := fs.ReadDir(f, path)
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		dirEntryPath := filepath.Join(path, dirEntry.Name())
		if dirEntry.IsDir() {
			CopyDir(f, root, dirEntryPath, ret)
		} else {
			fileContent, err := fs.ReadFile(f, dirEntryPath)
			if err != nil {
				panic(err)
			}
			sanitizedPath := strings.TrimPrefix(dirEntryPath, root)
			sanitizedPath = strings.TrimPrefix(sanitizedPath, "/")
			(*ret)[sanitizedPath] = string(fileContent)
		}
	}
}
