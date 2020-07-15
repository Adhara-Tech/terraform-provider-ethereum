# Terraform provider for ethereum accounts

Terraform provider to generate ethereum accounts. A simple provider that allows you to create ethereum accounts and tracked with terraform. useful if you want to generate keys and store in vault or aws secrets manager.

## Usage

1. Build the provider `go build -o terraform-provider-ethereum`
2. Copy the binary to `$HOME/.terraform.d/plugins`
3. Copy the following terraform code to a `main.tf` file.
```hcl
provider "ethereum" {}

resource "ethereum_keystore_account" "key_test_empty_passphrase" {
 scrypt_encryption = "lightkdf"
}

resource "ethereum_keystore_account" "key_test_passphrase" {
  scrypt_encryption = "lightkdf"
  passphrase = "mySuperSecurePassphrase"
}
```
4. Execute `terraform plan` and `terraform apply`

## Importing accounts

You can import existing ethereum accounts with the private key and it's passphrase.

1. Add the resource that you want to import to your terraform code.
```hcl
resource "ethereum_keystore_account" "key_you_want_import" {
  scrypt_encryption = "lightkdf"
  passphrase = "mySuperSecurePassphrase"
}
```
2. Execute `terraform import ethereum_keystore_account.key_you_want_import yourPrivateKey_yourPassphrase`
3. You can also import a ethereum key with an empty passphrase. `terraform import ethereum_keystore_account.key_you_want_import yourPrivateKey_`, make sure to set the underscore after the private key.

## TODO

* [] Implement tests
