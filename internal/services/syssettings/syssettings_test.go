package syssettings

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Confialink/wallet-users/internal/services/syssettings/mocks"
	settingsMock "github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/rpc/settings"
)

var _ = Describe("syssettings package", func() {
	var (
		clientFactory *mocks.ClientFactory
		client        *settingsMock.MockSettingsHandler
	)
	BeforeEach(func() {
		clientFactory = &mocks.ClientFactory{}
		client = &settingsMock.MockSettingsHandler{}
	})

	Context("GetTimeSettings", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetTimeSettings()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: "regional/general/%"}
					client.On("List", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetTimeSettings()
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			When("everything is ok", func() {
				It("should not return an error", func() {
					req := &pb.Request{Path: "regional/general/%"}
					timeZone := "Europe/Minsk"
					dateFormat := "DD/MM/YYYY"
					timeFormat := "hh:mm A"

					resp := &pb.Response{
						Settings: []*pb.Setting{
							{Path: defaultTimeZonePath, Value: timeZone},
							{Path: defaultDateFormatPath, Value: dateFormat},
							{Path: defaultTimeFormatPath, Value: timeFormat},
						},
					}
					client.On("List", context.Background(), req).Return(resp, nil)
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetTimeSettings()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(res.Timezone).Should(Equal(timeZone))
					Expect(res.TimeFormat).Should(Equal(timeFormat))
					Expect(res.DateFormat).Should(Equal(dateFormat))
					Expect(res.DateTimeFormat).Should(Equal(fmt.Sprintf("%s %s", dateFormat, timeFormat)))
				})
			})
		})
	})

	Context("GetAutologoutSettings", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetAutologoutSettings()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: "profile/autologout/%"}
					client.On("List", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetAutologoutSettings()
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			Context("the HTTP client returns response", func() {
				When("AutoLogoutTimeout has an invalid value", func() {
					It("should return an error", func() {
						req := &pb.Request{Path: "profile/autologout/%"}
						resp := &pb.Response{
							Settings: []*pb.Setting{
								{Path: autoLogoutStatusPath, Value: "yes"},
								{Path: autoLogoutTimeoutPath, Value: "not-integer-value"}, // wrong value
								{Path: autoLogoutPaddingPath, Value: "20"},
							},
						}
						client.On("List", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetAutologoutSettings()
						Expect(err).Should(HaveOccurred())
						Expect(res).Should(BeNil())
					})
				})
				When("AutoLogoutPadding has an invalid value", func() {
					It("should return an error", func() {
						req := &pb.Request{Path: "profile/autologout/%"}
						resp := &pb.Response{
							Settings: []*pb.Setting{
								{Path: autoLogoutStatusPath, Value: "yes"},
								{Path: autoLogoutTimeoutPath, Value: "20"},
								{Path: autoLogoutPaddingPath, Value: "not-integer-value"}, // wrong value
							},
						}
						client.On("List", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetAutologoutSettings()
						Expect(err).Should(HaveOccurred())
						Expect(res).Should(BeNil())
					})
				})

				When("everything is ok", func() {
					It("should not return an error", func() {
						req := &pb.Request{Path: "profile/autologout/%"}
						resp := &pb.Response{
							Settings: []*pb.Setting{
								{Path: autoLogoutStatusPath, Value: "yes"},
								{Path: autoLogoutTimeoutPath, Value: "20"},
								{Path: autoLogoutPaddingPath, Value: "21"},
							},
						}
						client.On("List", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetAutologoutSettings()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res.Enabled).Should(BeTrue())

						Expect(res.Timeout).Should(Equal(time.Duration(time.Duration(20) * time.Minute)))
						Expect(res.Padding).Should(Equal(time.Duration(time.Duration(21) * time.Second)))
					})
				})
			})
		})
	})

	Context("GetGDPRSettings", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetGDPRSettings()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: gdprSettingPath}
					client.On("Get", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetGDPRSettings()
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			When("everything is ok", func() {
				It("should not return an error", func() {
					req := &pb.Request{Path: gdprSettingPath}
					resp := &pb.Response{
						Setting: &pb.Setting{
							Path: gdprSettingPath, Value: "enable",
						},
					}
					client.On("Get", context.Background(), req).Return(resp, nil)
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetGDPRSettings()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(res.Enabled).Should(BeTrue())
				})
			})
		})
	})

	Context("GetLoginSecuritySettings", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetLoginSecuritySettings()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: "regional/login/%"}
					client.On("List", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					_, err := service.GetLoginSecuritySettings()
					Expect(err).Should(HaveOccurred())
				})
			})

			When("everything is ok", func() {
				It("should not return an error", func() {
					req := &pb.Request{Path: "regional/login/%"}
					resp := &pb.Response{
						Settings: []*pb.Setting{
							{Path: loginUsernameCleanupPath, Value: "5"},
							{Path: loginUsernameLimitPath, Value: "6"},
							{Path: loginUsernameUsePath, Value: "yes"},
							{Path: loginUserUsePath, Value: "100"},
							{Path: loginUserWindowPath, Value: "10"},
						},
					}
					client.On("List", context.Background(), req).Return(resp, nil)
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetLoginSecuritySettings()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(res.FailedLoginUsernameCleanup).Should(Equal(uint64(5)))
					Expect(res.FailedLoginUsernameLimit).Should(Equal(uint64(6)))
					Expect(res.FailedLoginUsernameUse).Should(BeTrue())
					Expect(res.FailedLoginUserUse).Should(Equal(uint64(100)))
					Expect(res.FailedLoginUserWindow).Should(Equal(uint64(10)))
				})
			})
		})
	})

	Context("GetDormantDuration", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetDormantDuration()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(Equal(time.Nanosecond))
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: userOptionsDormantPath}
					client.On("Get", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetDormantDuration()
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(Equal(time.Nanosecond))
				})
			})

			Context("the HTTP client returns response", func() {
				When("setting value is invalid", func() {
					It("should return an error", func() {
						req := &pb.Request{Path: userOptionsDormantPath}
						resp := &pb.Response{
							Setting: &pb.Setting{
								Path: userOptionsDormantPath, Value: "not-integer-value", // wrong value
							},
						}
						client.On("Get", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetDormantDuration()
						Expect(err).Should(HaveOccurred())
						Expect(res).Should(Equal(time.Nanosecond))
					})
				})

				When("everything is ok", func() {
					It("should not return an error", func() {
						req := &pb.Request{Path: userOptionsDormantPath}
						resp := &pb.Response{
							Setting: &pb.Setting{
								Path: userOptionsDormantPath, Value: "5",
							},
						}
						client.On("Get", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetDormantDuration()
						Expect(err).ShouldNot(HaveOccurred())

						hoursCount := 5 * 30 * 24
						r, _ := time.ParseDuration(fmt.Sprintf("%dh", hoursCount))
						Expect(res).Should(Equal(r))
					})
				})
			})
		})
	})

	Context("GetMaintenanceModeSettings", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetMaintenanceModeSettings()
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: maintenancePath}
					client.On("Get", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetMaintenanceModeSettings()
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			Context("the HTTP client returns response", func() {
				When("everything is ok", func() {
					It("should not return an error", func() {
						req := &pb.Request{Path: maintenancePath}
						resp := &pb.Response{
							Setting: &pb.Setting{
								Path: maintenancePath, Value: "enable",
							},
						}
						client.On("Get", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetMaintenanceModeSettings()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res.Enabled).Should(BeTrue())
					})
				})
			})
		})
	})

	Context("GetDefaultUserClassByRole", func() {
		roleName := "client" // random value
		settingPath := fmt.Sprintf("%s/%s", "profile/default-user-classes", roleName)
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewSysSettings(clientFactory)

				res, err := service.GetDefaultUserClassByRole(roleName)
				Expect(err).Should(HaveOccurred())
				Expect(res).Should(BeNil())
			})
		})

		Context("client factory returns an HTTP client", func() {
			When("Settings service returns an error", func() {
				It("should return an error", func() {
					req := &pb.Request{Path: settingPath}
					client.On("Get", context.Background(), req).Return(nil, errors.New("random err"))
					clientFactory.On("NewClient").Return(client, nil)
					service := NewSysSettings(clientFactory)
					res, err := service.GetDefaultUserClassByRole(roleName)
					Expect(err).Should(HaveOccurred())
					Expect(res).Should(BeNil())
				})
			})

			Context("the HTTP client returns response", func() {
				When("everything is ok", func() {
					It("should not return an error", func() {
						result := "12"
						req := &pb.Request{Path: settingPath}
						resp := &pb.Response{
							Setting: &pb.Setting{
								Path: settingPath, Value: result,
							},
						}
						client.On("Get", context.Background(), req).Return(resp, nil)
						clientFactory.On("NewClient").Return(client, nil)
						service := NewSysSettings(clientFactory)
						res, err := service.GetDefaultUserClassByRole(roleName)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(res).Should(Equal(&result))
					})
				})
			})
		})
	})
})
