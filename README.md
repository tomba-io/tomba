# Tomba Email Finder Cli 🔥

CLI utility to search or verify email addresses in minutes.

## Features ✨

- 🛡️ Instantly locate email addresses from any website.
- 🛡️ Email verify to confirm an email address' authenticity.
- 🛡️ Enrich email with data.
- 🛡️ Instantly discover the email addresses of Linkedin URLs.
- 🛡️ Instantly discover the email addresses of article authors.
- 🛡️ Search for phone numbers given an email, domain, or LinkedIn URL.
- 🛡️ Validate phone numbers.
- 🛡️ Retrieve domains similar to a specific domain.
- 🛡️ Discover technologies detected for a domain.
- 🛡️ Search for companies using natural language or structured filters.
- 🛡️ Print current account information.

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

### Using homebrew tap

[The formula](https://github.com/tomba-io/homebrew-tap/blob/master/Formula/tomba.rb)

```bash
brew install tomba-io/tap/tomba
```

### Using scoop

```bash
scoop bucket add tomba https://github.com/tomba-io/scoop-bucket.git
scoop install tomba
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

Slack Command

```bash
/search tomba.io
```

### Email Finder

Retrieves the most likely email address from a domain name, a first name and a last name.

```bash
tomba finder --target "tomba.io" --fist "mohamed" --last "ben rebia"
```

With phone number enrichment:

```bash
tomba finder --target "tomba.io" --fist "mohamed" --last "ben rebia" --enrich-mobile
```

### Enrichment

Locate and include data in your emails.

```bash
tomba enrich --target "b.mohamed@tomba.io"
```

With phone number enrichment:

```bash
tomba enrich --target "b.mohamed@tomba.io" --enrich-mobile
```

Slack Command

```bash
/enrich b.mohamed@tomba.io
```

![tomba enrich](svg/enrich.svg)

### Author Finder

Instantly discover the email addresses of article authors.

```bash
tomba author --target "https://clearbit.com/blog/company-name-to-domain-api"
```

Slack Command

```bash
/author https://clearbit.com/blog/company-name-to-domain-api
```

### Linkedin Finder

Instantly discover the email addresses of Linkedin URLs.

```bash
tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia"
```

With phone number enrichment:

```bash
tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia" --enrich-mobile
```

### Email Sources

Find email address source somewhere on the web

```bash
tomba sources --target "ab@tomba.io"
```

Slack Command

```bash
/linkedin https://www.linkedin.com/in/mohamed-ben-rebia
```

### Phone Finder

Search for phone numbers given an email, domain, or LinkedIn URL.

```bash
tomba phone-finder --email "info@stripe.com"
tomba phone-finder --domain "tomba.io"
tomba phone-finder --linkedin "https://www.linkedin.com/in/alex-maccaw-ab592978"
tomba phone-finder --domain "stripe.com" --full
```

### Phone Validator

Validate phone numbers.

```bash
tomba phone-validator --phone "+14155552671"
tomba phone-validator --phone "4155552671" --country-code US
```

### Reveal (Company Search)

Search for companies using natural language or structured filters.

```bash
tomba reveal --query "Real Estate in France"
tomba reveal --country US,UK --industry Technology
tomba reveal --country US --size 101-500,501-1000 --page 2
```

### Similar Domains

Retrieve domains similar to a specific domain.

```bash
tomba similar --target "tomba.io"
```

### Technology

Discover technologies detected for a domain.

```bash
tomba technology --target "tomba.io"
```

### Whoami

Print current account information.

```bash
tomba whoami
```

### Email Verifier

Verify the deliverability of an email address.

```bash
tomba verify --target "b.mohamed@tomba.io"
```

Slack Command

```bash
/checker b.mohamed@tomba.io
```

## Http

**_Tomba Reverse Proxy_**

```bash
tomba http
```

## Endpoints

| Name            | Route            | Body                                                 | State     | Slack | Method |
| --------------- | ---------------- | ---------------------------------------------------- | --------- | ----- | ------ |
| author finder   | /author          | `url`                                                | Completed | Yes   | Post   |
| email counter   | /count           | `domain`                                             | Completed | No    | Post   |
| enrichment      | /enrich          | `email`, `enrich_mobile`                             | Completed | Yes   | Post   |
| email finder    | /finder          | `domain`, `first_name`, `last_name`, `enrich_mobile` | Completed | No    | Post   |
| linkedin finder | /linkedin        | `url`, `enrich_mobile`                               | Completed | Yes   | Post   |
| phone finder    | /phone-finder    | `email`, `domain`, or `linkedin`                     | Completed | No    | Post   |
| phone validator | /phone-validator | `phone`, `country_code`                              | Completed | No    | Post   |
| reveal          | /reveal          | `query`, `country`, `industry`, `size`               | Completed | No    | Post   |
| domain search   | /search          | `domain`                                             | Completed | Yes   | Post   |
| similar domains | /similar         | `domain`                                             | Completed | No    | Post   |
| email sources   | /sources         | `email`                                              | Completed | No    | Post   |
| domain status   | /status          | `domain`                                             | Completed | No    | Post   |
| technology      | /technology      | `domain`                                             | Completed | No    | Post   |
| email verifier  | /verify          | `email`                                              | Completed | Yes   | Post   |
| logs            | /logs            | No                                                   | Completed | No    | Get    |
| usage           | /usage           | No                                                   | Completed | No    | Get    |
| whoami          | /whoami          | No                                                   | Completed | No    | Get    |

### Available Commands

| Command name    | Description                                                                               |
| --------------- | ----------------------------------------------------------------------------------------- |
| author          | Instantly discover the email addresses of article authors.                                |
| completion      | Generate the autocompletion script for the specified shell                                |
| count           | Returns total email addresses we have for one domain.                                     |
| enrich          | Locate and include data in your emails.                                                   |
| finder          | Retrieves the most likely email address from a domain name, a first name and a last name. |
| help            | Help about any command                                                                    |
| http            | Runs a HTTP server (reverse proxy).                                                       |
| linkedin        | Instantly discover the email addresses of Linkedin URLs.                                  |
| login           | Sign in to Tomba account                                                                  |
| logout          | delete your current KEY & SECRET API session.                                             |
| logs            | Check your last 1,000 requests you made during the last 3 months.                         |
| phone-finder    | Search for phone numbers given an email, domain, or LinkedIn URL.                         |
| phone-validator | Validate phone numbers.                                                                   |
| reveal          | Search for companies using natural language or structured filters.                        |
| search          | Instantly locate email addresses from any company name or website.                        |
| similar         | Retrieve domains similar to a specific domain.                                            |
| status          | Returns domain status if is webmail or disposable.                                        |
| technology      | Discover technologies detected for a domain.                                              |
| usage           | Check your monthly requests.                                                              |
| verify          | Verify the deliverability of an email address.                                            |
| version         | Print version number and build information.                                               |
| whoami          | Print current account information.                                                        |

### Command Global Flags

| shortopts | longopts   | Description                                                        |
| --------- | ---------- | ------------------------------------------------------------------ |
| `-h`      | `--help`   | help for tomba                                                     |
| `-j`      | `--json`   | output JSON format. (default true)                                 |
| `-k`      | `--key`    | Tomba API KEY.                                                     |
| `-o`      | `--output` | Save the results to file.                                          |
| `-p`      | `--prot`   | Sets the port on which the HTTP server should bind. (default 3000) |
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

See the [official documentation](https://docs.tomba.io/).

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
