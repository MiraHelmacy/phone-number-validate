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
	if valid := args.Validate(); !valid {

		err = fmt.Errorf("invalid args: %v", args)
		return
	}

	wrapper, createWrapperError := NewPinpointWrapper()
	if createWrapperError != nil {
		err = createWrapperError
		return
	}
	result = NewLongValidateResult(args.OnlyInvalid)
	if args.Short {
		result = NewShortValidateResult(args.OnlyInvalid)
	}

	for _, request := range args.PhoneNumberValidateRequests() {
		phoneNumberValidateResponse, phoneNumberValidateRequestErr := wrapper.PhoneNumberValidate(request)
		if phoneNumberValidateRequestErr != nil {
			err = phoneNumberValidateRequestErr
			return
		}

		numberValidateResponse := *phoneNumberValidateResponse.NumberValidateResponse
		result.Add(numberValidateResponse)
	}

	return
}

/*
var ValidateCmd = &cobra.Command{
	Use:   "validate [flags] phonenumber[...phonenumber]",
	Short: "Validates phone numbers",
	Long:  `Validates a list of phone numbers.`,
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.LocalFlags())
		viper.Set("phonenumbers", args)
		cliArgs := ValidateCmdCliArgs{}
		viper.Unmarshal(&cliArgs)
		result, err := ExecuteValidate(cliArgs)
		if err == nil {
			str, marshalErr := result.Result()
			if marshalErr != nil {
				return marshalErr
			}
			fmt.Println(str)
		}
		return
	},
}

func init() {
	ValidateCmd.Flags().Bool(config.SHORT_FORMAT_KEY, viper.GetBool(config.SHORT_FORMAT_KEY), "Toggle Short Format.")
	ValidateCmd.Flags().Bool(config.ONLY_INVALID_NUMBERS, viper.GetBool(config.ONLY_INVALID_NUMBERS), "Toggle only show invalid numbers")
}
*/
