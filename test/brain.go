package test

import (
	"github.com/iot-my-world/brain/test/system"
	"testing"
)

func Test(t *testing.T) {
	t.Log("System Tests")
	system.Test(t)
}
