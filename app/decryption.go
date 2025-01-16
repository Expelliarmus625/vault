package app

import (
	"fmt"
	"sync"
)

type EncryptedImage struct {
	path       string
	outputPath string
	ext        string
	plaintext  []byte
}

type EncryptedImageList []EncryptedImage

func NewEncryptedImage(path, outputPath, ext string) EncryptedImage {
	if ext == "" {
		ext = ".png"
	}
	return EncryptedImage{
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


func (e *EncryptedImageList) DecryptAll(key string) error {
	var wg sync.WaitGroup
	wg.Add(len(*e))

	for _, image := range *e {
		go func(image EncryptedImage) {
			defer wg.Done()
			err := image.Decrypt(key)
			if err != nil {
				fmt.Println(err)
			}
		}(image)
	}

	wg.Wait()
	return nil
}