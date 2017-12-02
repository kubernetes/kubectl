package test_test

import (
	"fmt"

	. "k8s.io/kubectl/pkg/framework/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TempDirManager", func() {
	var (
		manager            *TempDirManager
		removerError       error
		createError        error
		managedDirCount    int
		separateDirCounter int
	)
	BeforeEach(func() {
		managedDirCount = 0
		separateDirCounter = 0
		createError = nil
		removerError = nil
		manager = NewTempDirManager()
		manager.Maker = func(dir, prefix string) (string, error) {
			managedDirCount += 1
			separateDirCounter += 1
			return fmt.Sprintf("%d-%s-%s", separateDirCounter, dir, prefix), createError
		}
		manager.Remover = func(dir string) error {
			managedDirCount -= 1
			return removerError
		}
	})

	It("can creates and remove directories", func() {
		Expect(managedDirCount).To(Equal(0))
		manager.Create()
		Expect(managedDirCount).To(Equal(1))
		manager.Destroy()
		Expect(managedDirCount).To(Equal(0))
	})

	Context("when I call Create() multiple times on the same manager", func() {
		It("returns the same directory every time", func() {
			var dir1, dir2 string
			var err error

			Expect(managedDirCount).To(Equal(0))

			dir1, err = manager.Create()
			Expect(err).NotTo(HaveOccurred())
			Expect(managedDirCount).To(Equal(1))

			dir2, err = manager.Create()
			Expect(err).NotTo(HaveOccurred())
			Expect(managedDirCount).To(Equal(1))
			Expect(dir1).To(Equal(dir2))
		})

		It("deletes the managed directory as soon as Destroy() is called even once", func() {
			var err error

			Expect(managedDirCount).To(Equal(0))

			_, err = manager.Create()
			Expect(err).NotTo(HaveOccurred())
			_, err = manager.Create()
			Expect(err).NotTo(HaveOccurred())
			Expect(managedDirCount).To(Equal(1))

			manager.Destroy()
			Expect(managedDirCount).To(Equal(0))
		})
	})

	Context("when I call Destroy() without calling create first", func() {
		It("does nothing", func() {
			Expect(managedDirCount).To(Equal(0))
			manager.Destroy()
			Expect(managedDirCount).To(Equal(0))
		})
	})

	Context("when the remover returns an error", func() {
		JustBeforeEach(func() {
			removerError = fmt.Errorf("Error on removing dir")
		})
		It("handles that error depending on whether Create() has been called", func() {
			By("avoiding the error if Create() has not been called")
			err := manager.Destroy()
			Expect(err).NotTo(HaveOccurred())

			By("propagating the error if Create() has been called")
			manager.Create()
			err = manager.Destroy()
			Expect(err).To(MatchError("Error on removing dir"))
		})
	})

	Context("when the creater returns an error", func() {
		JustBeforeEach(func() {
			createError = fmt.Errorf("Error on creating dir")
		})
		It("bubbles up the error", func() {
			_, err := manager.Create()
			Expect(err).To(MatchError("Error on creating dir"))
		})
	})
})
