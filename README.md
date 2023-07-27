# Tomba Email Finder Cli 🔥

CLI utility to search or verify email addresses in minutes.

## Features ✨

- 🛡️ Instantly locate email addresses from any website.
- 🛡️ Email verify to confirm an email address' authenticity.
- 🛡️ Enrich email with data.
- 🛡️ Instantly discover the email addresses of Linkedin URLs.
- 🛡️ Instantly discover the email addresses of article authors.

## Installation

### Using Snap

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/tomba)

```bash
sudo snap install tomba
```

### Using Go

Make sure that `$GOPATH/bin` is in your `$PATH`, because that's where this gets
installed:

```bash
go install github.com/tomba-io/tomba@latest
```

## Get Started 🎉

By default, invoking the CLI shows a help message:

```bash
tomba
```

![tomba email](svg/default.svg)

### Login

Sign in to Tomba account

```bash
tomba login
```

![tomba email](svg/login.svg)

### Domain search

Instantly locate email addresses from any company name or website.

```bash
tomba search --target "tomba.io"
```

### Enrichment

Locate and include data in your emails.

```bash
tomba enrich --target "b.mohamed@tomba.io"
```

![tomba enrich](svg/enrich.svg)

### Author Finder

Instantly discover the email addresses of article authors.

```bash
tomba author --target "https://clearbit.com/blog/company-name-to-domain-api"
```

### Linkedin Finder

Instantly discover the email addresses of Linkedin URLs.

```bash
tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia"
```

### Email Verifier

Verify the deliverability of an email address.

```bash
tomba verify --target "b.mohamed@tomba.io"
```

## Http

***Tomba Reverse Proxy***

```bash
tomba http
```

## Endpoints

| Name            | Route     | Query    | State     | Authentication | Method |
| --------------- | --------- | -------- | --------- | -------------- | ------ |
| Home            | /         | No       | Completed | No             | Get    |
| author finder   | /author   | `url`    | Completed | No             | Get    |
| email counter   | /count    | `domain` | Completed | No             | Get    |
| enrichment      | /enrich   | `email`  | Completed | No             | Get    |
| linkedin finder | /linkedin | `url`    | Completed | No             | Get    |
| domain search   | /search   | `domain` | Completed | No             | Get    |
| domain status   | /status   | `domain` | Completed | No             | Get    |
| email verifier  | /verify   | `email`  | Completed | No             | Get    |

### Available Commands

| Command name | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| author       | Instantly discover the email addresses of article authors.         |
| completion   | Generate the autocompletion script for the specified shell         |
| count        | Returns total email addresses we have for one domain.              |
| enrich       | Locate and include data in your emails.                            |
| help         | Help about any command                                             |
| http         | Runs a HTTP server (reverse proxy).                                |
| linkedin     | Instantly discover the email addresses of Linkedin URLs.           |
| login        | Sign in to Tomba account                                           |
| logout       | delete your current KEY & SECRET API session.                      |
| search       | Instantly locate email addresses from any company name or website. |
| status       | Returns domain status if is webmail or disposable.                 |
| verify       | Verify the deliverability of an email address.                     |
| version      | Print version number and build information.                        |

### Command Flags

| shortopts | longopts   | Description                                                        |
| --------- | ---------- | ------------------------------------------------------------------ |
| `-h`      | `--help`   | help for tomba                                                     |
| `-j`      | `--json`   | output JSON format. (default true)                                 |
| `-k`      | `--key`    | Tomba API KEY.                                                     |
| `-o`      | `--output` | Save the results to file.                                          |
| `-op`     | `--prot`   | Sets the port on which the HTTP server should bind. (default 3000) |
| `-s`      | `--secret` | Tomba API SECRET.                                                  |
| `-t`      | `--target` | TARGET SPECIFICATION Can pass email, Domain, URL, Linkedin URL.    |
| `-y`      | `--yaml`   | output YAML format.                                                |

## Auto-Completion

Auto-completion is supported for at least the following shells:

```bash
bash
fish
powershell
zsh
```

NOTE: it may work for other shells as well because the implementation is in
Golang and is not necessarily shell-specific.

### Completion

Installing auto-completions is as simple as running one command (works for
`bash`, `fish`, `powershell` and `zsh` shells):

```bash
tomba completion zsh
```

## Changelog 📌

Detailed changes for each release are documented in the [release notes](https://github.com/tomba-io/tomba/releases).

## Documentation

See the [official documentation](https://developer.tomba.io/).

### About Tomba

Founded in 2021, Tomba prides itself on being the most reliable, accurate, and in-depth source of Email address data available anywhere. We process terabytes of data to produce our Email finder API, company.

[![image](https://avatars.githubusercontent.com/u/67979591?s=200&v=4)](https://tomba.io/)

## Contribution

1. Fork it (<https://github.com/tomba-io/tomba/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

Please see the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html) file for more information.