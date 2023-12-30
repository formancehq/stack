package core

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/fs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func HashFromConfigMap(configMap *corev1.ConfigMap) string {
	digest := sha256.New()
	if err := json.NewEncoder(digest).Encode(configMap.Data); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(digest.Sum(nil))
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

type ObjectMutator[T any] func(t T)

func WithController[T client.Object](scheme *runtime.Scheme, owner client.Object) ObjectMutator[T] {
	return func(t T) {
		if !metav1.IsControlledBy(t, owner) {
			if err := controllerutil.SetControllerReference(owner, t, scheme); err != nil {
				panic(err)
			}
		}
	}
}

func WithAnnotations[T client.Object](newAnnotations map[string]string) ObjectMutator[T] {
	return func(t T) {
		annotations := t.GetAnnotations()
		if annotations == nil {
			annotations = make(map[string]string)
		}
		for k, v := range newAnnotations {
			annotations[k] = v
		}
		t.SetAnnotations(annotations)
	}
}

func CreateOrUpdate[T client.Object](ctx Context,
	key types.NamespacedName, mutators ...ObjectMutator[T]) (T, controllerutil.OperationResult, error) {

	var ret T
	ret = reflect.New(reflect.TypeOf(ret).Elem()).Interface().(T)
	ret.SetNamespace(key.Namespace)
	ret.SetName(key.Name)
	operationResult, err := controllerutil.CreateOrUpdate(ctx, ctx.GetClient(), ret, func() error {
		for _, mutate := range mutators {
			mutate(ret)
		}
		return nil
	})
	return ret, operationResult, err
}

func DeleteIfExists[T client.Object](ctx Context, name types.NamespacedName) error {
	var t T
	t = reflect.New(reflect.TypeOf(t).Elem()).Interface().(T)
	if err := ctx.GetClient().Get(ctx, name, t); err != nil {
		if client.IgnoreNotFound(err) == nil {
			return nil
		}
		return err
	}
	return ctx.GetClient().Delete(ctx, t)
}

func ExtractItemsFromList(list client.ObjectList) []client.Object {

	items := reflect.ValueOf(list).
		Elem().
		FieldByName("Items")

	ret := make([]client.Object, 0)
	for i := 0; i < items.Len(); i++ {
		ret = append(ret, items.Index(i).Addr().Interface().(client.Object))
	}

	return ret
}
