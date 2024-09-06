//go:build it

package test_suite

import (
	"context"
	"fmt"
	ingester "github.com/formancehq/stack/ee/ingester/internal"
	. "github.com/formancehq/stack/ee/ingester/pkg/testserver"
	. "github.com/formancehq/stack/libs/go-libs/testing/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

func sendBatch(service *FakeModule, offset, size int) {
	GinkgoHelper()
	for i := 0; i < size; i++ {
		go func() {
			service.PushLogs(GinkgoT(), ingester.Log{
				ID: fmt.Sprintf("%d", offset+i),
			})
		}()
	}
}

func shouldHaveInsert(connector *Deferred[Connector], n int) {
	GinkgoHelper()
	Eventually(connector.GetValue().ReadMessages).
		WithContext(context.Background()).
		WithTimeout(10 * time.Second).
		Should(HaveLen(n))
}

func shouldHaveInsertConsistently(connector *Deferred[Connector], n int) func() {
	return func() {
		GinkgoHelper()
		Consistently(connector.GetValue().ReadMessages).
			WithArguments(context.Background()).
			Should(HaveLen(n))
	}
}
