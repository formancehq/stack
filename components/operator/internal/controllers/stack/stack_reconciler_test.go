package stack_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/formancehq/operator/internal/modules/auth"

	"github.com/davecgh/go-spew/spew"
	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/modules"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
	batchv1 "k8s.io/api/batch/v1"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

func init() {
	stackv1beta3.ClientSecretGenerator = func() string {
		return "mocked-secret"
	}
	auth.RSAKeyGenerator = func() string {
		return "fake-rsa-key"
	}
	modules.CreatePostgresDatabase = func(ctx context.Context, dsn, dbName string) error {
		return nil
	}
}

var _ = Describe("Check stack deployment", func() {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			spew.Dump(e)
			panic(e)
		}
	}()
	ls, err := os.ReadDir("testdata")
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range ls {
		dirEntry := dirEntry
		if !dirEntry.IsDir() {
			continue
		}

		Context(dirEntry.Name(), func() {
			dirName := dirEntry.Name()
			name := strings.ReplaceAll(dirName, ".", "-")

			tmp := make(map[string]any)

			// Get Version
			// todo: factorize 3 occurences of the same code
			versionsFile, err := os.ReadFile(filepath.Join("testdata", dirName, "versions.yaml"))
			if err != nil {
				panic(err)
			}
			if err := yaml.Unmarshal(versionsFile, &tmp); err != nil {
				panic(err)
			}
			data, err := json.Marshal(tmp)
			if err != nil {
				panic(err)
			}
			versions := &stackv1beta3.Versions{}
			if err := json.Unmarshal(data, versions); err != nil {
				panic(err)
			}
			versions.Name = name

			// Get Configuration
			configurationFile, err := os.ReadFile(filepath.Join("testdata", dirName, "configuration.yaml"))
			if err != nil {
				panic(err)
			}
			if err := yaml.Unmarshal(configurationFile, &tmp); err != nil {
				panic(err)
			}
			data, err = json.Marshal(tmp)
			if err != nil {
				panic(err)
			}
			configuration := &stackv1beta3.Configuration{}
			if err := json.Unmarshal(data, configuration); err != nil {
				panic(err)
			}
			configuration.Name = name

			// Get Stack
			stackFile, err := os.ReadFile(filepath.Join("testdata", dirName, "stack.yaml"))
			if err != nil {
				panic(err)
			}
			if err := yaml.Unmarshal(stackFile, &tmp); err != nil {
				panic(err)
			}
			data, err = json.Marshal(tmp)
			if err != nil {
				panic(err)
			}
			stack := &stackv1beta3.Stack{}
			if err := json.Unmarshal(data, stack); err != nil {
				panic(err)
			}
			stack.Name = name

			// Launch test
			Context(fmt.Sprintf("with config from dir '%s'", dirEntry.Name()), func() {
				BeforeEach(func() {
					stack.Spec.Seed = configuration.Name
					stack.Spec.Versions = versions.Name
					stack.Spec.Stargate = &stackv1beta3.StackStargateConfig{}
				})
				JustBeforeEach(func() {
					Expect(Create(configuration)).To(Succeed())
					Expect(Create(versions)).To(Succeed())
					Expect(Create(stack)).To(Succeed())
					Eventually(func() bool {
						Expect(Get(types.NamespacedName{
							Namespace: stack.GetNamespace(),
							Name:      stack.GetName(),
						}, stack)).To(BeNil())
						return stack.IsReady()
					}).WithTimeout(10 * time.Second).Should(BeTrue())
				})
				JustAfterEach(func() {
					Expect(Delete(stack)).To(Succeed())
					Expect(Delete(configuration)).To(Succeed())
					Expect(Delete(versions)).To(Succeed())
				})
				It("should be ok", func() {
					verifyResources(stack, filepath.Join(dirName, "results"))
				})
			})
		})
	}
})

func verifyResources(stack *stackv1beta3.Stack, directory string) {
	if value, ok := os.LookupEnv("UPDATE_TEST_DATA"); ok && (value == "true" || value == "1") {
		updateTestingData(stack, directory)
		return
	}
	testDataResourcesDir := filepath.Join("testdata", directory)

	entries, err := os.ReadDir(testDataResourcesDir)
	Expect(err).ToNot(HaveOccurred())

	for _, gvkEntry := range entries {

		resourceDirFilename := filepath.Join(testDataResourcesDir, gvkEntry.Name())
		resourceEntries, err := os.ReadDir(resourceDirFilename)
		Expect(err).ToNot(HaveOccurred())

		for _, resourceEntry := range resourceEntries {
			gvkParts := strings.SplitN(gvkEntry.Name(), "-", 3)

			kind, err := k8sClient.RESTMapper().ResourceSingularizer(gvkParts[0])
			Expect(err).ToNot(HaveOccurred())

			actualResource := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind": kind,
					"apiVersion": v1.GroupVersion{
						Group:   gvkParts[1],
						Version: gvkParts[2],
					}.String(),
				},
			}
			Expect(Get(types.NamespacedName{
				Namespace: stack.Name,
				Name:      strings.TrimSuffix(resourceEntry.Name(), ".yaml"),
			}, actualResource)).To(BeNil())

			resourceEntryFilename := filepath.Join(resourceDirFilename, resourceEntry.Name())

			expectedContent, err := os.ReadFile(resourceEntryFilename)
			Expect(err).ToNot(HaveOccurred())

			expectedContentAsMap := make(map[string]any)
			Expect(yaml.Unmarshal(expectedContent, &expectedContentAsMap)).To(BeNil())

			expectedContentAsJSON, err := json.Marshal(expectedContentAsMap)
			Expect(err).ToNot(HaveOccurred())

			expectedResource := &unstructured.Unstructured{}
			Expect(json.Unmarshal(expectedContentAsJSON, &expectedResource)).To(BeNil())

			actualResourceSpec := actualResource.UnstructuredContent()
			expectedResourceSpec := expectedResource.UnstructuredContent()

			ignored := []string{
				`["clusterIP"]`,
				`["clusterIPs"]`,
				`["kind"]`,
				`["managedFields"]`,
				`["uid"]`,
				`["resourceVersion"]`,
				`["creationTimestamp"]`,
				`["ownerReferences"]`,
				`["generation"]`,
				`["lastTransitionTime"]`,
				`["lastUpdateTime"]`,
				`["controller-uid"]`,
			}

			if diff := cmp.Diff(expectedResourceSpec, actualResourceSpec,
				// Filter ignored pass
				cmp.FilterPath(func(path cmp.Path) bool {
					for _, ignored := range ignored {
						if ignored == path.Last().String() {
							return true
						}
					}
					return false
				}, cmp.Ignore()),
				// Ignore POSTGRES_PORT as it is not stable
				cmp.FilterValues(func(f, f2 any) bool {
					fAsMap, ok := f.(map[string]any)
					if !ok {
						return false
					}
					fName, ok := fAsMap["name"]
					if !ok {
						return false
					}
					f2AsMap, ok := f2.(map[string]any)
					if !ok {
						return false
					}
					f2Name, ok := f2AsMap["name"]
					if !ok {
						return false
					}
					return fName.(string) == f2Name.(string) && strings.HasSuffix(fName.(string), "POSTGRES_PORT")
				}, cmp.Ignore()),
			); diff != "" {
				msg := fmt.Sprintf("Expected content for resource %s not matching", resourceEntry.Name())
				msg += "\n" + diff
				Fail(msg)
			}

			ptrBool := func(v bool) *bool {
				return &v
			}

			Expect(actualResource.GetOwnerReferences()).To(Equal([]v1.OwnerReference{{
				APIVersion:         stackv1beta3.GroupVersion.String(),
				Kind:               "Stack",
				Name:               stack.Name,
				UID:                stack.UID,
				Controller:         ptrBool(true),
				BlockOwnerDeletion: ptrBool(true),
			}}))
		}
	}
}

func updateTestingData(stack *stackv1beta3.Stack, directory string) {
	dynamic := dynamic.NewForConfigOrDie(restConfig)

	gvks := []schema.GroupVersionResource{
		{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		},
		{
			Group:    batchv1.GroupName,
			Version:  batchv1.SchemeGroupVersion.Version,
			Resource: "cronjobs",
		},
		{
			Group:    batchv1.GroupName,
			Version:  batchv1.SchemeGroupVersion.Version,
			Resource: "jobs",
		},
		{
			Group:    "",
			Version:  "v1",
			Resource: "configmaps",
		},
		{
			Group:    "",
			Version:  "v1",
			Resource: "services",
		},
		{
			Group:    "",
			Version:  "v1",
			Resource: "secrets",
		},
		{
			Group:    networkingv1.GroupName,
			Version:  "v1",
			Resource: "ingresses",
		},
		{
			Group:    "stack.formance.com",
			Version:  "v1beta3",
			Resource: "migrations",
		},
	}
	for _, gvk := range gvks {
		list, err := dynamic.Resource(gvk).Namespace(stack.Name).List(ctx, v1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		groupDir := filepath.Join("testdata", directory, fmt.Sprintf("%s-%s-%s",
			gvk.Resource, gvk.Group, gvk.Version))
		Expect(os.RemoveAll(groupDir)).To(Succeed())
		Expect(os.MkdirAll(groupDir, os.ModePerm)).To(BeNil())

		for _, item := range list.Items {

			sampleFile, err := os.Create(fmt.Sprintf("%s/%s.yaml", groupDir, item.GetName()))
			Expect(err).ToNot(HaveOccurred())

			itemAsJson, err := json.Marshal(item.UnstructuredContent())
			Expect(err).ToNot(HaveOccurred())

			itemAsMap := make(map[string]any)
			Expect(json.Unmarshal(itemAsJson, &itemAsMap)).To(BeNil())

			clearUnstructuredContent(itemAsMap)

			Expect(yaml.NewEncoder(sampleFile).Encode(itemAsMap)).To(BeNil())
		}
	}
}

func clearUnstructuredContent(v map[string]any) {
	deleteFromSpec(v, "metadata.managedFields")
	deleteFromSpec(v, "metadata.uid")
	deleteFromSpec(v, "metadata.resourceVersion")
	deleteFromSpec(v, "metadata.creationTimestamp")
	deleteFromSpec(v, "metadata.ownerReferences")
	deleteFromSpec(v, "metadata.labels.controller-uid")
	deleteFromSpec(v, "spec.clusterIP")
	deleteFromSpec(v, "spec.clusterIPs")
	deleteFromSpec(v, "spec.selector.matchLabels.controller-uid")
	deleteFromSpec(v, "status.conditions[].lastTransitionTime")
	deleteFromSpec(v, "status.conditions[].lastUpdateTime")
	deleteFromSpec(v, "spec.template.metadata.labels.controller-uid")
}

func deleteFromSpec(m map[string]any, key string) {
	deleteFromMap(m, strings.Split(key, "."))
}

func deleteFromMap(m map[string]any, keys []string) {
	if strings.HasSuffix(keys[0], "[]") {
		key := keys[0][:len(keys[0])-2]
		v, ok := m[key]
		if !ok {
			return
		}
		deleteFromArray(v.([]any), keys[1:])
		return
	}
	if len(keys) == 1 {
		delete(m, keys[0])
		return
	}
	v, ok := m[keys[0]]
	if !ok {
		return
	}
	deleteFromMap(v.(map[string]any), keys[1:])
}

func deleteFromArray(items []any, keys []string) {
	for _, item := range items {
		deleteFromMap(item.(map[string]any), keys)
	}
}
