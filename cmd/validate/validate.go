/*
Copyright © 2025 Alex Helmacy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package validate

import (
	"fmt"
)

func ExecuteValidate(args ValidateCmdCliArgs) (result ValidateResult, err error) {
	//validate args
	if valid := args.Validate(); !valid {

		err = fmt.Errorf("invalid args: %v", args)
		return
	}

	//create a pinpoint action wrapper
	wrapper, createWrapperError := NewPinpointWrapper()

	//check for errors
	if createWrapperError != nil {
		err = createWrapperError
		return
	}

	//create result. initially long validate result.
	result = NewLongValidateResult(args.OnlyInvalid)

	//if short option specified
	if args.Short {

		//create a short validate result object
		result = NewShortValidateResult(args.OnlyInvalid)
	}

	//for each Number Validate Request
	for _, request := range args.PhoneNumberValidateRequests() {

		//validate the request
		phoneNumberValidateResponse, phoneNumberValidateRequestErr := wrapper.PhoneNumberValidate(request)

		//check for errors
		if phoneNumberValidateRequestErr != nil {
			err = phoneNumberValidateRequestErr
			return
		}

		//get the number validate response
		numberValidateResponse := *phoneNumberValidateResponse.NumberValidateResponse

		//add the response to the result
		result.Add(numberValidateResponse)
	}

	return
}
