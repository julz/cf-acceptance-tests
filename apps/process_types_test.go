package apps

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers"
)

var _ = Describe("Process Types", func() {
	type AppResource struct {
		Entity struct {
			DetectedStartCommand string `json:"detected_start_command"`
		} `json:"entity"`
	}

	type AppsResponse struct {
		Resources []AppResource `json:"resources"`
	}

	Describe("Staging", func() {
		var appName string

		BeforeEach(func() {
			appName = generator.RandomName()
		})

		Describe("without a procfile", func() {
			BeforeEach(func() {
				Expect(cf.Cf("push", appName, "-p", helpers.NewAssets().Node, "-c", "node app.js").Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
			})

			AfterEach(func() {
				Expect(cf.Cf("delete", appName, "-f").Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
			})

			It("detects the use of the start command supplied on the command line", func() {
				var appsResponse AppsResponse
				cfResponse := cf.Cf("curl", fmt.Sprintf("/v2/apps?q=name:%s", appName)).Wait(DEFAULT_TIMEOUT).Out.Contents()
				json.Unmarshal(cfResponse, &appsResponse)

				Expect(appsResponse.Resources[0].Entity.DetectedStartCommand).To(Equal("node app.js"))
			})
		})

		Describe("with a procfile", func() {
			BeforeEach(func() {
				Expect(cf.Cf("push", appName, "-p", helpers.NewAssets().NodeWithProcfile).Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
			})

			AfterEach(func() {
				Expect(cf.Cf("delete", appName, "-f").Wait(CF_PUSH_TIMEOUT)).To(Exit(0))
			})

			It("detects the use of the start command in the 'web' process type", func() {
				var appsResponse AppsResponse
				cfResponse := cf.Cf("curl", fmt.Sprintf("/v2/apps?q=name:%s", appName)).Wait(DEFAULT_TIMEOUT).Out.Contents()
				json.Unmarshal(cfResponse, &appsResponse)

				Expect(appsResponse.Resources[0].Entity.DetectedStartCommand).To(Equal("node app.js"))
			})
		})
	})
})
