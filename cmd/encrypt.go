/*
Copyright © 2025 Ojas Barpande
Copyrights apply to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/expelliarmus625/vault/app"
	"github.com/spf13/cobra"
)

// encCmd represents the enc command
var encryptCmd = &cobra.Command{
	Use:     "encrypt",
	Aliases: []string{"e"},
	Short:   "Used for encrypting files",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// path, _ := cmd.Flags().GetString("path")
		path := args[0]
		// outputPath, _ := cmd.Flags().GetString("outputPath")
		outputPath := args[1]
		key, _ := cmd.Flags().GetString("key")
		if err := encryptAction(path, outputPath, key); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encCmd.PersistentFlags().String("foo", "", "A help for foo")

	// encryptCmd.Flags().StringP("path", "p", "", "Path to image")
	// encryptCmd.Flags().StringP("outputPath", "o", "", "Path to output image")
	encryptCmd.Flags().StringP("key", "k", "./privateKey.pem", "Path to private key")
}

func encryptAction(inputpath, outputPath, key string) error {

	//Get file info to determine if it is a file or directory
	info, err := os.Stat(inputpath)
	if err != nil {
		return err
	}

	images := &app.ImageList{}
	if info.IsDir() {
		filepath.WalkDir(inputpath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				targetPath,err := filepath.Rel(inputpath, path)
				if err != nil {
					return err
				}
				image := app.NewImage(path, filepath.Join(outputPath, filepath.Dir(targetPath)))
				// err = image.Encrypt(key)
				*images = append(*images, image)
				if err != nil {
					return err
				}
			}
			return nil
		})
		images.EncryptAll(key)

		return nil
	}

	image := app.NewImage(inputpath, outputPath)
	err = image.Encrypt(key)
	if err != nil {
		return err
	}
	return nil
}
