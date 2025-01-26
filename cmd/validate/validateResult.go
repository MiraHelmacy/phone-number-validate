package validate

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
)

// generate an AWS CLI command equivalent for this response
func awsCliCommand(response types.NumberValidateResponse) (cliString string) {
	cliString = fmt.Sprintf("aws pinpoint phone-number-validate --number-validate-request PhoneNumber=%v", *response.OriginalPhoneNumber)
	return
}

// Response from ExecuteValidate
type ValidateResult interface {
	Result() (result string, err error) //returns a string and an error
	Add(types.NumberValidateResponse)   //add a new NumberValidateResponse
}

// for a given validate result, marshal to JSON.
func marshalJson(v_result ValidateResult) (result string, err error) {
	//attempt to marshal to json
	bytes, marshalErr := json.MarshalIndent(v_result, "", "\t")

	//check error
	if marshalErr != nil {
		err = marshalErr
		return
	}

	//convert to string
	result = string(bytes)
	return
}

// long validate list item.
type LongValidateResultListItem struct {
	NumberValidateResponse types.NumberValidateResponse //the full response
	AwsCliCommand          string                       //the AWS CLI string
}

// A long validate response item
type LongValidateResult struct {
	Valid   *[]LongValidateResultListItem //list of valid phone numbers
	Invalid *[]LongValidateResultListItem //list of invalid phone numbers
}

// Implement Result
func (lvr LongValidateResult) Result() (result string, err error) {
	//marshal to json
	result, err = marshalJson(lvr)

	return
}

// implement Add
func (lvr LongValidateResult) Add(response types.NumberValidateResponse) {
	//create the AWS CLI string
	cliString := awsCliCommand(response)

	//create a new long result item
	listItem := LongValidateResultListItem{AwsCliCommand: cliString, NumberValidateResponse: response}

	//check if invalid
	if *response.PhoneType == "INVALID" {

		//number invalid
		*lvr.Invalid = append(*lvr.Invalid, listItem)
	} else {

		//number valid
		*lvr.Valid = append(*lvr.Valid, listItem)
	}
}

// Long Validate Result only invalid items
type LongValidateResultOnlyInvalid struct {
	Invalid *[]LongValidateResultListItem
}

func (lvroiv LongValidateResultOnlyInvalid) Result() (result string, err error) {
	result, err = marshalJson(lvroiv)

	return
}

// implement Add
func (lvroiv LongValidateResultOnlyInvalid) Add(response types.NumberValidateResponse) {
	//check if invalid
	if *response.PhoneType == "INVALID" {
		//create aws cli command
		cliString := awsCliCommand(response)

		//create list item
		listItem := LongValidateResultListItem{AwsCliCommand: cliString, NumberValidateResponse: response}

		//add to invalid list
		*lvroiv.Invalid = append(*lvroiv.Invalid, listItem)
	}
}

// short validate list item
type ShortValidateResultListItem struct {
	PhoneNumber   string //phone number
	AwsCliCommand string //cli command
}

// short validate result includes both valid and invalid
type ShortValidateResult struct {
	Valid   *[]ShortValidateResultListItem //valid items
	Invalid *[]ShortValidateResultListItem //invalid items
}

// Implement Result
func (svr ShortValidateResult) Result() (result string, err error) {

	//marshal to json
	result, err = marshalJson(svr)

	return
}

// Implement Add
func (svr ShortValidateResult) Add(response types.NumberValidateResponse) {

	//AWS CLI Command
	cliCommand := awsCliCommand(response)

	//new list item
	listItem := ShortValidateResultListItem{PhoneNumber: *response.OriginalPhoneNumber, AwsCliCommand: cliCommand}

	//check if invalid
	if *response.PhoneType == "INVALID" {

		//phone number invalid
		*svr.Invalid = append(*svr.Invalid, listItem)
	} else {

		//phone number valid
		*svr.Valid = append(*svr.Valid, listItem)
	}
}

// short Validate Result only invalid items
type ShortValidateResultOnlyInvalid struct {
	Invalid *[]ShortValidateResultListItem
}

// Implment Result
func (svroiv ShortValidateResultOnlyInvalid) Result() (result string, err error) {
	//marshal to json
	result, err = marshalJson(svroiv)

	return
}

// Implement Add
func (svrov ShortValidateResultOnlyInvalid) Add(response types.NumberValidateResponse) {

	//check if invalid
	if *response.PhoneType == "INVALID" {
		//AWS CLI command
		cliCommand := awsCliCommand(response)

		//list item
		listItem := ShortValidateResultListItem{PhoneNumber: *response.OriginalPhoneNumber, AwsCliCommand: cliCommand}

		//add to invalid
		*svrov.Invalid = append(*svrov.Invalid, listItem)
	}
}

// create a new short validate result
func NewShortValidateResult(onlyInvalid bool) (newResult ValidateResult) {

	//create a valid list
	valid := make([]ShortValidateResultListItem, 0)

	//create an invalid list
	invalid := make([]ShortValidateResultListItem, 0)

	//create a standard validate result
	newResult = ShortValidateResult{Valid: &valid, Invalid: &invalid}

	//if only invalid parameter detected
	if onlyInvalid {
		//create a new result with only invalid items
		newResult = ShortValidateResultOnlyInvalid{Invalid: &invalid}
	}
	return
}

// new long validate result
func NewLongValidateResult(onlyInvalid bool) (newResult ValidateResult) {

	//valid list
	valid := make([]LongValidateResultListItem, 0)

	//invalid list
	invalid := make([]LongValidateResultListItem, 0)

	//create standard long validate result
	newResult = LongValidateResult{Valid: &valid, Invalid: &invalid}

	//if only invalid requested
	if onlyInvalid {
		//recreate a new result with only invalid options
		newResult = LongValidateResultOnlyInvalid{Invalid: &invalid}
	}
	return
}
