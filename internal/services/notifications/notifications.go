package notifications

import (
	"context"

	pb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"
)

const (
	eventNamePasswordRecovery    = "PasswordRecovery"
	eventNameProfileCreate       = "ProfileCreate"
	eventNameChangePassword      = "ChangePassword"
	eventNamePhoneVerification   = "PhoneVerification"
	eventNameEmailVerification   = "EmailVerification"
	eventNameFailedLoginAttempts = "FailedLoginAttempts"
	eventNameInviteCreate        = "InviteCreate"
)

type Notifications struct {
	clientFactory ClientFactory
}

func NewNotifications(clientFactory ClientFactory) *Notifications {
	return &Notifications{clientFactory}
}

// PasswordRecovery sends a confirmation code to recover a password
func (s *Notifications) PasswordRecovery(userID, confirmationCode string, methods []string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNamePasswordRecovery,
		TemplateData: &pb.TemplateData{
			ConfirmationCode: confirmationCode,
		},
		Notifiers: methods,
	})
}

// ProfileCreated send a notification after an profile was created
func (s *Notifications) ProfileCreated(userID, password, confirmationCode string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNameProfileCreate,
		TemplateData: &pb.TemplateData{
			Password:                    password,
			SetPasswordConfirmationCode: confirmationCode,
		},
	})
}

// PasswordChanged sends a notification when a user's password was changed
func (s *Notifications) PasswordChanged(userID string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNameChangePassword,
	})
}

// VerifyPhone sends a confirmation code to the user to verify his phone number
func (s *Notifications) VerifyPhone(userID, confirmationCode string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNamePhoneVerification,
		TemplateData: &pb.TemplateData{
			ConfirmationCode: confirmationCode,
		},
	})
}

// VerifyEmail sends a confirmation code to the user to verify his email address
func (s *Notifications) VerifyEmail(userID, confirmationCode string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNameEmailVerification,
		TemplateData: &pb.TemplateData{
			ConfirmationCode: confirmationCode,
		},
	})
}

// FailLoginAttempts sends a notification when a user exceeded failed login attempts
func (s *Notifications) FailLoginAttempts(userID string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNameFailedLoginAttempts,
	})
}

// InviteCreated sends a notification when an invite was created
func (s *Notifications) InviteCreated(userID, password string) (*pb.Response, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	return client.Dispatch(context.Background(), &pb.Request{
		To:        userID,
		EventName: eventNameInviteCreate,
		TemplateData: &pb.TemplateData{
			Password: password,
		},
	})
}
