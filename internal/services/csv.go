package services

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Confialink/wallet-pkg-utils/pointer"
	"github.com/Confialink/wallet-pkg-utils/value"

	"github.com/Confialink/wallet-users/internal/db/models"
)

// Csv is an empty structure
type Csv struct{}

// CsvHandler is interface for csv functionality
type CsvHandler interface {
	UsersToCsv(users []*models.User) *bytes.Buffer
	CsvHeader() []string
}

// UsersToCsv exports users to csv
func (s Csv) UsersToCsv(items []*models.User) *bytes.Buffer {

	header := s.CsvHeader()

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.Write(header)
	for i, v := range items {
		if i == 0 { // headers
			continue
		}
		row := userToRow(v)
		wr.Write(row)
	}

	wr.Flush()
	return b
}

// CsvToUsers imports users from csv
func (s Csv) CsvToUsers(b *bytes.Buffer) ([]*models.User, error) {
	reader := csv.NewReader(bufio.NewReader(b))

	var users []*models.User

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		details := models.UserDetails{
			ClassId:                    json.Number(line[12]),
			CountryOfResidenceIsoTwo:   line[13],
			CountryOfCitizenshipIsoTwo: line[14],
			DocumentType:               &line[18],
			DocumentPersonalId:         line[19],
			Fax:                        line[20],
			HomePhoneNumber:            line[21],
			InternalNotes:              line[22],
			OfficePhoneNumber:          line[23],
			Position:                   line[24],
		}

		cd := models.Company{
			CompanyName: line[7],
		}

		groupID := stringToUint64(line[11])

		user := models.User{
			UID:            line[0],
			Email:          line[1],
			Username:       line[2],
			Password:       line[3],
			FirstName:      line[4],
			LastName:       line[5],
			PhoneNumber:    preparePhoneNumber(line[6]),
			IsCorporate:    pointer.ToBool(stringToBool(line[8])),
			RoleName:       line[9],
			Status:         line[10],
			UserGroupId:    &groupID,
			UserDetails:    details,
			CompanyDetails: cd,
		}

		users = append(users, &user)
	}

	return users, nil
}

// CsvHeader return header for csv
func (s Csv) CsvHeader() []string {
	return []string{
		"UID", // line 0
		"Email",
		"Username",
		"Password",
		"FirstName",
		"LastName",
		"PhoneNumber",
		"CompanyName",
		"IsCorporate",
		"RoleName",
		"Status", // line 10
		"UserGroupId",
		// User Details
		"ClassId",
		"CountryOfResidenceIso2",
		"CountryOfCitizenshipIso2",
		"DateOfBirthYear",
		"DateOfBirthMonth",
		"DateOfBirthDay",
		"DocumentType",
		"DocumentPersonalId",
		"Fax", // line 20
		"HomePhoneNumber",
		"InternalNotes",
		"OfficePhoneNumber",
		"Position",
		// Physical Adress
		"PaZipPostalCode",
		"PaAddress",
		"PaAddress2ndLine",
		"PaCity",
		"PaCountryIso2",
		"PaStateProvRegion", // line 30
		// Mailing Address
		"MaZipPostalCode",
		"MaStateProvRegion",
		"MaPhoneNumber",
		"MaName",
		"MaCountryIso2",
		"MaCity",
		"MaAddress",
		"MaAddress2ndLine",
		"MaAsPhysical",
		// Benificial Owner
		"BoFullName", // line 40
		"BoPhoneNumber",
		"BoDateOfBirthYear",
		"BoDateOfBirthMonth",
		"BoDateOfBirthDay",
		"BoDocumentPersonalId",
		"BoDocumentType",
		"BoAddress",
		"BoRealationship",
	}
}

func preparePhoneNumber(number string) string {
	if !strings.HasPrefix(number, "+") {
		var b strings.Builder
		b.WriteString("+")
		b.WriteString(number)
		return b.String()
	}
	return number
}

func userToRow(user *models.User) []string {
	return []string{
		user.UID,
		user.Email,
		user.Username,
		user.Password,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.CompanyDetails.CompanyName,
		fmt.Sprint(user.IsCorporate),
		user.RoleName,
		user.Status,
		fmt.Sprint(value.FromUint64(user.UserGroupId)),
		fmt.Sprint(user.ClassId),
		user.CountryOfResidenceIsoTwo,
		user.CountryOfCitizenshipIsoTwo,
		user.GetDocumentType(),
		user.DocumentPersonalId,
		user.Fax,
		user.HomePhoneNumber,
		user.InternalNotes,
		user.OfficePhoneNumber,
		user.Position,
	}
}

func stringToUint64(str string) uint64 {
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	return u64
}

func stringToBool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		fmt.Println(err)
	}
	return b
}
