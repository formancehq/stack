package delete

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/formancehq/operator/apis/stack/v1beta3"
	"github.com/formancehq/operator/internal/controllers/stack/storage/pg"
	operatorclient "github.com/formancehq/operator/pkg/client/v1beta3"
	sharedlogging "github.com/formancehq/stack/libs/go-libs/logging"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type StorageFn func(fileName string, data []byte) error

func (fn StorageFn) PutFile(fileName string, data []byte) error {
	return fn(fileName, data)
}

type File struct {
	data     []byte
	fileName string
}

var (
	conf           v1beta3.PostgresConfig
	storage        StorageFn
	file           *File
	kubeConfigPath = homedir.HomeDir() + "/.kube/config"
)

const (
	currentContext    = "k3d-formance"
	configurationName = "stacks"
)
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randomString(length int) string {
	return stringWithCharset(length, charset)
}

func getKubeConfig(kubeconfigPath string, overrides *clientcmd.ConfigOverrides) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		overrides,
	).ClientConfig()
}

func deleteStack(cli *operatorclient.Client, ctx context.Context, stackName string) error {
	return cli.Stacks().Delete(ctx, stackName)
}

func newClient(config *rest.Config) (*operatorclient.Client, error) {
	sharedlogging.FromContext(context.Background()).Infof("Connect to cluster")
	defer func() {
		sharedlogging.FromContext(context.Background()).Infof("Connect to cluster OK")
	}()
	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &v1beta3.GroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	if crdConfig.UserAgent == "" {
		crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	client, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}

	return operatorclient.NewClient(client), nil
}

func countServiceDbs(conf v1beta3.ConfigurationServicesSpec) int {
	values := reflect.ValueOf(conf)
	count := 0
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			if servicesValues.Type().Field(j).Name != "Postgres" {
				continue
			}
			count++

		}
	}
	return count
}

func findDeployment(items []v1.Deployment, serviceName string) *v1.Deployment {
	for _, i := range items {
		if strings.ToLower(i.Name) == serviceName {
			return &i
		}
	}
	return nil
}

func waitStackUp(cli *kubernetes.Clientset, out io.Writer, ctx context.Context, conf v1beta3.ConfigurationServicesSpec, stackName string) error {
	tries := 0
	countServiceDbs := countServiceDbs(conf)
	for {
		deploymentStacks, err := cli.AppsV1().Deployments(stackName).List(ctx, v1meta.ListOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(out, "Waiting for stack to be fully up")

		if len(deploymentStacks.Items) > 0 {
			countServiceDbsAvailable := 0

			values := reflect.ValueOf(conf)
			for i := 0; i < values.NumField(); i++ {
				servicesValues := reflect.ValueOf(values.Field(i).Interface())
				for j := 0; j < servicesValues.NumField(); j++ {

					if servicesValues.Type().Field(j).Name != "Postgres" {
						continue
					}

					serviceName := strings.ToLower(values.Type().Field(i).Name)
					deploy := findDeployment(deploymentStacks.Items, serviceName)

					if deploy == nil {
						continue
					}

					for _, status := range deploy.Status.Conditions {
						if status.Type == v1.DeploymentAvailable {
							countServiceDbsAvailable++
						}
					}

				}
			}

			fmt.Println("minAvailableCount", countServiceDbsAvailable)
			fmt.Println("LEN", countServiceDbs)
			if countServiceDbsAvailable == countServiceDbs {
				return nil
			}
		}

		time.Sleep(1 * time.Second)
		tries++

		if tries > 60 {
			return fmt.Errorf("timeout")
		}
	}
}

func waitStackDisable(cli *kubernetes.Clientset, out io.Writer, ctx context.Context, stackName string) error {
	count := 0
	for {
		stack, err := cli.AppsV1().Deployments(stackName).List(ctx, v1meta.ListOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(out, "Waiting for stack to be fully disabled")

		if len(stack.Items) == 0 {
			return nil
		}

		time.Sleep(1 * time.Second)
		count++

		if count > 60 {
			return fmt.Errorf("timeout")
		}
	}
}

func TestBackupPostgres(t *testing.T) {

	// Get Kubeconfig
	override := &clientcmd.ConfigOverrides{
		CurrentContext: currentContext,
	}

	config, err := getKubeConfig(kubeConfigPath, override)
	if err != nil {
		t.Error(err.Error())
	}

	// Get the operator client
	cli, err := newClient(config)
	if err != nil {
		t.Error(err.Error())
	}

	// Create a stack
	randName := strings.ToLower(randomString(13))
	stackSpec := v1beta3.NewStack(randName, v1beta3.StackSpec{
		Seed:     configurationName,
		Host:     "host.k3d.internal",
		Scheme:   "http",
		Versions: "default",
		DevProperties: v1beta3.DevProperties{
			Debug: true,
			Dev:   true,
		},
		Auth: v1beta3.StackAuthSpec{
			DelegatedOIDCServer: v1beta3.DelegatedOIDCServerConfiguration{
				Issuer:       "http://host.k3d.internal/api/dex",
				ClientID:     "dexclient",
				ClientSecret: "dexclient",
			},
			StaticClients: []v1beta3.StaticClient{
				v1beta3.StaticClient{
					ID:      "foo2",
					Secrets: []string{"bar2"},
				},
			},
		},
	})
	newStack, err := cli.Stacks().Create(context.TODO(), &stackSpec)
	if err != nil {
		t.Errorf("Cannot create the stack: %s", err.Error())
	}

	// Get Kube client & Wait stack up
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Error("Error while getting kube client")
	}

	// Get the configuration
	conf, err := cli.Configurations().Get(context.TODO(), configurationName, v1meta.GetOptions{
		TypeMeta: v1meta.TypeMeta{
			Kind:       "Configuration",
			APIVersion: "stack.formancehq.com/v1beta3",
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	err = waitStackUp(clientSet, os.Stdout, context.TODO(), conf.Spec.Services, newStack.Name)
	if err != nil {
		t.Error(err.Error())
	}
	// Mock the destination storage
	storage = func(fileName string, data []byte) error {
		file = &File{
			fileName: fileName,
			data:     data,
		}

		return nil
	}

	logger := log.FromContext(context.TODO(), "stack", stackSpec.Name)

	values := reflect.ValueOf(conf.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			if servicesValues.Type().Field(j).Name != "Postgres" {
				continue
			}

			postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
			if !ok {
				t.Errorf("%s", "cannot cas to postgresconfig")
			}
			serviceName := strings.ToLower(values.Type().Field(i).Name)

			databaseName := fmt.Sprintf("%s-%s", newStack.Name, serviceName)

			t.Log(databaseName, postgresConfig, storage)

			err := backupPostgres(databaseName, postgresConfig, storage, logger)
			if err != nil {
				t.Error(err)
			}

			if !strings.HasPrefix(file.fileName, databaseName) {
				t.Errorf("backup filename should begin with %s", databaseName)
			}

			if len(file.data) == 0 {
				t.Errorf("In memory file should not be empty")
			}

		}
	}

	err = deleteStack(cli, context.TODO(), newStack.Name)
	if err != nil {
		t.Log("Stack not deleted")
	}
}

func TestDeletePostgresDatabases(t *testing.T) {

	// Get Kubeconfig
	override := &clientcmd.ConfigOverrides{
		CurrentContext: currentContext,
	}

	config, err := getKubeConfig(kubeConfigPath, override)
	if err != nil {
		t.Error(err.Error())
	}

	// Get the operator client
	cli, err := newClient(config)
	if err != nil {
		t.Error(err.Error())
	}

	// Create a stack
	randName := strings.ToLower(randomString(13))
	stackSpec := v1beta3.NewStack(randName, v1beta3.StackSpec{
		Seed:     configurationName,
		Host:     "host.k3d.internal",
		Scheme:   "http",
		Versions: "default",
		DevProperties: v1beta3.DevProperties{
			Debug: true,
			Dev:   true,
		},
		Auth: v1beta3.StackAuthSpec{
			DelegatedOIDCServer: v1beta3.DelegatedOIDCServerConfiguration{
				Issuer:       "http://host.k3d.internal/api/dex",
				ClientID:     "dexclient",
				ClientSecret: "dexclient",
			},
			StaticClients: []v1beta3.StaticClient{
				v1beta3.StaticClient{
					ID:      "foo2",
					Secrets: []string{"bar2"},
				},
			},
		},
	})

	newStack, err := cli.Stacks().Create(context.TODO(), &stackSpec)
	if err != nil {
		t.Errorf("Cannot create the stack: %s", err.Error())
	}

	// Get Kube client & Wait stack up
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Error("Error while getting kube client")
	}

	// Get the configuration
	conf, err := cli.Configurations().Get(context.TODO(), configurationName, v1meta.GetOptions{
		TypeMeta: v1meta.TypeMeta{
			Kind:       "Configuration",
			APIVersion: "stack.formancehq.com/v1beta3",
		},
	})
	if err != nil {
		t.Error(err.Error())
	}

	// Wait stack fully up and initialized
	if err = waitStackUp(clientSet, os.Stdout, context.TODO(), conf.Spec.Services, newStack.Name); err != nil {
		t.Error(err)
	}

	//Getting the stack again because it has been modified by the Reconcilier
	newStack, err = cli.Stacks().Get(context.TODO(), newStack.Name, v1meta.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	//Disable the stack
	newStack.Spec.Disabled = true
	newStack, err = cli.Stacks().Update(context.TODO(), newStack)
	if err != nil {
		t.Error(err)
	}

	// Wait all deployents deleted
	// So nobody is still accessing it
	err = waitStackDisable(clientSet, os.Stdout, context.TODO(), newStack.Name)
	if err != nil {
		t.Error(err.Error())
	}

	values := reflect.ValueOf(conf.Spec.Services)
	for i := 0; i < values.NumField(); i++ {
		servicesValues := reflect.ValueOf(values.Field(i).Interface())
		for j := 0; j < servicesValues.NumField(); j++ {
			if servicesValues.Type().Field(j).Name != "Postgres" {
				continue
			}

			postgresConfig, ok := servicesValues.Field(j).Interface().(v1beta3.PostgresConfig)
			if !ok {
				t.Errorf("%s", "cannot cas to postgresconfig")
			}
			serviceName := strings.ToLower(values.Type().Field(i).Name)

			databaseName := fmt.Sprintf("%s-%s", newStack.Name, serviceName)

			client, err := pg.OpenClient(postgresConfig)
			if err != nil {
				t.Error(err)
			}

			defer client.Close()
			if err := pg.DropDB(client, newStack.Name, serviceName); err != nil {
				t.Error(err)
			}
			t.Logf("Database droped %s", databaseName)
		}
	}
	err = deleteStack(cli, context.TODO(), newStack.Name)
	if err != nil {
		t.Log("Stack not deleted")
	}
}
