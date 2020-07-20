# Terraform provider for ethereum accounts

Terraform provider to generate ethereum accounts. A simple provider that allows you to create ethereum accounts and tracked with terraform. useful if you want to generate keys and store in vault or aws secrets manager.

## Usage

1. Build the provider `make build`
2. Copy the binary from `bin` to `$HOME/.terraform.d/plugins`, remove the archOS sufix.
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

### What if you don't have the private key ?

You can use a small tool inside scripts folder ([scripts/private_key.go](./scripts/private_key.go)) to get your private key from a kestore JSON account.

```bash
go run private_key.go <path-to-your-json-keystore> <passphrase>
```
The output is your private key. After this, you can follow the import instructions to import your key into terraform code.

## TODO

* [ ] Implement tests
