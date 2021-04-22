package models

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	notificationspb "github.com/Confialink/wallet-notifications/rpc/proto/notifications"
	env_mods "github.com/Confialink/wallet-pkg-env_mods"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Confialink/wallet-users/internal/config"
	"github.com/Confialink/wallet-users/internal/db/types"
)

const (
	// TODO: remove
	RoleRoot   = "root"
	RoleAdmin  = "admin"
	RoleClient = "client"

	StatusPending  = "pending"
	StatusActive   = "active"
	StatusBlocked  = "blocked"
	StatusDormant  = "dormant"
	StatusCanceled = "canceled"

	DocumentTypePassport         = "passport"
	DocumentTypeDriverLicense    = "driver-license"
	DocumentTypeGovIssuedPhotoId = "gov-issued-photo-id"

	ProfileTypePersonal  = "Personal"
	ProfileTypeCorporate = "Corporate"

	letterBytes       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	minPasswordLength = 8

	ChallengeNameNewPasswordRequired = "new_password_required"
)

// User is the abstract user model
// Models should only be concerned with database schema, more strict checking should be put in validator.
// More detail you can find here: http://jinzhu.me/gorm/models.html#model-definition
// NOTE: If you want to split null and "", you should use *string instead of string.
type User struct {
	UID              string     `gorm:"primary_key:yes;column:uid;unique_index" json:"uid"`
	Email            string     `gorm:"column:email;unique_index" json:"email"`
	Username         string     `gorm:"column:username;unique;not null" json:"username"`
	Password         string     `gorm:"column:password;not null" json:"-"`
	FirstName        string     `gorm:"column:first_name" json:"firstName"`
	LastName         string     `gorm:"column:last_name" json:"lastName"`
	MiddleName       string     `gorm:"column:middle_name" json:"middleName"`
	Nickname         string     `gorm:"column:nickname" json:"nickname"`
	PhoneNumber      string     `gorm:"column:phone_number" json:"phoneNumber"`
	SmsPhoneNumber   *string    `gorm:"column:sms_phone_number" json:"smsPhoneNumber"`
	IsCorporate      *bool      `gorm:"not null; default:false" json:"isCorporate"`
	RoleName         string     `gorm:"column:role_name" json:"roleName"`
	ParentId         string     `gorm:"column:parent_id;default:null" json:"parentId"`
	Status           string     `gorm:"column:status" json:"status"`
	UserGroup        *UserGroup `gorm:"foreignkey:UserGroupId;association_foreignkey:ID;association_autoupdate:false" json:"userGroup"`
	UserGroupId      *uint64    `json:"userGroupId"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	LastLoginAt      *time.Time `gorm:"column:last_login_at;default:null" json:"lastLoginAt"`
	LastLoginIp      string     `json:"lastLoginIp"`
	ChallengeName    *string    `json:"challengeName"`
	IsPhoneConfirmed bool       `gorm:"not null; default:false" json:"isPhoneConfirmed"`
	IsEmailConfirmed bool       `gorm:"not null; default:false" json:"isEmailConfirmed"`
	CompanyDetails   Company    `gorm:"foreignkey:CompanyID;association_foreignkey:ID;association_autoupdate:false;association_autocreate:false;association_save_reference:false" json:"companyDetails"`
	CompanyID        *uint64    `json:"companyID"`
	ProfileImageID   *uint64    `json:"profileImageId"`

	PhysicalAddresses []*Address             `gorm:"-" json:"physicalAddresses"`
	MailingAddresses  []*Address             `gorm:"-" json:"mailingAddresses"`
	Attributes        map[string]interface{} `gorm:"-" json:"attributes"`

	UserDetails

	PermissionGroup *PermissionGroup `json:"-"` // TODO: check and remove if it is not needed
}

// UserDetails is
type UserDetails struct {
	ClassId                    json.Number `gorm:"column:class_id" json:"classId"`
	CountryOfResidenceIsoTwo   string      `gorm:"column:country_of_residence_iso_two"  json:"countryOfResidenceIsoTwo"`
	CountryOfCitizenshipIsoTwo string      `gorm:"column:country_of_citizenship_iso_two" json:"countryOfCitizenshipIsoTwo"`
	DateOfBirth                *types.Date `gorm:"column:date_of_birth" json:"dateOfBirth"`
	DocumentType               *string     `gorm:"column:document_type" json:"documentType"`
	DocumentPersonalId         string      `gorm:"column:document_personal_id" json:"documentPersonalId"`
	Fax                        string      `gorm:"column:fax" json:"fax"`
	HomePhoneNumber            string      `gorm:"column:home_phone_number" json:"homePhoneNumber"`
	InternalNotes              string      `gorm:"column:internal_notes" json:"internalNotes"`
	OfficePhoneNumber          string      `gorm:"column:office_phone_number" json:"officePhoneNumber"`
	Position                   string      `gorm:"column:position" json:"position"`
	BlockedUntil               *time.Time  `json:"blockedUntil"`
	LastActedAt                time.Time   `json:"lastActedAct"`
}

type Company struct {
	ID                uint64     `gorm:"column:id" json:"id"`
	CompanyName       string     `gorm:"column:company_name" json:"companyName"`
	CompanyType       string     `gorm:"column:company_type" json:"companyType"`
	CompanyRole       string     `gorm:"column:company_role" json:"companyRole"`
	DirectorFirstName string     `gorm:"column:director_first_name" json:"directorFirstName"`
	DirectorLastName  string     `gorm:"column:director_last_name" json:"directorLastName"`
	MaskName          string     `gorm:"column:mask_name" json:"-"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetProfileType returns profile type string
func (user *User) GetProfileType() string {
	if *user.IsCorporate {
		return ProfileTypeCorporate
	}
	return ProfileTypePersonal
}

// IsPasswordEncrypted checks if password encrypted
// U$..... - drupal md5
// $S$.... - drupal sha512
// $H$.... - drupal md5
// $P$.... - drupal md5
// $2a$... - bcrypt
func (user *User) IsPasswordEncrypted() bool {
	var password = user.Password
	var isbcrypt = regexp.MustCompile(`^U\$|\$S\$|\$H\$|\$P\$|\$2a\$`)
	var hashLength = 55 // number of characters in a hashed password

	if len(password) >= hashLength && isbcrypt.MatchString(password) {
		return true
	}
	return false
}

// GeneratePassword generates new password
func (user *User) GeneratePassword() error {
	b := make([]byte, minPasswordLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	password := string(b)

	if env_mods.GetMode(config.GetConf().GetServer().GetEnv()) == gin.DebugMode {
		log.Printf("The password is %s for %s", password, user.Username)
	}

	user.Password = password
	return nil
}

// GetAvailableStatuses returns available statuses
func GetAvailableStatuses() map[string]bool {
	return map[string]bool{
		StatusActive:   true,
		StatusPending:  true,
		StatusBlocked:  true,
		StatusDormant:  true,
		StatusCanceled: true,
	}
}

// GetAvailableDocumentTypes returns available document types
func GetAvailableDocumentTypes() map[string]bool {
	return map[string]bool{
		DocumentTypePassport:         true,
		DocumentTypeDriverLicense:    true,
		DocumentTypeGovIssuedPhotoId: true,
	}
}

// IsValidStatus checks if role is valid
func IsValidStatus(status string) bool {
	validStatuses := GetAvailableStatuses()
	return validStatuses[status]
}

// IsActive checks if user active
func (user *User) IsActive() bool {
	return user.Status == StatusActive
}

// IsInactive checks if user is not active
func (user *User) IsInactive() bool {
	return user.Status == StatusBlocked
}

// IsBlocked checks if user is blocked
func (user *User) IsBlocked() bool {
	return user.BlockedUntil != nil && user.BlockedUntil.After(time.Now())
}

// IsCanceled checks if user canceled
func (user *User) IsCanceled() bool {
	return user.Status == StatusCanceled
}

// IsPending checks if user pending
func (user *User) IsPending() bool {
	return user.Status == StatusPending
}

// IsDormant checks if user dormant
func (user *User) IsDormant() bool {
	return user.Status == StatusDormant
}

// IsAdmin checks if user rolename is admin
func (user *User) IsAdmin() bool {
	return user.RoleName == RoleAdmin
}

// IsRoot checks if user rolename is root
func (user *User) IsRoot() bool {
	return user.RoleName == RoleRoot
}

// IsClient checks if user rolename is root
func (user *User) IsClient() bool {
	return user.RoleName == RoleClient
}

// IsValidDocumentType checks if document type is valid
func IsValidDocumentType(docType string) bool {
	validDocumentTypes := GetAvailableDocumentTypes()
	if validDocumentTypes[docType] {
		return true
	}
	return false
}

// Approve checks users status is not pending and sets it to active
func (user *User) Approve() error {
	if !user.IsPending() {
		return errors.New("could not approve user with status different from 'pending'")
	}
	user.Status = StatusActive
	return nil
}

// Cancel checks users status is not pending and sets it to canceled
func (user *User) Cancel() error {
	if !user.IsPending() {
		return errors.New("could not cancel user with status different from 'pending'")
	}

	user.Status = StatusCanceled
	return nil
}

// Clear BlockedUntil property
func (user *User) ClearBlockedUntil() {
	user.BlockedUntil = nil
}

// CreateNotificationRequest creates notification request struct by event
func (user *User) CreateNotificationRequest(eventName string) *notificationspb.Request {
	return &notificationspb.Request{
		To:        user.Email,
		EventName: eventName,
		TemplateData: &notificationspb.TemplateData{
			UserName:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}
}

// GetClassIDAsInt64 returns ClassId as int64
func (user *User) GetClassIDAsInt64() int64 {
	strID := string(user.ClassId)
	intID, _ := strconv.ParseInt(strID, 10, 64)
	return intID
}

func (u *User) BeforeCreate() (err error) {
	u.UID = uuid.New().String()
	u.Username = u.UID
	return
}

func (*User) TableName() string {
	return "users"
}

func (u *UserDetails) GetDocumentType() string {
	if u.DocumentType != nil {
		return *u.DocumentType
	}
	return ""
}

func (u *User) GetCompanyDetails() *Company {
	return &Company{
		CompanyName:       u.CompanyDetails.CompanyName,
		CompanyRole:       u.CompanyDetails.CompanyRole,
		CompanyType:       u.CompanyDetails.CompanyType,
		DirectorLastName:  u.CompanyDetails.DirectorLastName,
		DirectorFirstName: u.CompanyDetails.DirectorFirstName,
	}
}

func (u *User) SetCompanyDetails(company *Company) {
	u.CompanyDetails.CompanyName = company.CompanyName
	u.CompanyDetails.CompanyType = company.CompanyType
	u.CompanyDetails.CompanyRole = company.CompanyRole
	u.CompanyDetails.DirectorFirstName = company.DirectorFirstName
	u.CompanyDetails.DirectorLastName = company.DirectorLastName
	return
}
