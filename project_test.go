package mcapi_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/materials-commons/gomcapi"
)

var _ = Describe("Project", func() {
	It("Should get all projects", func() {
		projects, err := GetAllProjects()
		Expect(err).ShouldNot(HaveOccurred())
		Expect(projects).ShouldNot(BeNil())
		for _, proj := range projects {
			fmt.Printf("%+v\n", proj)
		}
	})
})
