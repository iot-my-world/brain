package company

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestCompany(t *testing.T) {
	suite.Run(t, New())
}
