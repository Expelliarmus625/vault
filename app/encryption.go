package app

import (
	"crypto/rsa"
	"fmt"
	"path/filepath"
)


type Image struct {
	path string
	outputPath string
	ext string
	cipher []byte
}


func NewImage(path, outputPath string) Image {
	return Image{
		path: path,
		outputPath: outputPath,
		ext: filepath.Ext(path),
	}
}


func (i *Image) Encrypt(key string) error {
	//Parse image
	img, err := ParseImage(i.path)
	if err != nil {
		return err
	}

	var pkey *rsa.PrivateKey
	if key == "" {
		//Generate and save key pair
		pkey, err = GenerateRSAKeypair()
		if err != nil {
			return err
		}
		//Save private key
		savePrivateKey(pkey, "pkey.pem")
	} else {
		pkey, err = loadPrivateKey(key)
		if err != nil {
			return err
		}
	}

	i.cipher, err = EncryptImage(img, &pkey.PublicKey)
	if err != nil {
		return err
	}

	//Save encrypted image
	err = SaveEncryptedImage(*i)
	if err != nil {
		return err
	}
	fmt.Printf("\nImage encrypted and saved as %s\n", i.outputPath)

	return nil
}