package services

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	// hashCount is the standard log2 number of iterations for password stretching.
	hashCount = 15
	// minHashCount is the minimum allowed log2 number of iterations for password stretching
	minHashCount = 7
	// maxHashCount is the maximum allowed log2 number of iterations for password stretching.
	maxHashCount = 30
	// hashLength is the expected (and maximum) number of characters in a hashed password.
	hashLength = 55
)

// Password is secure password hashing for user authentication
type Password struct{}

// NewPassword creates new password service
func NewPassword() *Password {
	return &Password{}
}

// passwordItoa64 returns a string for mapping an int to the corresponding base 64 character.
func (p *Password) passwordItoa64() string {
	return "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
}

// UserHashPassword hash a password using a secure hash.
func (p *Password) UserHashPassword(password string) (string, error) {
	salt := p.passwordGenerateSalt(hashCount)

	hash, err := p.passwordCrypt("sha512", password, salt)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// UserCheckPassword check whether a plain text password matches a stored hashed password
func (p *Password) UserCheckPassword(password, hashedPassword string) error {
	var passwordType = hashedPassword[0:3]

	if passwordType[0:2] == "U$" {
		hashedPassword = hashedPassword[1:]
		password = p.getMD5Hash(password)
	}

	switch passwordType {
	case "$S$":
		return p.compareHashAndPasswordByAlgo("sha512", password, hashedPassword)
	case "$H$", "$P$":
		return p.compareHashAndPasswordByAlgo("md5", password, hashedPassword)
	case "$2a":
		bytePassword := []byte(password)
		byteHashedPassword := []byte(hashedPassword)
		return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	default:
		return fmt.Errorf("unknown type of password %s", passwordType)
	}
}

// compareHashAndPassword compares a hashed password with its possible hashed equivalent.
func (p *Password) compareHashAndPasswordByAlgo(algo, password, hashedPassword string) error {
	hash, err := p.passwordCrypt(algo, password, hashedPassword)
	if nil != err {
		return err
	}

	if len(hash) < minHashCount {
		return errors.New("hash too short to be a bcrypted password")
	}

	if hashedPassword != hash {
		return errors.New("hashedPassword is not the hash of the given password")
	}

	return nil
}

// passwordCrypt hash a password using a secure stretched hash.
func (p *Password) passwordCrypt(algo, password, setting string) (string, error) {
	var output bytes.Buffer

	if len(password) > 512 { // Prevent DoS attacks by refusing to hash large passwords.
		return "", errors.New("password cannot be longer than 512")
	}
	// The first 12 characters of an existing hash are its setting string.
	setting = setting[0:12]
	if string(setting[0]) != "$" || string(setting[2]) != "$" {
		return "", errors.New("password is not valid")
	}

	output.WriteString(setting)

	countLog2 := p.passwordGetCountLog2(setting)
	// Hashes may be imported from elsewhere, so we allow != DRUPAL_HASH_COUNT
	if countLog2 < minHashCount || countLog2 > maxHashCount {
		return "", errors.New("password is not valid")
	}

	salt := setting[4:12]
	// Hashes must have an 8 character salt.
	if len(salt) != 8 {
		return "", errors.New("password is not valid")
	}

	// Convert the base 2 logarithm into an integer.
	count := 1 << countLog2

	data := fmt.Sprintf("%s%s", salt, password)

	hash := p.hash(algo, data)
	for i := 0; i < count; i++ {
		data := fmt.Sprintf("%s%s", hash, password)
		hash = p.hash(algo, data)
	}

	hashlength := len(hash)
	division := 8 * float64(hashlength) / 6
	expected := 12 + math.Ceil(division)

	output.WriteString(p.passwordBase64Encode(hash, hashlength))

	outputStr := output.String()
	if len(outputStr) == int(expected) {
		return outputStr[0:hashLength], nil
	}

	return "", errors.New("password is not valid")
}

func (p *Password) passwordBase64Encode(input []byte, hashlength int) string {
	var output bytes.Buffer
	var byteChar string
	i := 0
	itoa64 := p.passwordItoa64()

	for {
		value := rune(input[i])
		i++

		byteChar = string(itoa64[value&0x3f])
		output.WriteString(byteChar)

		if i < hashlength {
			value |= rune(input[i]) << 8
		}

		byteChar = string(itoa64[(value>>6)&0x3f])
		output.WriteString(byteChar)

		if i >= hashlength {
			break
		} else {
			i++
		}

		if i < hashlength {
			value |= rune(input[i]) << 16
		}

		output.WriteString(string(itoa64[(value>>12)&0x3f]))

		if i >= hashlength {
			break
		} else {
			i++
		}

		output.WriteString(string(itoa64[(value>>18)&0x3f]))

		if i >= hashlength {
			break
		}
	}

	return output.String()
}

func (p *Password) hash(algorithm, data string) []byte {
	var hash []byte
	switch algorithm {
	case "sha512":
		hash = p.getSHA512Sum(data)
		break
	case "md5":
		hash = p.getMD5Sum(data)
		break
	// case "bcrypt":
	// 	hash, _ = bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	// 	break
	default:
		hash = p.getSHA512Sum(data)
	}
	return hash
}

func (p *Password) passwordGenerateSalt(countLog2 int) string {
	var output bytes.Buffer

	output.WriteString("$S$")

	// Ensure that countLog2 is within set bounds.
	countLog2 = p.passwordEnforceLog2Boundaries(countLog2)

	// We encode the final log2 iteration count in base 64.
	itoa64 := p.passwordItoa64()
	output.WriteString(string(itoa64[countLog2]))

	// 6 bytes is the standard salt for a portable phpass hash.
	output.WriteString(p.passwordBase64Encode(p.randomBytes(6), 6))
	return output.String()
}

func (p *Password) randomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func (p *Password) passwordEnforceLog2Boundaries(countLog2 int) int {
	if countLog2 < minHashCount {
		return minHashCount
	} else if countLog2 > maxHashCount {
		return maxHashCount
	}
	return countLog2
}

func (p *Password) getMD5Hash(text string) string {
	sum := p.getMD5Sum(text)
	return hex.EncodeToString(sum)
}

func (p *Password) getMD5Sum(text string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

func (p *Password) getSHA512Hash(text string) string {
	sum := p.getSHA512Sum(text)
	return hex.EncodeToString(sum)
}

func (p *Password) getSHA512Sum(text string) []byte {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hasher.Sum(nil)
}

// passwordGetCountLog2 parse the log2 iteration count from a stored hash or setting string.
func (p *Password) passwordGetCountLog2(setting string) uint64 {
	itoa64 := p.passwordItoa64()
	substr := string(setting[3])

	i := strings.Index(itoa64, substr)

	u64 := uint64(i)
	return u64
}
