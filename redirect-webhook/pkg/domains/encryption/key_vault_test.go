package encryption

import (
	"fmt"
	"testing"
)

func TestAesEncryptionAndDecryption(t *testing.T) {

	vault := KeyVault{}
	vault.CreateAesKey()

	testInputs := []string{"simplepass", "", "verylongvalueWithDistinctCharacters!@##$(*#&%(!*#%()*!&#@($#!(*&$(!*#ˆˆÄS:DLM:L M:AC<X>ROIEVNQPOCverylongvalueWithDistinctCharacters!@##$(*#&%(!*#%()*!&#@($#!(*&$(!*#ˆˆÄS:DLM:L M:AC<X>ROIEVNQPOCverylongvalueWithDistinctCharacters!@##$(*#&%(!*#%()*!&#@($#!(*&$(!*#ˆˆÄS:DLM:L M:AC<X>ROIEVNQPOCverylongvalueWithDistinctCharacters!@##$(*#&%(!*#%()*!&#@($#!(*&$(!*#ˆˆÄS:DLM:L M:AC<X>ROIEVNQPOC"}

	for _, textToEncrypt := range testInputs {

		bytes, e := vault.EncryptWithAes(textToEncrypt)
		if e != nil {
			t.Error(e)
		}

		result, e := vault.DecryptWithAes(bytes)
		if e != nil {
			t.Error(e)
		}

		if result != textToEncrypt {
			t.Errorf("values did not match.\n\tExpected: %s \n\tGot: %s", textToEncrypt, result)
		}
	}

}

func TestAesEncryptWithExternalData(t *testing.T) {
	vault := KeyVault{
		AESKey: []byte("c88ba867994f440963f55b727cdd3cb7"),
	}
	textToEncrypt := "Encrypt me"
	expectedOutput := "k49TM5LTflKyWXIo1BQ1Rg"

	hashedValue, err := vault.EncryptWithAes(textToEncrypt)
	if err != nil {
		t.Error(err)
	}

	if hashedValue != expectedOutput {
		t.Errorf("values did not match.\n\tExpected: %s \n\tGot: %s", expectedOutput, hashedValue)
	}
}

func TestAesDecryptWithExternalData(t *testing.T) {
	vault := KeyVault{
		AESKey: []byte("c88ba867994f440963f55b727cdd3cb7"),
	}
	hashToDecrypt := "/tsfpjgD5rj/CRhw5T0OCA=="
	expectedOutput := "Decrypted!!"

	result, err := vault.DecryptWithAes(hashToDecrypt)
	if err != nil {
		t.Error(err)
	}

	if result != expectedOutput {
		t.Errorf("values did not match.\n\tExpected: %s \n\tGot: %s", expectedOutput, result)
	}
}

func TestRsaEncryption(t *testing.T) {
	// private key used in test below:
	_ = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

	vault := KeyVault{
		IssuerPublicKey: []byte(`-----BEGIN PUBLIC KEY-----
MIICITANBgkqhkiG9w0BAQEFAAOCAg4AMIICCQKCAgB+abZRz3SFhzGSBemtlJac
+GmRU2jO+9iQFzXOE27k0/vgGHEcHAqd67w0MWC6Pm7jmDcIVL5JtTcobvYVaf3p
bSs+5Keh8UnC98ghIp325VWFln7B9PcoDZscWj+olRd4EiqhGCyfUq44O+/SXULu
pdnmeZVainBvGmoOfcy+Wuaq4DLmwWVMA/eaETGEgGxiYLMqwbDbcQlpRO0jKutb
lZnPuc7bvk86Bb8YyT8A/jCnB/fLVIcw4S8mhhtmE8NjCtMkZFNSUlGN6z1EBfAM
Cu2h/29V5o7D+jLdOyflX7plLgYXme9Pv6XE7ob6FEPHrW9u5VzCZ6zVJ1DsE+ec
zNPtsvRJXTspNOdzu9sgb/ffwmWZPx6AWZI+rBpAAGOxdv0Lbvmb9uXs4XZj54BF
M5XcpwEf43Ifn+9PbgUQMmehgv7yjY0wt9dfzyFsvjv4n5VVhl68+yVP9rpQwQV/
0+SqGj4urW8QehsvFXTGnyd1xOfkUNirIbyxGWzEyaLTmfpFwuezmDDICGE7Vyur
G/DRZBnHQOgS3iWI2r+laZWs/A/3mWAGEfPrkdmu+enabmzWjmDPi6DqElespn6/
fmZ9f03+yTuW5EutS/Re4uMFnxelYWIu4x1UUnHQqW6RFCMNixdh1bTej+pDOB80
GKiXQMAoRbSAK99DZbJkewIDAQAB
-----END PUBLIC KEY-----`),
	}

	textToEncrypt := "Encrypt me!"

	resultString, err := vault.EncryptWithPublicKey([]byte(textToEncrypt))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resultString)
}
