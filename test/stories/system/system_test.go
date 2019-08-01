package system

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, New())
}
