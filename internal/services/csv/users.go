package csv

import (
	"fmt"
	"time"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-pkg-utils/csv"
	"github.com/Confialink/wallet-pkg-utils/timefmt"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/csv/userprofilesrow"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

// Users service for generating csv file for user profiles
type Users struct {
	repository  repositories.RepositoryInterface
	sysSettings *syssettings.SysSettings
}

// NewUsers returns new Users csv service
func NewUsers(repository repositories.RepositoryInterface, sysSettings *syssettings.SysSettings) *Users {
	return &Users{repository, sysSettings}
}

// GetFile returns generated files for users by passed params
func (s *Users) GetFile(params *list_params.ListParams) (*csv.File, error) {
	users, err := s.repository.GetUsersRepository().GetList(params)
	if err != nil {
		return nil, err
	}
	currentTime := time.Now()
	timeSettings, err := s.sysSettings.GetTimeSettings()
	if err != nil {
		return nil, err
	}

	file := csv.NewFile()
	formattedCurrentTime := timefmt.FormatFilenameWithTime(currentTime, timeSettings.Timezone)
	file.Name = fmt.Sprintf("user-profiles-%s.csv", formattedCurrentTime)

	file.WriteRow(getUsersHeader())

	for _, v := range users {
		rowBuilder := userprofilesrow.NewUserProfileRowBuilder(v, timeSettings)
		row := rowBuilder.Call()
		file.WriteRow(row)
	}

	return file, nil
}
