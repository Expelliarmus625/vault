package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"path/filepath"
)

// Parse Image
func ParseImage(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("input file path is empty")
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Parse Encrypted Image
func ParseEncryptedImage(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("input encrypted file path is empty")
	}
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Generate Key Pair
func GenerateRSAKeypair() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Encrypt Image
func EncryptImage(image []byte, pub *rsa.PublicKey) ([]byte, error) {
	chunkSize := pub.Size() - 11

	// bar := progressbar.Default(int64(len(image)), "Encrypting image...")
	bar := progressbar.NewOptions(len(image),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan]Encrypting file...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	var encryptedImage []byte
	for i := 0; i < len(image); i += chunkSize {
		end := i + chunkSize
		if end > len(image) {
			end = len(image)
		}
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, pub, image[i:end])
		if err != nil {
			return nil, err
		}
		encryptedImage = append(encryptedImage, encryptedChunk...)
		bar.Add(chunkSize)
	}
	return encryptedImage, nil
}

// Decrypt Image
func DecryptImage(encryptedImage []byte, priv *rsa.PrivateKey) ([]byte, error) {
	chunkSize := priv.Size()
	var decryptedImage []byte

	// bar := progressbar.Default(int64(len(encryptedImage)), "Decrypting image...")
	bar := progressbar.NewOptions(len(encryptedImage),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan]Decrypting file...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for i := 0; i < len(encryptedImage); i += chunkSize {
		end := i + chunkSize
		if end > len(encryptedImage) {
			end = len(encryptedImage)
		}
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, priv, encryptedImage[i:end])
		if err != nil {
			return nil, err
		}
		decryptedImage = append(decryptedImage, decryptedChunk...)
		bar.Add(chunkSize)
	}
	return decryptedImage, nil
}

// Save Image
func SaveEncryptedImage(image Image) error {
	outPath := filepath.Join(image.outputPath, filepath.Base(image.path))

	// Create the outputPath if it doesn't exist
	if _, err := os.Stat(image.outputPath); os.IsNotExist(err) {
		err = os.MkdirAll(image.outputPath, 0755)
		if err != nil {
			return err
		}
	}
	err := os.WriteFile(outPath, image.cipher, 0644)
	if err != nil {
		return err
	}
	return nil
}

// SaveDecryptedImage will save the decrypted image to the output path 
func SaveDecryptedImage(e EncryptedImage) error {
	// filename := strings.TrimSuffix(filepath.Base(e.path), filepath.Ext(e.path))
	outfile := filepath.Join(e.outputPath, filepath.Base(e.path))

	// Create the outputPath if it doesn't exist
	if _, err := os.Stat(e.outputPath); os.IsNotExist(err) {
		err = os.MkdirAll(e.outputPath, 0755)
		if err != nil {
			return err
		}
	}

	err := os.WriteFile(outfile, e.plaintext, 0644)
	if err != nil {
		return err
	}
	return nil
}
func savePrivateKey(privateKey *rsa.PrivateKey, filename string) {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err := os.WriteFile(filename, privateKeyPEM, 0600)
	if err != nil {
		fmt.Println("Error saving private key:", err)
		return
	}
	fmt.Printf("Private key saved to : %s", filename)
}

func loadPrivateKey(filename string) (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error loading private key:", err)
		return nil, err
	}
	block, _ := pem.Decode(keyData)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing private key:", err)
		return nil, err
	}
	return privateKey, nil
}
