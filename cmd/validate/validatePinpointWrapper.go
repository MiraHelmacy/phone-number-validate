package validate

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
)

// pinpoint client wrapper
type PinpointActionWrapper struct {
	client *pinpoint.Client
}

// validate the request
func (wrapper *PinpointActionWrapper) PhoneNumberValidate(phoneNumberValidateRequestInput *pinpoint.PhoneNumberValidateInput) (phoneNumberValidateResponse *pinpoint.PhoneNumberValidateOutput, err error) {
	phoneNumberValidateResponse, err = wrapper.client.PhoneNumberValidate(context.TODO(), phoneNumberValidateRequestInput)
	return
}

// create a new wrapper
func NewPinpointWrapper() (wrapper *PinpointActionWrapper, err error) {
	//create a config with context.TODO
	cfg, loadErr := config.LoadDefaultConfig(context.TODO())

	//check the error
	if loadErr != nil {
		wrapper = nil
		err = loadErr
		return
	}

	//create a new client
	client := pinpoint.NewFromConfig(cfg)

	//create wrapper
	wrapper = &PinpointActionWrapper{client: client}
	return
}
