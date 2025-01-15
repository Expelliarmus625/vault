/*
Copyright © 2025 Ojas Barpande
Copyrights apply to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vault",
	Short: "A cli tool to encrypt and decrypt files",
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	if err := rootAction(); err != nil {
	// 		return err
	// 	} // call rootAction()
	// 	return nil
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vault.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// func rootAction(path, outputPath, key, ext *bool) error {
// 	switch *enc {
// 	case true:

// 	case false:
// 		encImage := newEncryptedImage(*path, *outputPath, *ext)
// 		err := encImage.Decrypt(*key)
// 		if err != nil {
// 			fmt.Fprintln(os.Stdout, err.Error())
// 			os.Exit(2)
// 		}

// 	default:
// 		fmt.Println("Please specify either -encrypt or -decrypt")
// 		os.Exit(1)
// 	}
// }
