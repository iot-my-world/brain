package system

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestAll(t *testing.T) {
	suite.Run(t, New())
}
