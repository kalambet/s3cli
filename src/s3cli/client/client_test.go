package client

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlobstoreClient", func() {
	Describe(".New", func() {
		Context("with SSL cert verification turned off", func() {
			It("returns client that uses a HTTP Client with SSL cert verification off", func() {
				config := Config{
					AccessKeyID:       "fake-access-key",
					SecretAccessKey:   "fake-secret-key",
					CredentialsSource: "static",
					SSLVerifyPeer:     false,
				}

				blobstoreClient, err := New(config)
				Expect(err).ToNot(HaveOccurred())
				roundTripper := blobstoreClient.s3Client.Config.HTTPClient.Transport
				Expect(roundTripper).ToNot(BeNil())
				transport := roundTripper.(*http.Transport)
				Expect(transport.TLSClientConfig).ToNot(BeNil())
				Expect(transport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
				Expect(transport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())

				Expect(http.DefaultTransport.(*http.Transport).TLSClientConfig).To(BeNil())
			})
		})

		It("configures the region", func() {
			config := Config{
				AccessKeyID:       "fake-access-key",
				SecretAccessKey:   "fake-secret-key",
				CredentialsSource: "static",
				Region:            "my-custom-region",
			}

			blobstoreClient, err := New(config)
			Expect(err).ToNot(HaveOccurred())
			s3Config := blobstoreClient.s3Client.Config
			Expect(*s3Config.Region).To(Equal("my-custom-region"))
		})

		It("uses the 'path' style access to bucket", func() {
			config := Config{
				AccessKeyID:       "fake-access-key",
				SecretAccessKey:   "fake-secret-key",
				CredentialsSource: "static",
			}

			blobstoreClient, err := New(config)
			Expect(err).ToNot(HaveOccurred())
			s3Config := blobstoreClient.s3Client.Config
			Expect(s3Config.S3ForcePathStyle).ToNot(BeNil())
			Expect(*s3Config.S3ForcePathStyle).To(BeTrue())
		})

		Describe("config.CredentialsSource", func() {
			var config Config
			Context("with CredentialsSource set to static", func() {
				BeforeEach(func() {
					config = Config{
						CredentialsSource: "static",
						AccessKeyID:       "fake-access-key",
						SecretAccessKey:   "fake-secret-key",
					}
				})

				It("returns properly configured client", func() {
					client, err := New(config)
					Expect(err).ToNot(HaveOccurred())

					credentials, err := client.s3Client.Config.Credentials.Get()
					Expect(err).ToNot(HaveOccurred())

					Expect(credentials.AccessKeyID).To(Equal("fake-access-key"))
					Expect(credentials.SecretAccessKey).To(Equal("fake-secret-key"))
				})

				It("raises an error when AccessKeyID is not provided", func() {
					config.AccessKeyID = ""

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeyMissing))
				})

				It("raises an error when SecretAccessKey is not provided", func() {
					config.SecretAccessKey = ""

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeyMissing))
				})

				It("raises an error when both AccessKeyID and SecretAccessKey are not provided", func() {
					config.SecretAccessKey = ""
					config.AccessKeyID = ""

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeyMissing))
				})
			})

			Context("with CredentialsSource set to env_or_profile", func() {
				BeforeEach(func() {
					config = Config{
						CredentialsSource: "env_or_profile",
					}
				})

				It("returns properly configured client", func() {
					_, err := New(config)
					Expect(err).ToNot(HaveOccurred())
				})

				It("raises an error when AccessKeyID is also provided", func() {
					config.AccessKeyID = "fake-access-key"

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeysPresent))
				})

				It("raises an error when SecretAccessKey is also provided", func() {
					config.SecretAccessKey = "fake-secret-key"

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeysPresent))
				})

				It("raises an error when AccessKeyID and SecretAccessKey are also provided", func() {
					config.AccessKeyID = "fake-access-key"
					config.SecretAccessKey = "fake-secret-key"

					_, err := New(config)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(errorKeysPresent))
				})
			})

			Context("with CredentialsSource set to incorrect value", func() {
				It("raises an error", func() {
					config := Config{
						CredentialsSource: "incorrect-value",
					}

					_, err := New(config)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
