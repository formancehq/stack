package streams

import (
	"embed"
	"fmt"
	"github.com/formancehq/operator/v2/api/v1beta1"
	"github.com/formancehq/operator/v2/internal/core"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"strings"
)

func LoadFromFileSystem(ctx core.Context, fs embed.FS,
	stackName string, streamDirectory string) error {
	streamFiles, err := fs.ReadDir(streamDirectory)
	if err != nil {
		return err
	}

	// TODO: Only if search enabled
	for _, file := range streamFiles {
		streamContent, err := fs.ReadFile(streamDirectory + "/" + file.Name())
		if err != nil {
			return err
		}

		sanitizedName := strings.ReplaceAll(file.Name(), "_", "-")

		_, _, err = core.CreateOrUpdate[*v1beta1.Stream](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stackName, sanitizedName),
		}, func(t *v1beta1.Stream) {
			t.Spec.Data = string(streamContent)
			t.Spec.Stack = stackName
		})
		if err != nil {
			return errors.Wrap(err, "creating stream")
		}
	}

	return nil
}
