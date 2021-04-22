package syssettings

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSysSettings(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SysSettings Suite")
}
