package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Confialink/wallet-pkg-utils"
	pbSettings "github.com/Confialink/wallet-settings/rpc/proto/settings"
	"github.com/stretchr/testify/assert"

	"github.com/Confialink/wallet-users/internal/services/syssettings"
	"github.com/Confialink/wallet-users/internal/services/syssettings/mocks"
	"github.com/Confialink/wallet-users/internal/tests/mocks/vendor-mocks/rpc/settings"
)

// AutologoutTTLResolver
func TestAutologoutResolveByTokenSubject(t *testing.T) {
	var testTable = []struct {
		caseNumber        uint16
		autologoutStatus  string
		autologoutTimeout string // minutes
		autologoutPadding string // seconds
		isError           bool   // expected behavior
		expectedTtl       string
		claim             string
		rpcClientError    error
	}{
		{
			1,
			"no",
			"10",  // random
			"100", // random
			false,
			ClaimAccessTokenExp,
			ClaimAccessSub,
			nil,
		},
		{
			2,
			"yes",
			"40",
			"120",
			false,
			"21m10s", // (40m+120s/2) + 10s
			ClaimAccessSub,
			nil,
		},
		{
			3,
			"no",
			"10",  // random
			"120", // random
			false,
			ClaimRefreshTokenExp,
			ClaimRefreshSub,
			nil,
		},
		{
			4,
			"yes",
			"40",
			"120",
			false,
			"1h3m10s", // (40m+120s/2) + 10s + 40m+120s
			ClaimRefreshSub,
			nil,
		},
		{
			5,
			"no", // random
			"22",
			"120",
			true,
			"0",
			ClaimAccessSub,
			errors.New("random text"),
		},
		{
			6,
			"no",
			"40", // random
			"60",
			false,
			"0",
			"random",
			nil,
		},
		{
			7,
			"yes",
			"40", // random
			"60",
			false,
			"0",
			"random",
			nil,
		},
	}

	for _, testData := range testTable {
		client := &settings.MockSettingsHandler{}
		req := &pbSettings.Request{Path: "profile/autologout/%"}
		statusSetting := &pbSettings.Setting{Path: "profile/autologout/status", Value: testData.autologoutStatus}
		timeoutSetting := &pbSettings.Setting{Path: "profile/autologout/timeout", Value: testData.autologoutTimeout}
		paddingSettings := &pbSettings.Setting{Path: "profile/autologout/padding", Value: testData.autologoutPadding}
		respSettings := []*pbSettings.Setting{statusSetting, timeoutSetting, paddingSettings}
		resp := &pbSettings.Response{Settings: respSettings}
		client.On("List", context.Background(), req).Return(resp, testData.rpcClientError)
		clientFactory := &mocks.ClientFactory{}
		clientFactory.On("NewClient").Return(client, nil)
		sysSettings := syssettings.NewSysSettings(clientFactory)
		resolver := NewAutologoutTTLResolver(sysSettings)
		ttl, err := resolver.ResolveByTokenSubject(testData.claim)

		errorIsExists := err != nil
		assert.False(t, errorIsExists && !testData.isError, fmt.Sprintf("there must be no any errors. case #: %d", testData.caseNumber))
		assert.False(t, !errorIsExists && testData.isError, fmt.Sprintf("there must be an error. case #: %d", testData.caseNumber))
		assert.Equal(t, utils.MustParseDuration(testData.expectedTtl), ttl, fmt.Sprintf("ttl is not expected. case #: %d", testData.caseNumber))
	}
}

// FixedValueTTLResolver
func TestFixedValueResolveByTokenSubject(t *testing.T) {
	var testTable = []struct {
		caseNumber      uint16
		claim           string
		refreshTokenTTL string
		accessTokenTTL  string
		expectedTtl     string
	}{
		{
			1,
			ClaimAccessSub,
			"10m", // random
			"11m",
			"11m",
		},
		{
			2,
			ClaimRefreshSub,
			"10m",
			"11m", // random
			"10m",
		},
		{
			3,
			"random",
			"10m", // random
			"11m", // random
			"0",
		},
	}

	for _, testData := range testTable {
		resolver := NewFixedValueTTLResolver(testData.refreshTokenTTL, testData.accessTokenTTL)
		ttl, err := resolver.ResolveByTokenSubject(testData.claim)

		assert.True(t, err == nil, fmt.Sprintf("there must be no any errors. case #: %d", testData.caseNumber))
		assert.Equal(t, utils.MustParseDuration(testData.expectedTtl), ttl, fmt.Sprintf("ttl is not expected. case #: %d", testData.caseNumber))
	}
}
