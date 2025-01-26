package validate

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
)

type ValidateCmdCliArgs struct {
	PhoneNumbers []string
	Short        bool
	OnlyInvalid  bool
}

func (args *ValidateCmdCliArgs) PhoneNumberValidateRequests() (requests []*pinpoint.PhoneNumberValidateInput) {
	requests = make([]*pinpoint.PhoneNumberValidateInput, 0)
	for _, phoneNumber := range args.PhoneNumbers {
		numberValidateRequest := &types.NumberValidateRequest{}
		numberValidateRequest.PhoneNumber = aws.String(phoneNumber)
		requestInput := &pinpoint.PhoneNumberValidateInput{}
		requestInput.NumberValidateRequest = numberValidateRequest
		requests = append(requests, requestInput)
	}
	return
}

func (args ValidateCmdCliArgs) Validate() (isValid bool) {
	phoneNumbersLength := len(args.PhoneNumbers)
	validPhoneNumberLength := phoneNumbersLength > 0

	isValid = validPhoneNumberLength
	return
}
