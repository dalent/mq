package rabbit

import (
	"github.com/dalent/mq"
	"testing"
)

func TestingManager(t *testing.T) {
	pro, err := mq.NewProvider(mq.RABBIT)
	if err != nil {
		t.Fail()
	}
	err = pro.SetUp("amqp://guest:guest@192.168.11.125:5672/", "sms")
	if err != nil {
		t.Fail()
	}
	defer pro.Close()
}
