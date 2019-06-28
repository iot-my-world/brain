package test

import (
	"github.com/iot-my-world/brain/test/company"
	"github.com/iot-my-world/brain/test/system"
	"testing"
)

func Test(t *testing.T) {
	t.Log("System Tests")
	system.Test(t)
	t.Log("Company Tests")
	company.Test(t)
}
