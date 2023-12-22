package prover

import (
	"encoding/json"
	"fmt"
	"os"
)

func CreateVkPk() {
	verificationKeyFile := "verificationKey.json"
	provingKeyFile := "provingKey.txt"

	if _, err := os.Stat(provingKeyFile); os.IsNotExist(err) {
		if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
			provingKey, verificationKey, error := GenerateVerificationKey()
			if error != nil {
				fmt.Println("Error generating verification key:", error)
			}
			vkJSON, _ := json.Marshal(verificationKey)
			vk_err := os.WriteFile(verificationKeyFile, vkJSON, 0644)
			if vk_err != nil {
				fmt.Println("Error writing verification key to file:", vk_err)
			}
			file, err := os.Create(provingKeyFile) // Creates a file named provingKeyData.txt
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close() // Make sure to close the file when done
			n, err := provingKey.WriteTo(file)
			if err != nil {
				fmt.Println("Error writing proving key to buffer:", err)
			}
			fmt.Printf("Wrote %d bytes to file\n", n)
		} else {
			return
		}
	}
	if _, err := os.Stat(verificationKeyFile); os.IsNotExist(err) {
		_, verificationKey, error := GenerateVerificationKey()
		if error != nil {
			fmt.Println("Error generating verification key:", error)
		}
		// Writing verification key to a file
		vkJSON, _ := json.Marshal(verificationKey)
		vk_err := os.WriteFile(verificationKeyFile, vkJSON, 0644)
		if vk_err != nil {
			fmt.Println("Error writing verification key to file:", vk_err)
		}
	} else {
		fmt.Println("Verification key already exists. No action needed.")
	}
	if _, err := os.Stat(provingKeyFile); os.IsNotExist(err) {
		provingKey, _, error := GenerateVerificationKey()
		if error != nil {
			fmt.Println("Error generating verification key:", error)
		}
		// Writing proving key to a file
		file, err := os.Create(provingKeyFile) // Creates a file named provingKeyData.txt
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close() // Make sure to close the file when done
		n, err := provingKey.WriteTo(file)
		if err != nil {
			fmt.Println("Error writing proving key to buffer:", err)
		}
		fmt.Printf("Wrote %d bytes to file\n", n)
	} else {
		fmt.Println("Proving key already exists. No action needed.")
	}
}
