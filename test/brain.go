package test

import (
	"github.com/iot-my-world/brain/test/stories/company"
	"github.com/iot-my-world/brain/test/stories/public"
	"github.com/iot-my-world/brain/test/stories/system"
	"testing"
)

func Test(t *testing.T) {
	t.Log("System Tests")
	system.Test(t)
	t.Log("Company Tests")
	company.Test(t)
	t.Log("Public Tests")
	public.Test(t)
}
