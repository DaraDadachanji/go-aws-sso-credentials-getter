# Pcreds

When logging into the commandline using Control Tower, you typically need to
click the copy button, open your credentials file, paste the contents
in the appropriate location, and then potentially change the profile name.

pcreds is a simple commandline utility
to do all of these things with one command.

pcreds:

* Reads a credentials profile from the clipboard
* Parses the profile name and looks up an alias in pcreds.yaml (stored in your .aws folder)
* Parses your credentials file, updates the corresponding profile and saves an updated version

## Installation

Install Go from the [official website](https://go.dev/)

clone this repository and build the executable. Then move it to your bin folder

NOTE: if you are building for WSL you should include the `WSL_DISTRO_NAME` environment variable when building 

```bash
git clone https://github.com/DaraDadachanji/pcreds.git
cd pcreds
go mod tidy
go build
mv ./go-aws-sso-credentials /usr/local/bin/sso-cred
```

## Configuration

You may create a pcreds.yaml file in `~/.aws` to store aliases for your profiles

```
[profile my-profile]
sso_account_id = 123456654321
sso_role_name = my_aws_role_name
sso_start_url = https://my-organization.awsapps.com/start
sso_region = us-east-1
region = us-east-1

[profile my-second-profile]
sso_account_id = 654321123456
sso_role_name = my_aws_role_name
sso_start_url = https://my-organization.awsapps.com/start
sso_region = us-east-1
region = eu-west-2
```

The key should be the profile name generated by SSO while the value should be
your preferred alias for it.

This step is optional

## Usage

simply run `sso-cred {profile}` in your terminal
