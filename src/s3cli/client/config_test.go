package client

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe(".newConfigFromPath", func() {
		It("has sensible defaults for all the properties", func() {
			configPath := writeConfigFile(`{}`)

			defer os.Remove(configPath)

			config, err := newConfigFromPath(configPath)
			Expect(err).ToNot(HaveOccurred())
			Expect(config.AccessKeyID).To(Equal(""))
			Expect(config.SecretAccessKey).To(Equal(""))
			Expect(config.BucketName).To(Equal(""))
			Expect(config.CredentialsSource).To(Equal("static"))
			Expect(config.Host).To(Equal(""))
			Expect(config.Port).To(Equal(443))
			Expect(config.Region).To(Equal("us-east-1"))
			Expect(config.SSLVerifyPeer).To(BeTrue())
			Expect(config.UseSSL).To(BeTrue())
			Expect(config.UseSSL).To(BeTrue())
		})

		Describe("#AccessKeyID", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"access_key_id": "secret"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.AccessKeyID).To(Equal("secret"))
			})
		})

		Describe("#SecretAccessKey", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"secret_access_key": "secret"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.SecretAccessKey).To(Equal("secret"))
			})
		})

		Describe("#BucketName", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"bucket_name": "foo"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.BucketName).To(Equal("foo"))
			})
		})

		Describe("#CredentialsSource", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"credentials_source": "any-source"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.CredentialsSource).To(Equal("any-source"))
			})

			It("will convert a blank value to staic", func() {
				configPath := writeConfigFile(`{
					"credentials_source": ""
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.CredentialsSource).To(Equal("static"))
			})
		})

		Describe("#Host", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"host": "host.example.com"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Host).To(Equal("host.example.com"))
			})
		})

		Describe("#Port", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"port": 123
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Port).To(Equal(123))
			})
		})

		Describe("#Region", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"region": "ca-central-2"
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Region).To(Equal("ca-central-2"))
			})
		})

		Describe("#SSLVerifyPeer", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"ssl_verify_peer": false
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.SSLVerifyPeer).To(BeFalse())
			})
		})

		Describe("#UseSSL", func() {
			It("can be overridden", func() {
				configPath := writeConfigFile(`{
					"use_ssl": false
				}`)

				defer os.Remove(configPath)

				config, err := newConfigFromPath(configPath)
				Expect(err).ToNot(HaveOccurred())
				Expect(config.UseSSL).To(BeFalse())
			})
		})

		Describe("Error handling", func() {
			Context("when config file is missing", func() {
				It("returns an error", func() {
					config, err := newConfigFromPath("/non/existant/file")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("no such file or directory"))
					Expect(config).To(BeNil())
				})
			})

			Context("when the config file has invalid json", func() {
				It("returns an error when the type of a key is invalid", func() {
					configPath := writeConfigFile(`{"use_ssl": "Yes"}`)

					defer os.Remove(configPath)

					config, err := newConfigFromPath(configPath)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("json: cannot unmarshal string"))
					Expect(config).To(BeNil())
				})

				It("returns an error when the file is not correct json", func() {
					configPath := writeConfigFile(`{"UseSSL: true`)

					defer os.Remove(configPath)

					config, err := newConfigFromPath(configPath)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(MatchRegexp("unexpected end of JSON input"))
					Expect(config).To(BeNil())
				})
			})
		})
	})
})

func writeConfigFile(contents string) string {
	file, err := ioutil.TempFile("", "client_test")
	Expect(err).ToNot(HaveOccurred())

	err = file.Close()
	Expect(err).ToNot(HaveOccurred())

	err = ioutil.WriteFile(file.Name(), []byte(contents), os.ModeTemporary)
	Expect(err).ToNot(HaveOccurred())

	return file.Name()
}
