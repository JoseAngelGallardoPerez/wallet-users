package syssettings

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/Confialink/wallet-settings/rpc/proto/settings"
)

const (
	gdprSettingPath = "regional/modules/velmie_wallet_gdpr"

	defaultTimeZonePath   = "regional/general/default_timezone"
	defaultDateFormatPath = "regional/general/default_date_format"
	defaultTimeFormatPath = "regional/general/default_time_format"

	maintenancePath = "regional/general/maintenance"

	autoLogoutStatusPath  = "profile/autologout/status"
	autoLogoutTimeoutPath = "profile/autologout/timeout"
	autoLogoutPaddingPath = "profile/autologout/padding"

	loginUsernameCleanupPath = "regional/login/failed_login_username_cleanup"
	loginUsernameLimitPath   = "regional/login/failed_login_username_limit"
	loginUsernameUsePath     = "regional/login/failed_login_username_use"
	loginUserUsePath         = "regional/login/failed_login_user_use"
	loginUserWindowPath      = "regional/login/failed_login_user_window"

	userOptionsDormantPath = "profile/user-options/dormant"
)

type SysSettings struct {
	clientFactory ClientFactory
}

func NewSysSettings(clientFactory ClientFactory) *SysSettings {
	return &SysSettings{clientFactory}
}

// TimeSettings struct has timezone and date format
type TimeSettings struct {
	Timezone       string
	DateFormat     string
	TimeFormat     string
	DateTimeFormat string
}

type LoginSecuritySettings struct {
	FailedLoginUsernameCleanup uint64
	FailedLoginUsernameLimit   uint64
	FailedLoginUsernameUse     bool
	FailedLoginUserUse         uint64
	FailedLoginUserWindow      uint64
}

type GDPRSettings struct {
	Enabled bool
}

type MaintenanceModeSettings struct {
	Enabled bool
}

type AutologoutSettings struct {
	Enabled bool
	Timeout time.Duration
	Padding time.Duration
}

// GetTimeSettings returns new TimeSettings from settings service or err if can not get it
func (s *SysSettings) GetTimeSettings() (*TimeSettings, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	response, err := client.List(context.Background(), &pb.Request{Path: "regional/general/%"})
	if err != nil {
		return nil, err
	}

	timeSettings := TimeSettings{}
	timeSettings.Timezone = getSettingValue(response.Settings, defaultTimeZonePath)
	timeSettings.DateFormat = getSettingValue(response.Settings, defaultDateFormatPath)
	timeSettings.TimeFormat = getSettingValue(response.Settings, defaultTimeFormatPath)
	timeSettings.DateTimeFormat = fmt.Sprintf("%s %s", timeSettings.DateFormat, timeSettings.TimeFormat)
	return &timeSettings, nil
}

func (s *SysSettings) GetAutologoutSettings() (*AutologoutSettings, error) {
	autologoutSettings := &AutologoutSettings{}
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	response, err := client.List(context.Background(), &pb.Request{Path: "profile/autologout/%"})
	if err != nil {
		return nil, err
	}

	autologoutSettings.Enabled = getSettingValue(response.Settings, autoLogoutStatusPath) == "yes"

	timeout, err := strconv.ParseUint(getSettingValue(response.Settings, autoLogoutTimeoutPath), 10, 32)
	if err != nil {
		return nil, err
	}

	padding, err := strconv.ParseUint(getSettingValue(response.Settings, autoLogoutPaddingPath), 10, 32)
	if err != nil {
		return nil, err
	}

	autologoutSettings.Timeout = time.Duration(time.Duration(timeout) * time.Minute)
	autologoutSettings.Padding = time.Duration(time.Duration(padding) * time.Second)

	return autologoutSettings, nil
}

// GetGDPRSettings returns GDPR module settings from settings service or err if can not get it
func (s *SysSettings) GetGDPRSettings() (*GDPRSettings, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Get(context.Background(), &pb.Request{Path: gdprSettingPath})
	if err != nil {
		return nil, err
	}
	settings := GDPRSettings{}
	if response.Setting != nil {
		settings.Enabled = response.Setting.Value == "enable"
	}

	return &settings, nil
}

func (s *SysSettings) GetLoginSecuritySettings() (*LoginSecuritySettings, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}
	settings := LoginSecuritySettings{}

	response, err := client.List(context.Background(), &pb.Request{Path: "regional/login/%"})
	if err != nil {
		return &settings, err
	}

	settings.FailedLoginUsernameCleanup, _ = strconv.ParseUint(getSettingValue(response.Settings, loginUsernameCleanupPath), 10, 16)
	settings.FailedLoginUsernameLimit, _ = strconv.ParseUint(getSettingValue(response.Settings, loginUsernameLimitPath), 10, 16)
	settings.FailedLoginUsernameUse = "yes" == getSettingValue(response.Settings, loginUsernameUsePath)
	settings.FailedLoginUserUse, _ = strconv.ParseUint(getSettingValue(response.Settings, loginUserUsePath), 10, 16)
	settings.FailedLoginUserWindow, _ = strconv.ParseUint(getSettingValue(response.Settings, loginUserWindowPath), 10, 16)

	return &settings, nil
}

func (s *SysSettings) GetDormantDuration() (time.Duration, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return time.Nanosecond, err
	}
	response, err := client.Get(context.Background(), &pb.Request{Path: userOptionsDormantPath})
	if err != nil {
		return time.Nanosecond, err
	}

	strMonthDuration := response.Setting.Value
	monthCount, err := strconv.Atoi(strMonthDuration)
	if err != nil {
		return time.Nanosecond, err
	}
	hoursCount := monthCount * 30 * 24
	return time.ParseDuration(fmt.Sprintf("%dh", hoursCount))
}

// GetMaintenanceModeSettings returns maintenance mode settings from settings service or err if can not get it
func (s *SysSettings) GetMaintenanceModeSettings() (*MaintenanceModeSettings, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}

	response, err := client.Get(context.Background(), &pb.Request{Path: maintenancePath})
	if err != nil {
		return nil, err
	}

	settings := MaintenanceModeSettings{}
	if response.Setting != nil {
		settings.Enabled = response.Setting.Value == "enable"
	}

	return &settings, nil
}

// GetDefaulUserClassByRole returns default user class settings from settings service or err if can not get it
func (s *SysSettings) GetDefaultUserClassByRole(roleName string) (*string, error) {
	client, err := s.clientFactory.NewClient()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%s", "profile/default-user-classes", roleName)
	response, err := client.Get(context.Background(), &pb.Request{Path: path})
	if err != nil {
		return nil, err
	}
	return &response.Setting.Value, nil
}

func getSettingValue(settings []*pb.Setting, path string) string {
	for _, v := range settings {
		if v.Path == path {
			return v.Value
		}
	}
	return ""
}
