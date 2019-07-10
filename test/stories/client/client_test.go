package client

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestClient(t *testing.T) {
	suite.Run(t, New())
}
