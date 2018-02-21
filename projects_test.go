package mcapi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/materials-commons/gomcapi"
)

var _ = Describe("Project", func() {
	var createdProjectID string

	Describe("GetAllProjects", func() {
		It("Should get all projects for user", func() {
			projects, err := GetAllProjects()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(projects).ShouldNot(BeNil())
		})
	})

	Describe("CreateProject", func() {
		It("Should create a project", func() {
			projName := "Proj1"
			projDescription := "Project Created With Test"
			proj, err := CreateProject(projName, projDescription)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(proj).ShouldNot(BeNil())
			Expect(proj.Name).Should(Equal(projName))
			Expect(proj.Description).Should(Equal(projDescription))

			createdProjectID = proj.ID
		})

		It("Should delete the created project", func() {
			err := DeleteProject(createdProjectID)
			Expect(err).Should(BeNil())
		})
	})
})
