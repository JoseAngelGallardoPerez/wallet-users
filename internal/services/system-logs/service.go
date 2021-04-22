package system_logs

import (
	"encoding/json"
	"time"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/inconshreveable/log15"
)

type SystemLogsService struct {
	logsServiceWrap *logsServiceWrap
	recoverer       func()
	logger          log15.Logger
}

func NewSystemLogsService(
	recoverer func(),
	logger log15.Logger,
) *SystemLogsService {
	return &SystemLogsService{
		logsServiceWrap: newLogsServiceWrap(),
		recoverer:       recoverer,
		logger:          logger,
	}
}

func (self *SystemLogsService) LogCreateUserProfileAsync(
	user *models.User, userId string,
) {
	go self.LogCreateUserProfile(user, userId)
}

func (self *SystemLogsService) LogModifyUserProfileAsync(
	old *models.User, new *models.User, userId string,
) {
	go self.LogModifyUserProfile(old, new, userId)
}

func (self *SystemLogsService) LogCreateUserProfile(
	user *models.User, userId string,
) {
	data, err := json.Marshal(user)
	if err != nil {
		self.logger.Error("Can't marshal json", err)
		return
	}

	defer self.recoverer()
	self.logsServiceWrap.createLog(
		SubjectCreateUserProfile,
		userId,
		time.Now().Format(time.RFC3339),
		DataTitleProfileDetails,
		data,
	)
}

func (self *SystemLogsService) LogModifyUserProfile(
	old *models.User, new *models.User, userId string,
) {
	defer self.recoverer()

	data, err := json.Marshal(map[string]interface{}{
		"old": old,
		"new": new,
	})

	if err != nil {
		self.logger.Error("Can't marshall json", err)
		return
	}

	self.logsServiceWrap.createLog(
		SubjectModifyUserProfile,
		userId,
		time.Now().Format(time.RFC3339),
		DataTitleProfileDetails+": "+old.Username,
		data,
	)
}
