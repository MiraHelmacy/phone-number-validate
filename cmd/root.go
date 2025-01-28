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
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"phone-number-validate/cmd/config"
	"phone-number-validate/cmd/validate"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	SHORT_FORMAT_KEY     = "short"
	ONLY_INVALID_NUMBERS = "onlyinvalid"
	PHONE_NUMBERS_FILE   = "cfg"
	E164_REGEX_STRING    = `^\+[1-9]\d{1,14}$` //E.164 Regex as provided by Twilio. https://www.twilio.com/en-us/blog/validate-e164-phone-number-in-go
)

func e_164Regex() *regexp.Regexp {

	//attempt to compile regex for E.164 format.
	if regex, err := regexp.Compile(E164_REGEX_STRING); err == nil {

		//regex compiled successfully
		return regex
	} else {

		// regex failed to compile
		panic("E.164 regex did not compile: " + err.Error() + " Regex String: " + E164_REGEX_STRING)
	}

	//unreachable. regex returned or panic occurs
}

func cleanPhoneNumbersList(phoneNumbers []string) (cleanedPhoneNumbers []string) {

	//make new phone number list
	cleanedPhoneNumbers = make([]string, 0)

	//map to track duplicate numbers
	duplicatePhoneNumber := make(map[string]bool)

	//E.164 regex
	e_164_regex := e_164Regex()

	//for each phone number
	for _, phoneNumber := range phoneNumbers {
		//check the phone number is in E.164 format
		if match := e_164_regex.Match([]byte(phoneNumber)); match {
			//regex match.

			//check if the number is a duplicate
			if _, checked := duplicatePhoneNumber[phoneNumber]; !checked {
				//number is not a duplicate

				//add new phone number
				cleanedPhoneNumbers = append(cleanedPhoneNumbers, phoneNumber)
			} else {
				//number valid duplicate detected
				fmt.Println("Warning: duplicate phone number detected: ", phoneNumber, " Skipping.")
			}

		} else {
			//phone number not E.164 format

			//check if the phone number is a blank string
			if phoneNumber == "" {
				//phone number is blank. ignore
				continue
			}

			//no match and phone number is not blank
			fmt.Fprintln(os.Stderr, "Warning: non e.164 string detected: "+phoneNumber, "Ignoring")
		}

		//add phone number to duplicates
		duplicatePhoneNumber[phoneNumber] = true
	}
	return
}

func processFileContents(contents string) (phoneNumbers []string) {
	phoneNumbers = strings.Split(contents, "\n")
	return
}

func addPhoneNumbersFromConfigFile(originalPhoneNumbers []string) (originalAndPhoneNumbersFromConfig []string, err error) {

	//assign original to output
	originalAndPhoneNumbersFromConfig = originalPhoneNumbers

	//get the config file from viper and check if file name is provided
	if cfgFile := viper.GetString(PHONE_NUMBERS_FILE); cfgFile != "" {

		//clean file path
		cleanCfgFilePath := filepath.Clean(cfgFile)

		//read file contents
		if contents, readErr := config.ReadFileContents(cleanCfgFilePath); readErr == nil {

			//get list of phone numbers
			phoneNumbersFromFile := processFileContents(contents)

			//append processed file contents to output.
			originalAndPhoneNumbersFromConfig = append(originalAndPhoneNumbersFromConfig, phoneNumbersFromFile...)
		} else {
			//failed to read file
			err = readErr
			return
		}
	}

	return
}

var rootCmd = &cobra.Command{
	Use:     "phone-number-validate [flags] phonenumber[...phonenumber]",
	Short:   "Validates phone numbers",
	Long:    `Validates a list of phone numbers.`,
	Args:    cobra.ArbitraryArgs,
	Version: "3.0.4",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//bind local flags
		viper.BindPFlags(cmd.LocalFlags())

		//make empty list of phone numbers
		phoneNumbers := make([]string, 0)

		//append args to empty list
		phoneNumbers = append(phoneNumbers, args...)

		//add phone numbers from config file
		phoneNumbers, err = addPhoneNumbersFromConfigFile(phoneNumbers)

		//check adding phone numbers succeeded
		if err != nil {
			return
		}

		//clean phone number list of duplicates and invalid phone numbers
		phoneNumbers = cleanPhoneNumbersList(phoneNumbers)

		//no phone numbers to validate
		if len(phoneNumbers) == 0 {
			fmt.Println("No Phone Numbers to validate.")
			cmd.Usage()
			return
		}

		//set phone numbers in viper
		viper.Set("phonenumbers", phoneNumbers)

		//ValidateCmdCliArgs object
		cliArgs := validate.ValidateCmdCliArgs{}

		//bind viper config to validate args
		viper.Unmarshal(&cliArgs)

		//validate phone numbers
		result, err := validate.ExecuteValidate(cliArgs)

		//if no err
		if err == nil {

			//get the result and marshal error
			str, marshalErr := result.Result()

			//check for error
			if marshalErr != nil {
				return marshalErr
			}

			//print the result
			fmt.Println(str)
		}
		return
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool(SHORT_FORMAT_KEY, viper.GetBool(SHORT_FORMAT_KEY), "Toggle Short Format.")
	rootCmd.Flags().Bool(ONLY_INVALID_NUMBERS, viper.GetBool(ONLY_INVALID_NUMBERS), "Toggle only show invalid numbers")
	rootCmd.Flags().String(PHONE_NUMBERS_FILE, "", "config file containing a list of phone numbers")
}
