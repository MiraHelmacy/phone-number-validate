package validate

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
)

func awsCliCommand(response types.NumberValidateResponse) (cliString string) {
	cliString = fmt.Sprintf("aws pinpoint phone-number-validate --number-validate-request PhoneNumber=%v", *response.OriginalPhoneNumber)
	return
}

type ValidateResult interface {
	Result() (result string, err error)
	Add(types.NumberValidateResponse)
}

func marshalJson(v_result ValidateResult) (result string, err error) {
	bytes, marshalErr := json.MarshalIndent(v_result, "", "\t")
	if marshalErr != nil {
		err = marshalErr
		return
	}

	result = string(bytes)
	return
}

type LongValidateResultListItem struct {
	NumberValidateResponse types.NumberValidateResponse
	AwsCliCommand          string
}

type LongValidateResult struct {
	Valid   *[]LongValidateResultListItem
	Invalid *[]LongValidateResultListItem
}

func (lvr LongValidateResult) Result() (result string, err error) {
	result, err = marshalJson(lvr)

	return
}

func (lvr LongValidateResult) Add(response types.NumberValidateResponse) {
	cliString := awsCliCommand(response)
	listItem := LongValidateResultListItem{AwsCliCommand: cliString, NumberValidateResponse: response}
	if *response.PhoneType == "INVALID" {
		*lvr.Invalid = append(*lvr.Invalid, listItem)
	} else {
		*lvr.Valid = append(*lvr.Valid, listItem)
	}
}

type LongValidateResultOnlyInvalid struct {
	Invalid *[]LongValidateResultListItem
}

func (lvrov LongValidateResultOnlyInvalid) Result() (result string, err error) {
	result, err = marshalJson(lvrov)

	return
}

func (lvrov LongValidateResultOnlyInvalid) Add(response types.NumberValidateResponse) {

	if *response.PhoneType == "INVALID" {
		cliString := awsCliCommand(response)
		listItem := LongValidateResultListItem{AwsCliCommand: cliString, NumberValidateResponse: response}
		*lvrov.Invalid = append(*lvrov.Invalid, listItem)
	}
}

type ShortValidateResultListItem struct {
	PhoneNumber   string
	AwsCliCommand string
}

type ShortValidateResult struct {
	Valid   *[]ShortValidateResultListItem
	Invalid *[]ShortValidateResultListItem
}

func (svr ShortValidateResult) Result() (result string, err error) {

	result, err = marshalJson(svr)

	return
}

func (svr ShortValidateResult) Add(response types.NumberValidateResponse) {
	cliCommand := awsCliCommand(response)
	listItem := ShortValidateResultListItem{PhoneNumber: *response.OriginalPhoneNumber, AwsCliCommand: cliCommand}
	if *response.PhoneType == "INVALID" {
		*svr.Invalid = append(*svr.Invalid, listItem)
	} else {
		*svr.Valid = append(*svr.Valid, listItem)
	}
}

type ShortValidateResultOnlyInvalid struct {
	Invalid *[]ShortValidateResultListItem
}

func (svrov ShortValidateResultOnlyInvalid) Result() (result string, err error) {
	result, err = marshalJson(svrov)

	return
}

func (svrov ShortValidateResultOnlyInvalid) Add(response types.NumberValidateResponse) {
	if *response.PhoneType == "INVALID" {
		cliCommand := awsCliCommand(response)
		listItem := ShortValidateResultListItem{PhoneNumber: *response.OriginalPhoneNumber, AwsCliCommand: cliCommand}
		*svrov.Invalid = append(*svrov.Invalid, listItem)
	}
}

func NewShortValidateResult(onlyInvalid bool) (newResult ValidateResult) {
	valid := make([]ShortValidateResultListItem, 0)
	invalid := make([]ShortValidateResultListItem, 0)
	newResult = ShortValidateResult{Valid: &valid, Invalid: &invalid}
	if onlyInvalid {
		newResult = ShortValidateResultOnlyInvalid{Invalid: &invalid}
	}
	return newResult
}

func NewLongValidateResult(onlyInvalid bool) (newResult ValidateResult) {
	valid := make([]LongValidateResultListItem, 0)
	invalid := make([]LongValidateResultListItem, 0)
	newResult = LongValidateResult{Valid: &valid, Invalid: &invalid}
	if onlyInvalid {
		newResult = LongValidateResultOnlyInvalid{Invalid: &invalid}
	}
	return newResult
}
