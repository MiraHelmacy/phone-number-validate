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

	"phone-number-validate/cmd/validate"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	SHORT_FORMAT_KEY     = "short"
	ONLY_INVALID_NUMBERS = "onlyinvalid"
)

var rootCmd = &cobra.Command{
	Use:     "phone-number-validate [flags] phonenumber[...phonenumber]",
	Short:   "Validates phone numbers",
	Long:    `Validates a list of phone numbers.`,
	Args:    cobra.ArbitraryArgs,
	Version: "2.0.0",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		viper.BindPFlags(cmd.LocalFlags())
		viper.Set("phonenumbers", args)
		cliArgs := validate.ValidateCmdCliArgs{}
		viper.Unmarshal(&cliArgs)
		result, err := validate.ExecuteValidate(cliArgs)
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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool(SHORT_FORMAT_KEY, viper.GetBool(SHORT_FORMAT_KEY), "Toggle Short Format.")
	rootCmd.Flags().Bool(ONLY_INVALID_NUMBERS, viper.GetBool(ONLY_INVALID_NUMBERS), "Toggle only show invalid numbers")
}
