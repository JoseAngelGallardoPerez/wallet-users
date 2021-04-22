package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/inconshreveable/log15"
	"go.uber.org/dig"

	"github.com/Confialink/wallet-users/internal/db/models"
	"github.com/Confialink/wallet-users/internal/db/repositories"
	"github.com/Confialink/wallet-users/internal/services"
	"github.com/Confialink/wallet-users/internal/services/users"
	"github.com/Confialink/wallet-users/internal/validators"
)

// CreateRootUser create root user or change email & password for root user if root user already exists
type CreateRootUser struct {
	name        string
	usage       string
	description string
	container   *dig.Container
	logger      log15.Logger
}

type user struct {
	email           string
	password        string
	usersRepository *repositories.UsersRepository
	userService     *users.UserService
	passwordService *services.Password
}

func Init(container *dig.Container) *CreateRootUser {
	return &CreateRootUser{
		name:        "create-root-user",
		usage:       "create-root-user",
		description: "Create root user or change email & password for root user if root user already exists. Usage: \"create-root-user?email=emailValue&password=passwordValue\"",
		container:   container,
	}
}

func (c *CreateRootUser) Name() string {
	return c.name
}

func (c *CreateRootUser) Description() string {
	return fmt.Sprintf("%s:\nUsage:\t%s\nDescription:\t%s\n\n", c.name, c.usage, c.description)
}

func (c *CreateRootUser) Handle(args url.Values) {
	err := c.container.Invoke(func(
		userService *users.UserService,
		passwordService *services.Password,
		repo *repositories.UsersRepository,
		logger log15.Logger,
	) {
		c.logger = logger
		err := c.createUser(args, &user{usersRepository: repo, userService: userService, passwordService: passwordService})
		if err != nil {
			c.logger.Error(err.Error())
			os.Exit(1)
		}
	})
	if err != nil {
		c.logger.Error(err.Error())
		os.Exit(1)
	}
}

func (c *CreateRootUser) createUser(args url.Values, rootUser *user) error {
	email := args.Get("email")
	if email == "" {
		c.logger.Error("parameter \"email\" is required\n usage: \"create-root-user?email=emailValue&password=passwordValue\"")
		os.Exit(1)
	}
	password := args.Get("password")
	if password == "" {
		c.logger.Error("parameter \"password\" is required\n usage: \"create-root-user?email=emailValue&password=passwordValue\"")
		os.Exit(1)
	}
	rootUser.email = email
	rootUser.password = password
	// Find user
	user, err := rootUser.findRootUser()
	if err != nil {
		return err
	}
	if user != nil {
		// Update user
		_, err := rootUser.updateRootUser(user)
		if err != nil {
			return err
		}
		return nil
	}
	// Create user
	_, err = rootUser.createRootUser()
	if err != nil {
		return err
	}

	return nil
}

func (u *user) findRootUser() (*models.User, error) {
	result, err := u.usersRepository.FindByRoleName(models.RoleRoot)
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}

	return nil, nil
}

func (u *user) createRootUser() (*models.User, error) {
	// Checks if the query entry is valid
	rootValidator := validators.RootValidator{}
	rootValidator.Email = u.email
	rootValidator.Password = u.password
	if err := rootValidator.Call(); err != nil {
		return nil, err
	}
	// Create new user
	user, err := u.userService.Create(&rootValidator.UserModel, true, false, nil)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) updateRootUser(user *models.User) (*models.User, error) {
	// Checks if the query entry is valid
	rootValidator := validators.RootValidator{}
	rootValidator.Email = u.email
	rootValidator.Password = u.password
	if err := rootValidator.Call(); err != nil {
		return nil, err
	}
	hash, err := u.passwordService.UserHashPassword(rootValidator.UserModel.Password)

	if err != nil {
		return nil, err
	}

	user.Email = rootValidator.Email
	user.Password = hash

	// Update user
	user, err = u.usersRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
