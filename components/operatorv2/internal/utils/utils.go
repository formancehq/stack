package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/fs"
	"path/filepath"
	"strings"

	pkgError "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

var (
	ErrNoConfigurationFound   = pkgError.New("no configuration found")
	ErrMultipleInstancesFound = pkgError.New("multiple resources found")
)

func HashFromConfigMap(configMap *corev1.ConfigMap) string {
	digest := sha256.New()
	if err := json.NewEncoder(digest).Encode(configMap.Data); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(digest.Sum(nil))
}

//func GetSingleStackDependencyObject[T client.Object](ctx context.Context, scheme *runtime.Scheme, _client client.Client, stackName string) (T, error) {
//
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//	fmt.Println("find object")
//
//	var t T
//	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
//	kinds, _, err := scheme.ObjectKinds(t)
//	if err != nil {
//		return t, err
//	}
//	spew.Dump(kinds)
//
//	list := &unstructured.UnstructuredList{}
//	list.SetGroupVersionKind(kinds[0])
//	err = ctx.GetClient().List(ctx, list, client.MatchingFields{
//		".spec.stack": stackName,
//	})
//	if err != nil {
//		return t, err
//	}
//
//	switch len(list.Items) {
//	case 0:
//		return t, nil
//	case 1:
//		if err := runtime.DefaultUnstructuredConverter.
//			FromUnstructured(list.Items[0].UnstructuredContent(), t); err != nil {
//			return t, err
//		}
//		return t, nil
//	default:
//		return t, ErrMultipleInstancesFound
//	}
//}

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
