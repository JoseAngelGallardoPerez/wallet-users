package auth

import (
	"testing"
)

func TestCheckIp(t *testing.T) {
	// prepare data
	// ip := "172.16.12.13"          // random
	// randomCreatedAt := time.Now() // random
	// blockedIp := &models.BlockedIp{ID: 1, IP: ip, CreatedAt: randomCreatedAt, BlockedUntil: randomCreatedAt}
	// service := newTestBlocker()

	/* blockedIp is not found */
	// createDbSelectQueryMockBlockedIp(nil)
	// assert.Nil(t, service.CheckIp(ip), "IP must not be blocked")

	/* BlockedUntil is in the future */
	// blockedIp.BlockedUntil = time.Now().Local().Add(time.Hour * time.Duration(1)) // +1 hour
	// createDbSelectQueryMockBlockedIp(blockedIp)
	// assert.Error(t, service.CheckIp(ip), "IP must be blocked")

	/* BlockedUntil is in the past */
	// blockedIp.BlockedUntil = time.Now().Local().Add(-1 * time.Hour) // -1 hour
	// createDbSelectQueryMockBlockedIp(blockedIp)
	// helpers.DbMock.GetDbMock().ExpectExec("DELETE FROM `" + blockedIp.TableName() + "`").
	// 	WithArgs(blockedIp.ID).
	// 	WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.Nil(t, service.CheckIp(ip), "IP must be unblocked") // IP will be unblocked
}

func TestCheckUser(t *testing.T) {
	// service := newTestBlocker()
	// dbMock := helpers.DbMock.GetDbMock()
	// user := &models.User{
	// 	Username: "random_user_name",
	// 	Status:   models.StatusBlocked,
	// }

	/* a user is blocked and BlockedUntil is in the past */
	// blockedUntil := time.Now().Local().Add(time.Hour * time.Duration(-1)) // -1 hour
	// user.BlockedUntil = &blockedUntil
	// createDbSelectQueryMockUser(user)
	// dbMock.ExpectExec("UPDATE `" + user.TableName() + "`").WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.Nil(t, service.CheckUser(user.Username), "User must be unblocked")

	/* a user is blocked and BlockedUntil is in the future */
	// blockedUntil = time.Now().Local().Add(time.Hour * time.Duration(1)) // +1 hour
	// user.BlockedUntil = &blockedUntil
	// createDbSelectQueryMockUser(user)
	// assert.Error(t, service.CheckUser(user.Username), "User must not be unblocked")

	/* a user is active */
	// user.Status = models.StatusActive
	// createDbSelectQueryMockUser(user)
	// assert.Nil(t, service.CheckUser(user.Username), "User must be unblocked")

	/* a user is not found */
	// createDbSelectQueryMockUser(nil)
	// assert.Nil(t, service.CheckUser(user.Username), "User must be not found")
}

func TestAddIpFailAttempt(t *testing.T) {
	// ip := "172.16.12.13" //random
	// failAttempt := &models.FailAuthAttempt{}
	// service := newTestBlocker()
	// service.loginSettings = &syssettings.LoginSecuritySettings{}
	// dbMock := helpers.DbMock.GetDbMock()

	/* There is first login attempt. FailedLoginUserUse is zero */
	// dbMock.ExpectExec("INSERT INTO " + failAttempt.TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.True(t, service.AddIpFailAttempt(ip), "IpFailAttempt must be added")

	/* There is first login attempt. FailedLoginUserUse is not zero */
	// createDbSelectCountQueryMockFailAttempt(0)
	// service.loginSettings.FailedLoginUserUse = 5
	// dbMock.ExpectExec("INSERT INTO " + failAttempt.TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.True(t, service.AddIpFailAttempt(ip), "IpFailAttempt must be added")

	/* There is not first login attempt. FailedLoginUserUse is not zero and more then count of attempts. */
	// createDbSelectCountQueryMockFailAttempt(3)
	// service.loginSettings.FailedLoginUserUse = 5
	// dbMock.ExpectExec("INSERT INTO " + failAttempt.TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.True(t, service.AddIpFailAttempt(ip), "IpFailAttempt must be added")

	/* There is not first login attempt. FailedLoginUserUse is not zero. */
	// createDbSelectCountQueryMockFailAttempt(4) // 4(in DB) + 1(current attempt) = 5
	// service.loginSettings.FailedLoginUserUse = 5
	// dbMock.ExpectExec("INSERT INTO " + (&models.BlockedIp{}).TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// dbMock.ExpectExec("DELETE FROM `" + failAttempt.TableName() + "`").WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.False(t, service.AddIpFailAttempt(ip), "IP must be blocked")
}

// add User failAttempt to DB
func TestAddUserFailAttempt(t *testing.T) {
	// service := newTestBlocker()
	// failAttempt := &models.FailAuthAttempt{}
	// dbMock := helpers.DbMock.GetDbMock()
	// service.loginSettings = &syssettings.LoginSecuritySettings{}
	// service.loginSettings.FailedLoginUsernameUse = true
	// service.loginSettings.FailedLoginUsernameLimit = 5
	// user := &models.User{
	// 	Email:  "random_email@example.com",
	// 	Status: models.StatusBlocked,
	// }

	/* a user has 0 fail attempts */
	// createDbSelectQueryMockUser(user)
	// createDbSelectCountQueryMockFailAttempt(0)
	// dbMock.ExpectExec("INSERT INTO " + failAttempt.TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.True(t, service.AddUserFailAttempt(user.Email, ""), "UserFailAttempt must be added")

	/* a user has failed attempts */
	// mockNotificationHandler(user, service)
	// createDbSelectQueryMockUser(user)
	// service.loginSettings.FailedLoginUserUse = 5
	// createDbSelectCountQueryMockFailAttempt(4)
	// dbMock.ExpectExec("UPDATE `" + user.TableName() + "`").WillReturnResult(sqlmock.NewResult(1, 1))
	// dbMock.ExpectExec("UPDATE " + failAttempt.TableName()).WillReturnResult(sqlmock.NewResult(1, 1))
	// assert.False(t, service.AddUserFailAttempt(user.Email, ""), "User must be blocked")

	/* FailedLoginUsernameUse is disabled */
	// service.loginSettings.FailedLoginUsernameUse = false
	// createDbSelectQueryMockUser(user)
	// assert.False(t, service.AddUserFailAttempt(user.Email, ""), "UserFailAttempt must not be added")
}
