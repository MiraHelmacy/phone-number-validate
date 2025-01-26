package validate

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint"
	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
)

// Validate Cli Args struct
type ValidateCmdCliArgs struct {
	PhoneNumbers []string //phone numbers to validate
	Short        bool     //should short option be used?
	OnlyInvalid  bool     //should valid numbers be shown?
}

// validate the cli args
func (args ValidateCmdCliArgs) Validate() (isValid bool) {
	//get the length of the phone numbers
	phoneNumbersLength := len(args.PhoneNumbers)

	//no phone numbers provided if false
	validPhoneNumberLength := phoneNumbersLength > 0

	//valid only if phone numbers is greater than 0
	isValid = validPhoneNumberLength
	return
}

// for the cli args, create a list of phone number validate requests
func (args *ValidateCmdCliArgs) PhoneNumberValidateRequests() (requests []*pinpoint.PhoneNumberValidateInput) {
	//empty list of requests
	requests = make([]*pinpoint.PhoneNumberValidateInput, 0)

	//for each phone number
	for _, phoneNumber := range args.PhoneNumbers {
		//create a NumberValidateRequest
		numberValidateRequest := &types.NumberValidateRequest{}

		//assign the phone number
		numberValidateRequest.PhoneNumber = aws.String(phoneNumber)

		//create PhoneNumberValidateInput
		requestInput := &pinpoint.PhoneNumberValidateInput{}

		//assign request
		requestInput.NumberValidateRequest = numberValidateRequest

		//add request
		requests = append(requests, requestInput)
	}
	return
}
