package stack_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	stackv1beta3 "github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/handlers"
	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
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
	handlers.RSAKeyGenerator = func() string {
		return "fake-rsa-key"
	}
}

var _ = Describe("When creating a stack", func() {
	var (
		configuration *stackv1beta3.Configuration
		versions      *stackv1beta3.Versions
		stack         = &stackv1beta3.Stack{}
	)
	BeforeEach(func() {
		configuration = NewDumbConfiguration()
		versions = NewDumbVersions()
		*stack = stackv1beta3.Stack{
			ObjectMeta: v1.ObjectMeta{
				Name: "stack1",
			},
			Spec: stackv1beta3.StackSpec{
				Seed:     configuration.Name,
				Versions: versions.Name,
				Host:     "example.net",
				Scheme:   "http",
			},
		}

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
		}).WithTimeout(5 * time.Second).Should(BeTrue())
	})
	JustAfterEach(func() {
		Expect(Delete(stack)).To(Succeed())
		Expect(Delete(configuration)).To(Succeed())
		Expect(Delete(versions)).To(Succeed())
	})
	It("should be ok", func() {
		verifyResources(stack, "multipod")
	})
	Context("with light mode", func() {
		BeforeEach(func() {
			stack.Name = "stack2"
			configuration.Spec.LightMode = true
		})
		It("should be ok", func() {
			verifyResources(stack, "monopod")
		})
	})
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
	}
	for _, gvk := range gvks {
		list, err := dynamic.Resource(gvk).Namespace(stack.Name).List(ctx, v1.ListOptions{})
		Expect(err).ToNot(HaveOccurred())

		for _, item := range list.Items {
			groupDir := filepath.Join("testdata", directory, fmt.Sprintf("%s-%s-%s",
				gvk.Resource, gvk.Group, gvk.Version))

			Expect(os.MkdirAll(groupDir, os.ModePerm)).To(BeNil())
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
	metadata := v["metadata"].(map[string]any)
	delete(metadata, "managedFields")
	delete(metadata, "uid")
	delete(metadata, "resourceVersion")
	delete(metadata, "creationTimestamp")
	delete(metadata, "ownerReferences")
}
