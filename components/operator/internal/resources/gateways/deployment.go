package gateways

import (
	"archive/zip"
	"bytes"
	"github.com/formancehq/operator/api/formance.com/v1beta1"
	"github.com/formancehq/operator/internal/core"
	"github.com/formancehq/operator/internal/resources/deployments"
	"github.com/formancehq/operator/internal/resources/registries"
	"github.com/formancehq/operator/internal/resources/settings"
	"github.com/formancehq/stack/libs/events"
	"github.com/pkg/errors"
	"io"
	fs "io/fs"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
)

func createDeployment(ctx core.Context, stack *v1beta1.Stack,
	gateway *v1beta1.Gateway, caddyfileConfigMap *v1.ConfigMap,
	auditTopic *v1beta1.BrokerTopic, version string) error {

	env := GetEnvVars(gateway)
	otlpEnv, err := settings.GetOTELEnvVars(ctx, stack.Name, core.LowerCamelCaseKind(ctx, gateway))
	if err != nil {
		return err
	}
	env = append(env, otlpEnv...)
	env = append(env, core.GetDevEnvVars(stack, gateway)...)

	if stack.Spec.EnableAudit && auditTopic != nil {
		env = append(env, settings.GetBrokerEnvVars(auditTopic.Status.URI, stack.Name, "gateway")...)
	}

	image, err := registries.GetImage(ctx, stack, "gateway", version)
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")
	writer := zip.NewWriter(buf)
	eventsFS, err := fs.Sub(events.Computed, "generated")
	if err != nil {
		return err
	}

	if err := addFSToZipAtDir(writer, eventsFS, "assets/events/"); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	cm, _, err := core.CreateOrUpdate[*v1.ConfigMap](ctx, types.NamespacedName{
		Namespace: stack.Name,
		Name:      "assets",
	}, func(t *v1.ConfigMap) error {
		t.BinaryData = map[string][]byte{
			"assets.zip": buf.Bytes(),
		}

		return nil
	}, core.WithController[*v1.ConfigMap](ctx.GetScheme(), gateway))
	if err != nil {
		return err
	}

	_, err = deployments.CreateOrUpdate(ctx, gateway, "gateway",
		deployments.WithReplicasFromSettings(ctx, stack),
		deployments.WithMatchingLabels("gateway"),
		settings.ConfigureCaddy(caddyfileConfigMap, image, env, func(t *v1.Container) error {
			t.VolumeMounts = append(t.VolumeMounts, v1.VolumeMount{
				Name:      "assets",
				ReadOnly:  true,
				MountPath: "/assets",
			})
			return nil
		}),
		deployments.WithVolumes(v1.Volume{
			Name: "assets",
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: cm.Name,
					},
				},
			},
		}),
		deployments.WithTemplateAnnotations(map[string]string{
			"assets-hash": core.HashFromConfigMaps(cm),
		}),
	)

	return err
}

func addFSToZipAtDir(w *zip.Writer, fsys fs.FS, location string) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return errors.New("zip: cannot add non-regular file")
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = filepath.Join(location, name)
		h.Method = zip.Deflate
		fw, err := w.CreateHeader(h)
		if err != nil {
			return err
		}
		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(fw, f)
		return err
	})
}
