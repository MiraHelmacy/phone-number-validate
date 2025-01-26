package validate

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
)

type PinpointActionWrapper struct {
	client *pinpoint.Client
}

func (wrapper *PinpointActionWrapper) PhoneNumberValidate(phoneNumberValidateRequestInput *pinpoint.PhoneNumberValidateInput) (phoneNumberValidateResponse *pinpoint.PhoneNumberValidateOutput, err error) {
	phoneNumberValidateResponse, err = wrapper.client.PhoneNumberValidate(context.TODO(), phoneNumberValidateRequestInput)
	return
}

func NewPinpointWrapper() (wrapper *PinpointActionWrapper, err error) {
	cfg, loadErr := config.LoadDefaultConfig(context.TODO())
	if loadErr != nil {
		wrapper = nil
		err = loadErr
		return
	}
	client := pinpoint.NewFromConfig(cfg)
	wrapper = &PinpointActionWrapper{client: client}
	return
}
