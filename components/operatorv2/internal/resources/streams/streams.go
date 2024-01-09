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
	stackName string, streamDirectory string, opts ...core.ObjectMutator[*v1beta1.Stream]) error {
	streamFiles, err := fs.ReadDir(streamDirectory)
	if err != nil {
		return err
	}

	for _, file := range streamFiles {
		streamContent, err := fs.ReadFile(streamDirectory + "/" + file.Name())
		if err != nil {
			return err
		}

		sanitizedName := strings.ReplaceAll(file.Name(), "_", "-")

		opts = append(opts, func(stream *v1beta1.Stream) {
			stream.Spec.Data = string(streamContent)
			stream.Spec.Stack = stackName
		})
		_, _, err = core.CreateOrUpdate[*v1beta1.Stream](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", stackName, sanitizedName),
		}, opts...)
		if err != nil {
			return errors.Wrap(err, "creating stream")
		}
	}

	return nil
}
