package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var fileName string

var generateKeyCmd = &cobra.Command{
	Use:   "generate-key <directory>",
	Short: "Generates a pub/priv keypair to specified location",
	Long: `Generates a pub/priv keypair to specified location otherwise defaults to current directory.
	You can also specify a name for it using the -n flag`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := ""
		if len(args) > 0 {
			path = args[0]
			path = fixFilePath(path)
		}
		err := generateKey(path, fileName)
		if err != nil {
			fmt.Println("Failed to update config file")
			os.Exit(1)
		}
		fmt.Println("SSH keys have been generated and the config file has been updated")
	},
}

func init() {
	rootCmd.AddCommand(generateKeyCmd)
	generateKeyCmd.Flags().StringVarP(&fileName, "filename", "n", "holepunch_key", "The name your new key files will have")
}

func generateKey(keyPath string, fileName string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	if keyPath == "" {
		ex, errEx := os.Executable()
		if errEx != nil {
			return errEx
		}
		keyPath = filepath.Dir(ex) + string(os.PathSeparator)
	}
	if !strings.HasSuffix(keyPath, string(os.PathSeparator)) {
		keyPath += string(os.PathSeparator)
	}
	// generate and write private key as PEM
	privateKeyFile, err := os.Create(keyPath + fileName + ".pem")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if errEncode := pem.Encode(privateKeyFile, privateKeyPEM); errEncode != nil {
		return errEncode
	}
	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(keyPath+fileName+".pub", ssh.MarshalAuthorizedKey(pub), 0655)
	if err != nil {
		return err
	}
	viper.Set("privatekeypath", keyPath+fileName+".pem")
	viper.Set("publickeypath", keyPath+fileName+".pub")

	return viper.WriteConfig()
}
