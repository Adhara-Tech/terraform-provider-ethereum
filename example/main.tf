provider "ethereum" {}

provider "random" {
  version = "2.3.0"
}

resource "random_password" "passphrase" {
  length = 16
  special = false
}

resource "ethereum_keystore_account" "example_account" {
 scrypt_encryption = "standartkdf"
 passphrase = random_password.passphrase.result
}
