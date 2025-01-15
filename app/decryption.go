package app

import (
	"fmt"
)

type EncryptedImage struct {
	path       string
	outputPath string
	ext        string
	plaintext  []byte
}

func NewEncryptedImage(path, outputPath, ext string) *EncryptedImage {
	if ext == "" {
		ext = ".png"
	}
	return &EncryptedImage{
		path:       path,
		outputPath: outputPath,
		ext:        ext,
	}
}

func (e *EncryptedImage) Decrypt(key string) error {
	fmt.Printf("Decrypting %s... \n", e.path)

	imgData, err := ParseEncryptedImage(e.path)
	if err != nil {
		fmt.Println("Error parsing encrypted image:", err)
		return err
	}

	pkey, err := loadPrivateKey(key)
	if err != nil {
		return err
	}
	e.plaintext, err = DecryptImage(imgData, pkey)
	if err != nil {
		fmt.Println("Error decrypting image:", err)
		return err
	}

	err = SaveDecryptedImage(*e)
	if err != nil {
		fmt.Println("Error saving decrypted image:", err)
		return err
	}
	fmt.Printf("\nImage decrypted and saved to %s\n", e.outputPath)
	return nil
}
