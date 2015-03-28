package rabbit

import (
	"fmt"
	"testing"
)

type smsService struct {
}

func (*smsService) Do(_ []byte) {
}
func GetSmsService() *smsService {
	return &smsService{}
}

func TestClose(t *testing.T) {
	provider := new(rabbitAdapter)
	err := provider.SetUp("amqp://guest:guest@192.168.11.125:5672/", "sms")
	if err != nil {
		return
	}
	fmt.Println("connect success")
	defer provider.Close()
}
