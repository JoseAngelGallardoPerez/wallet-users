package csv

import (
	"fmt"
	"time"

	"github.com/Confialink/wallet-pkg-list_params"
	"github.com/Confialink/wallet-pkg-utils/csv"
	"github.com/Confialink/wallet-pkg-utils/timefmt"

	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services/csv/adminprofilesrow"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

// AdminProfiles service to generate csv file with admin profiles
type AdminProfiles struct {
	repository  repositories.RepositoryInterface
	sysSettings *syssettings.SysSettings
}

// NewAdminProfiles returns new AdminProfiles service
func NewAdminProfiles(repository repositories.RepositoryInterface, sysSettings *syssettings.SysSettings) *AdminProfiles {
	return &AdminProfiles{repository, sysSettings}
}

// GetFile returns generated csv file
func (s *AdminProfiles) GetFile(params *list_params.ListParams) (*csv.File, error) {
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
	file.Name = fmt.Sprintf("manager-profiles-%s.csv", formattedCurrentTime)

	file.WriteRow(adminProfilesHeader())

	for _, v := range users {
		rowBuilder := adminprofilesrow.NewRowBuilder(v, timeSettings)
		row := rowBuilder.Call()
		file.WriteRow(row)
	}

	return file, nil
}

func adminProfilesHeader() []string {
	return []string{
		"[User information] Profile Type",
		"[User information] First Name",
		"[User information] Last name",
		"[User information] Username",
		"[User information] Email",
		"[User information] Created",
		"[User information] Position",
		"[User information] Status",
		"[User information] Phone Number",
		"[User information] Class",
		"[Other] Internal Notes",
	}
}
