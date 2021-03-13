# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

For each changelong entry sections are:

* `Added` for new features.
* `Changed` for changes in existing functionality.
* `Deprecated` for soon-to-be removed features.
* `Removed` for now removed features.
* `Fixed` for any bug fixes.
* `Security` in case of vulnerabilities.

## [Unreleased]

## [v0.0.2] - 2021/03/10

Switch to github actions and goreleaser to be able to push to terrafrom registry

## [0.0.1] - 2020/09/16

First stable version of terraform-ethereum-provider

## [0.0.1-rc3] - 2020/07/30

### Changed
2. change paremeter `public_address` to `address`

### Fixed
1. use the correct elliptic curve secp256k1

## [0.0.1-rc2] - 2020/07/21

### Added
1. Decrypt script to get the private key from a keystore account

### Fixed
1. Passphrase parametes was in plain text. Update with sensitive=true.
2. Fix issue with passphrase when importing a key
3. compile statically to fix issues with alpine images

## [0.0.1-rc1] - 2020/07/20
First release candidate of terraform-provider-ethereum

### Added
1. Create ethereum accounts
2. Update ethereum accounts
3. Import ethereum accounts
