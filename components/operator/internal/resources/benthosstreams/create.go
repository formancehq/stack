package benthosstreams

import (
	"embed"
	"fmt"

	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/formancehq/operator/api/formance.com/v1beta1"
	. "github.com/formancehq/operator/internal/core"

	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
)

func LoadFromFileSystem(ctx Context, fs embed.FS,
	owner v1beta1.Module, streamDirectory, discr string, opts ...ObjectMutator[*v1beta1.BenthosStream]) error {
	streamFiles, err := fs.ReadDir(streamDirectory)
	if err != nil {
		return err
	}

	search := &v1beta1.Search{}
	hasSearch, err := HasDependency(ctx, owner.GetStack(), search)
	if err != nil {
		return err
	}

	if hasSearch {
		names := make([]string, 0)
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
					"service": string(owner.GetUID()) + "-" + discr,
				}),
				WithController[*v1beta1.BenthosStream](ctx.GetScheme(), owner))

			name := fmt.Sprintf("%s-%s", owner.GetStack(), sanitizedName)
			_, _, err = CreateOrUpdate[*v1beta1.BenthosStream](ctx, types.NamespacedName{
				Name: name,
			}, opts...)
			if err != nil {
				return errors.Wrap(err, "creating stream")
			}

			names = append(names, name)
		}

		// Clean potential orphan streams
		l := &v1beta1.BenthosStreamList{}
		if err := ctx.GetClient().List(ctx, l, client.MatchingLabels{
			"service": string(owner.GetUID()) + "-" + discr,
		}); err != nil {
			return err
		}

		for _, stream := range l.Items {
			if !collectionutils.Contains(names, stream.Name) {
				if err := ctx.GetClient().Delete(ctx, &stream); err != nil {
					return err
				}
			}
		}
	} else {
		if err := ctx.GetClient().DeleteAllOf(ctx, &v1beta1.BenthosStream{}, client.MatchingLabels{
			"service": string(owner.GetUID()) + "-" + discr,
		}); err != nil {
			return err
		}
	}

	return nil
}
