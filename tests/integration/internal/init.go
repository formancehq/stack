package internal

import (
	"context"
	formance "github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/go-libs/logging"
	"github.com/oauth2-proxy/mockoidc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var (
	ctx         context.Context
	currentEnv  *Env
	currentTest *Test
	mockOIDC    *mockoidc.MockOIDC
)

func CurrentEnv() *Env {
	return currentEnv
}

func CurrentTest() *Test {
	return currentTest
}

func TestContext() context.Context {
	return ctx
}

func Client(options ...formance.SDKOption) *formance.Formance {
	gatewayUrl, err := url.Parse(currentTest.httpServer.URL)
	if err != nil {
		panic(err)
	}

	Expect(err).To(BeNil())

	options = append([]formance.SDKOption{
		formance.WithServerURL(gatewayUrl.String()),
		formance.WithClient(
			&http.Client{
				Transport: currentTest.httpTransport,
			},
		),
	}, options...)

	return formance.New(options...)
}

func GatewayURL() string {
	return currentTest.GatewayURL()
}

func HTTPClient() *http.Client {
	return &http.Client{
		Transport: currentTest.httpTransport,
	}
}

func OIDCServer() *mockoidc.MockOIDC {
	return mockOIDC
}

var _ = BeforeSuite(func() {
	ctx = context.TODO()

	l := logrus.New()
	l.Out = GinkgoWriter
	l.Level = logrus.DebugLevel
	ctx = logging.ContextWithLogger(ctx, logging.NewLogrus(l))

	var err error
	mockOIDC, err = mockoidc.Run()
	Expect(err).To(BeNil())

	// Some defaults
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(500 * time.Millisecond)

	Eventually(func() error {
		_, err := http.Get("http://" + GetOpenSearchUrl())
		return err
	}).Should(Succeed())

	currentEnv = newEnv(GinkgoWriter)
	Expect(currentEnv.Setup(ctx)).To(Succeed())

	Eventually(func() error {
		_, err := http.Get("http://" + GetOpenSearchUrl())
		return err
	}).Should(Succeed())
})

var _ = AfterSuite(func() {
	Expect(currentEnv.Teardown(ctx)).To(Succeed())
})

var _ = BeforeEach(func() {
	currentTest = currentEnv.newTest()
	Expect(currentTest.setup()).To(Succeed())
})

var _ = AfterEach(func() {
	Expect(currentTest.tearDown()).To(Succeed())
})
