package backend

import "github.com/stretchr/testify/suite"

func New() *test {
	return &test{}
}

type test struct {
	suite.Suite
}

func (t *test) SetupTest() {

}
