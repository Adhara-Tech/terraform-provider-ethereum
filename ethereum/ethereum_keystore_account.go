package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/Adhara-Tech/terraform-provider-ethereum/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEthereumKeystore() *schema.Resource {

	return &schema.Resource{
		Create: CreateAccount,
		Update: UpdateAccount,
		Delete: DeleteAccount,
		Read:   ReadAccount,

		Schema: map[string]*schema.Schema{
			"scrypt_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "scrypt encryption difficulty ( lightkdf => 4MB and 100ms CPU approximately, standartkdf => 256MB and 1s CPU approximately",
				Default:     "lightkdf",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "passphrase to encrypt the ethereum account",
				Default:     "",
			},
			"public_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"encrypted_json": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func getScryptValues(scryptEncryption string) (scriptN, scriptP int, err error) {

	switch scryptEncryption {

	case "lightkdf":
		scriptN = 4096
		scriptP = 6
		err = nil
		return
	case "standartkdf":
		scriptN = 262144
		scriptP = 1
		err = nil
		return
	default:
		return 0, 0, fmt.Errorf("invalid script encryption. Should be lightkdf or standartkdf")
	}
}

func CreateAccount(d *schema.ResourceData, meta interface{}) error {

	scryptEncryption := d.Get("scrypt_encryption").(string)
	passphrase := d.Get("passphrase").(string)

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("error generating ecdsa key: %s", err)
	}

	scriptN, scriptP, err := getScryptValues(scryptEncryption)
	if err != nil {
		return err
	}

	crytoJSON, err := utils.EncryptJSONV3(k.D.Bytes(), []byte(passphrase), scriptN, scriptP)
	if err != nil {
		return err
	}

	privateKey := hex.EncodeToString(k.D.Bytes())
	publicAddress := hex.EncodeToString(utils.PublicKeyToAddress(k.PublicKey))
	uuid := uuid.New()

	keystoreJSONV3 := utils.EncryptedKeyV3{
		Address: publicAddress,
		Crypto:  crytoJSON,
		Id:      uuid.String(),
		Version: 3,
	}

	encryptedJSON, err := json.Marshal(keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error parsing EncryptedKeyV3 json")
	}

	d.SetId(publicAddress)
	d.Set("private_key", privateKey)
	d.Set("public_address", publicAddress)
	d.Set("encrypted_json", string(encryptedJSON))

	return nil
}

func UpdateAccount(d *schema.ResourceData, meta interface{}) error {

	scryptEncryption := d.Get("scrypt_encryption").(string)
	passphrase := d.Get("passphrase").(string)

	scriptN, scriptP, err := getScryptValues(scryptEncryption)
	if err != nil {
		return err
	}
	privateKeyBytes, err := hex.DecodeString(d.Get("private_key").(string))
	if err != nil {
		return fmt.Errorf("error decoding private key")
	}

	var keystoreJSONV3 utils.EncryptedKeyV3
	err = json.Unmarshal([]byte(d.Get("encrypted_json").(string)), &keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error unmarshaling the encryptedKeyV3")
	}

	crytoJSON, err := utils.EncryptJSONV3(privateKeyBytes, []byte(passphrase), scriptN, scriptP)
	if err != nil {
		return err
	}

	keystoreJSONV3.Crypto = crytoJSON

	encryptedJSON, err := json.Marshal(keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error parsing EncryptedKeyV3 json")
	}

	d.Set("encrypted_json", string(encryptedJSON))

	return nil
}

func DeleteAccount(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ReadAccount(d *schema.ResourceData, meta interface{}) error {
	return nil
}
