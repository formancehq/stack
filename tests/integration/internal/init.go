package internal

import (
	"context"
	formance "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/stack/libs/go-libs/logging"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	ctx         context.Context
	currentEnv  *env
	currentTest *Test
)

func TestContext() context.Context {
	return ctx
}

func Client() *formance.Formance {
	return currentTest.sdkClient
}

var _ = BeforeSuite(func() {
	ctx = context.TODO()

	l := logrus.New()
	l.Out = GinkgoWriter
	l.Level = logrus.DebugLevel
	ctx = logging.ContextWithLogger(ctx, logging.NewLogrus(l))

	// Some defaults
	SetDefaultEventuallyTimeout(10 * time.Second)
	SetDefaultEventuallyPollingInterval(500 * time.Millisecond)

	Eventually(func() error {
		_, err := http.Get("http://" + GetOpenSearchUrl())
		return err
	}).Should(Succeed())

	currentEnv = newEnv(GinkgoWriter)
	Expect(currentEnv.Setup(ctx)).To(Succeed())
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
