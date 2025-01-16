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

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt filePath outputPath",
	Aliases: []string{"d"},
	Short: "Used for decrypting files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// path, _ := cmd.Flags().GetString("path")
		// outputPath, _ := cmd.Flags().GetString("outputPath")
		path := args[0]
		outputPath := args[1]
		key, _ := cmd.Flags().GetString("key")
		ext, _ := cmd.Flags().GetString("ext")
		if err := decryptAction(path, outputPath, key, ext); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// decryptCmd.Flags().StringP("path", "p", "", "Path to image")
	// decryptCmd.Flags().StringP("outputPath", "o", "", "Path to output image")
	decryptCmd.Flags().StringP("key", "k", "./privateKey.pem", "Path to private key")
	decryptCmd.Flags().StringP("ext", "e", "", "File extension")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func decryptAction(path, outputPath, key, ext string) error {
	
//Get file info to determine if it's a file or directory
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	//initialize empty image list
	images := &app.EncryptedImageList{}

	//if it's a directory, decrypt all files
	if info.IsDir() {
		filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				targetPath, err := filepath.Rel(path, path)
				if err != nil {
					return err
				}
				image := app.NewEncryptedImage(path, filepath.Join(outputPath, filepath.Dir(targetPath)), ext)
				*images = append(*images, image)
			}
			return nil
		})
	}
	images.DecryptAll(key)

	encImage := app.NewEncryptedImage(path, outputPath, ext)
	err = encImage.Decrypt(key)
	if err != nil {
		return err
	}
	return nil
}