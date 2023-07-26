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

### Installation

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