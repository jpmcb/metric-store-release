package rules_test

import (
	"crypto/tls"

	. "github.com/cloudfoundry/metric-store-release/src/internal/rules"
	"github.com/cloudfoundry/metric-store-release/src/internal/testing"
	shared "github.com/cloudfoundry/metric-store-release/src/internal/tls"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type remoteRuleManagerTestContext struct {
	rulesApiSpy     *testing.RulesApiSpy
	tlsClientConfig *tls.Config
}

var _ = Describe("RemoteRuleManager", func() {
	var setup = func() (*remoteRuleManagerTestContext, func()) {
		tlsServerConfig, err := shared.NewMutualTLSServerConfig(
			testing.Cert("metric-store-ca.crt"),
			testing.Cert("metric-store.crt"),
			testing.Cert("metric-store.key"),
		)
		Expect(err).ToNot(HaveOccurred())

		tlsClientConfig, err := shared.NewMutualTLSClientConfig(
			testing.Cert("metric-store-ca.crt"),
			testing.Cert("metric-store.crt"),
			testing.Cert("metric-store.key"),
			"metric-store",
		)
		Expect(err).ToNot(HaveOccurred())

		rulesApiSpy, err := testing.NewRulesApiSpy(tlsServerConfig)
		Expect(err).ToNot(HaveOccurred())
		rulesApiSpy.Start()

		tc := &remoteRuleManagerTestContext{
			rulesApiSpy:     rulesApiSpy,
			tlsClientConfig: tlsClientConfig,
		}

		return tc, func() {
			rulesApiSpy.Stop()
		}
	}

	FDescribe("#CreateManager", func() {
		It("returns nothing when the api does not return an error", func() {
			tc, cleanup := setup()
			defer cleanup()

			remoteRuleManager := NewRemoteRuleManager(tc.rulesApiSpy.Addr(), tc.tlsClientConfig)

			err := remoteRuleManager.CreateManager("app-metrics", "")

			Expect(err).NotTo(HaveOccurred())
		})

		It("returns ManagerExistsError when api returns ErrorNotCreated", func() {
			tc, cleanup := setup()
			defer cleanup()

			remoteRuleManager := NewRemoteRuleManager(tc.rulesApiSpy.Addr(), tc.tlsClientConfig)

			tc.rulesApiSpy.NextRequestError(&testing.RulesApiHttpError{
				Status: 409,
			})
			err := remoteRuleManager.CreateManager("app-metrics", "")

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ManagerExistsError))
		})

		It("returns api error by default", func() {
			tc, cleanup := setup()
			defer cleanup()

			remoteRuleManager := NewRemoteRuleManager(tc.rulesApiSpy.Addr(), tc.tlsClientConfig)

			tc.rulesApiSpy.NextRequestError(&testing.RulesApiHttpError{
				Status: 201,
			})
			err := remoteRuleManager.CreateManager("app-metrics", "")

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ManagerExistsError))
		})
	})
})
