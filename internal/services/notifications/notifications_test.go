package notifications

import (
	"context"
	"errors"

	pb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Confialink/wallet-users/internal/services/notifications/mocks"
	notificationsMock "github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/rpc/notifications"
)

var _ = Describe("notifications package", func() {
	var (
		clientFactory *mocks.ClientFactory
		client        *notificationsMock.NotificationHandler
	)
	BeforeEach(func() {
		clientFactory = &mocks.ClientFactory{}
		client = &notificationsMock.NotificationHandler{}
	})

	userID := "random-uid"

	Context("PasswordRecovery", func() {
		confirmationCode := "random-code"
		notifyMethods := []string{"sms"}

		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.PasswordRecovery(userID, confirmationCode, notifyMethods)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNamePasswordRecovery,
					TemplateData: &pb.TemplateData{
						ConfirmationCode: confirmationCode,
					},
					Notifiers: notifyMethods,
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.PasswordRecovery(userID, confirmationCode, notifyMethods)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("ProfileCreated", func() {
		confirmationCode := "random-code"
		password := "random password"

		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.ProfileCreated(userID, password, confirmationCode)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNameProfileCreate,
					TemplateData: &pb.TemplateData{
						SetPasswordConfirmationCode: confirmationCode,
						Password:                    password,
					},
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.ProfileCreated(userID, password, confirmationCode)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("PasswordChanged", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.PasswordChanged(userID)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNameChangePassword,
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.PasswordChanged(userID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("VerifyPhone", func() {
		confirmationCode := "random-code"
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.VerifyPhone(userID, confirmationCode)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNamePhoneVerification,
					TemplateData: &pb.TemplateData{
						ConfirmationCode: confirmationCode,
					},
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.VerifyPhone(userID, confirmationCode)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("VerifyEmail", func() {
		confirmationCode := "random-code"
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.VerifyEmail(userID, confirmationCode)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNameEmailVerification,
					TemplateData: &pb.TemplateData{
						ConfirmationCode: confirmationCode,
					},
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.VerifyEmail(userID, confirmationCode)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("FailLoginAttempts", func() {
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.FailLoginAttempts(userID)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNameFailedLoginAttempts,
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.FailLoginAttempts(userID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})

	Context("InviteCreated", func() {
		password := "random-code"
		When("client factory returns an error", func() {
			It("should return an error", func() {
				clientFactory.On("NewClient").Return(nil, errors.New("random err"))
				service := NewNotifications(clientFactory)

				_, err := service.InviteCreated(userID, password)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("notification is successfully sent", func() {
			It("should not return an error", func() {
				req := &pb.Request{
					To:        userID,
					EventName: eventNameInviteCreate,
					TemplateData: &pb.TemplateData{
						Password: password,
					},
				}
				resp := &pb.Response{}
				client.On("Dispatch", context.Background(), req).Return(resp, nil)
				clientFactory.On("NewClient").Return(client, nil)
				service := NewNotifications(clientFactory)
				res, err := service.InviteCreated(userID, password)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(res).Should(Equal(resp))
			})
		})
	})
})
