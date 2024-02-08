package benthosstreams

import (
	"embed"
	"fmt"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"

	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

func LoadFromFileSystem(ctx Context, fs embed.FS,
	owner v1beta1.Module, streamDirectory string, opts ...ObjectMutator[*v1beta1.BenthosStream]) error {
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

		opts = append(opts,
			func(stream *v1beta1.BenthosStream) error {
				stream.Spec.Data = string(streamContent)
				stream.Spec.Stack = owner.GetStack()
				stream.Spec.Name = strings.TrimSuffix(file.Name(), ".yaml")

				return nil
			},
			WithLabels[*v1beta1.BenthosStream](map[string]string{
				"service": LowerCamelCaseKind(ctx, owner),
			}),
			WithController[*v1beta1.BenthosStream](ctx.GetScheme(), owner))
		_, _, err = CreateOrUpdate[*v1beta1.BenthosStream](ctx, types.NamespacedName{
			Name: fmt.Sprintf("%s-%s", owner.GetStack(), sanitizedName),
		}, opts...)
		if err != nil {
			return errors.Wrap(err, "creating stream")
		}
	}

	return nil
}
