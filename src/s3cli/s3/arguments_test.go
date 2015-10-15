package s3

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Argument handling", func() {
	Describe("validating arguments", func() {
		It("returns an error if the wrong number of arugments are supplied", func() {
			dummyArgs := []string{"get"}
			_, _, _, err := validateArguments(dummyArgs)
			Expect(err).To(HaveOccurred())
		})

		for _, cmd := range []string{"get", "put"} {
			It("returns the "+cmd+" command, the source, and the destination", func() {
				dummyArgs := []string{cmd, "dummyKey", "localPath"}
				command, src, dst, err := validateArguments(dummyArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(command).To(Equal(cmd))
				Expect(src).To(Equal("dummyKey"))
				Expect(dst).To(Equal("localPath"))
			})
		}

		It("returns an error if the command is not get or put", func() {
			dummyArgs := []string{"badCommand", "dummyKey", "localPath"}
			_, _, _, err := validateArguments(dummyArgs)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("loading configuration", func() {
		Context("when configuration path is provided", func() {
			It("returns a valid configuration path and the non-flag arguments", func() {
				dummyArgs := []string{"-c", "/some/path", "get"}
				path, nonFlagArgs, err := fetchConfigurationPath(dummyArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(Equal("/some/path"))
				Expect(nonFlagArgs).To(Equal([]string{"get"}))
			})
		})

		Context("when a configuration path is not provided", func() {
			It("returns a valid configuration path to $HOME/.s3cli and the non-flag arguments", func() {
				dummyArgs := []string{"get", "source", "destination"}
				path, nonFlagArgs, err := fetchConfigurationPath(dummyArgs)
				Expect(err).ToNot(HaveOccurred())
				Expect(path).To(ContainSubstring("/.s3cli"))
				Expect(nonFlagArgs).To(Equal([]string{"get", "source", "destination"}))
			})
		})
	})
})
