<p align="center"><img src="https://tomba.io/logo.svg" height='60' /></p>

<h3 align="center">Tomba Email Finder.</h3>

# Tomba Email Finder Cli

CLI utility to search or verify email addresses in minutes.

## Installation

### Using Snap

```bash
sudo snap install email
```

### Using Go

```bash
go install github.com/tomba-io/email@latest
```

## Get Started 🎉

```bash
cli utility to search or verify lists of email addresses in minutes can significantly improve productivity and efficiency for individuals and businesses dealing with large email databases.


    ████████╗ ██████╗ ███╗   ███╗██████╗  █████╗    ██╗ ██████╗ 
    ╚══██╔══╝██╔═══██╗████╗ ████║██╔══██╗██╔══██╗   ██║██╔═══██╗
       ██║   ██║   ██║██╔████╔██║██████╔╝███████║   ██║██║   ██║
       ██║   ██║   ██║██║╚██╔╝██║██╔══██╗██╔══██║   ██║██║   ██║
       ██║   ╚██████╔╝██║ ╚═╝ ██║██████╔╝██║  ██║██╗██║╚██████╔╝
       ╚═╝    ╚═════╝ ╚═╝     ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝╚═╝ ╚═════╝ 
       
    

Usage:
  email [command]

Available Commands:
  author      Instantly discover the email addresses of article authors.
  completion  Generate the autocompletion script for the specified shell
  count       Returns total email addresses we have for one domain.
  enrich      Locate and include data in your emails.
  help        Help about any command
  linkedin    Instantly discover the email addresses of Linkedin URLs.
  login       Sign in to Tomba account
  logout      delete your current KEY & SECRET API session.
  search      Instantly locate email addresses from any company name or website.
  status      Returns domain status if is webmail or disposable.
  verify      Verify the deliverability of an email address.
  version     Print version number and build information.

Flags:
  -N, --color           disable color output. (default true)
  -c, --csv             output CSV format.
  -h, --help            help for email
  -j, --json            output JSON format. (default true)
  -k, --key string      Tomba API KEY.
  -o, --output string   Save the results to file.
  -p, --pretty          output pretty format. (default true)
  -s, --secret string   Tomba API SECRET.
  -t, --target string   TARGET SPECIFICATION Can pass email, Domain, URL, Linkedin URL.
  -y, --yaml            output YAML format.

Use "email [command] --help" for more information about a command.
```

## Changelog 📌

Detailed changes for each release are documented in the [release notes](https://github.com/tomba-io/email/releases).

## Documentation

See the [official documentation](https://developer.tomba.io/).

### About Tomba

Founded in 2021, Tomba prides itself on being the most reliable, accurate, and in-depth source of Email address data available anywhere. We process terabytes of data to produce our Email finder API, company.

[![image](https://avatars.githubusercontent.com/u/67979591?s=200&v=4)](https://tomba.io/)

## Contribution

1. Fork it (<https://github.com/tomba-io/email/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## License

Please see the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html) file for more information.