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
- 🤖 AI Chat — natural language queries powered by OpenAI with Tomba tools.
- 📦 Bulk CSV processing — enrich, verify, or find emails from CSV files with auto column mapping.
- 🧩 ClawHub Skills — reusable multi-step workflows that chain API calls into one command.

## Installation

### Quick Install (Linux / macOS)

```bash
curl -sSL https://raw.githubusercontent.com/tomba-io/tomba/master/res/package/scripts/install.sh | sh
```

### Quick Install (Windows PowerShell)

```powershell
irm https://raw.githubusercontent.com/tomba-io/tomba/master/res/package/scripts/install.ps1 | iex
```

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

![tomba help](svg/default.svg)

### Login

Sign in to Tomba account.

```bash
tomba login
```

<details>
<summary>Demo: tomba login</summary>

![tomba login](svg/login.svg)

</details>

### Domain Search

Instantly locate email addresses from any company name or website.

```bash
tomba search --target "tomba.io"
```

<details>
<summary>Demo: tomba search</summary>

![tomba search](svg/search.svg)

</details>

<details>
<summary>Demo: tomba search --json</summary>

![tomba search json](svg/search-json.svg)

</details>

<details>
<summary>Demo: tomba search --yaml</summary>

![tomba search yaml](svg/search-yaml.svg)

</details>

Slack Command: `/search tomba.io`

### Email Finder

Retrieves the most likely email address from a domain name, a first name and a last name.

```bash
tomba finder --target "tomba.io" --fist "mohamed" --last "ben rebia"
tomba finder --target "tomba.io" --fist "mohamed" --last "ben rebia" --enrich-mobile
```

<details>
<summary>Demo: tomba finder</summary>

![tomba finder](svg/finder.svg)

</details>

### Enrichment

Locate and include data in your emails.

```bash
tomba enrich --target "b.mohamed@tomba.io"
tomba enrich --target "b.mohamed@tomba.io" --enrich-mobile
```

<details>
<summary>Demo: tomba enrich</summary>

![tomba enrich](svg/enrich.svg)

</details>

Slack Command: `/enrich b.mohamed@tomba.io`

### Author Finder

Instantly discover the email addresses of article authors.

```bash
tomba author --target "https://clearbit.com/blog/company-name-to-domain-api"
```

<details>
<summary>Demo: tomba author</summary>

![tomba author](svg/author.svg)

</details>

Slack Command: `/author https://clearbit.com/blog/company-name-to-domain-api`

### Linkedin Finder

Instantly discover the email addresses of Linkedin URLs.

```bash
tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia"
tomba linkedin --target "https://www.linkedin.com/in/mohamed-ben-rebia" --enrich-mobile
```

<details>
<summary>Demo: tomba linkedin</summary>

![tomba linkedin](svg/linkedin.svg)

</details>

### Email Sources

Find email address source somewhere on the web.

```bash
tomba sources --target "ab@tomba.io"
```

<details>
<summary>Demo: tomba sources</summary>

![tomba sources](svg/sources.svg)

</details>

### Email Verifier

Verify the deliverability of an email address.

```bash
tomba verify --target "b.mohamed@tomba.io"
```

<details>
<summary>Demo: tomba verify</summary>

![tomba verify](svg/verify.svg)

</details>

Slack Command: `/checker b.mohamed@tomba.io`

### Count

Returns total email addresses we have for one domain.

```bash
tomba count --target "tomba.io"
```

<details>
<summary>Demo: tomba count</summary>

![tomba count](svg/count.svg)

</details>

### Status

Returns domain status if is webmail or disposable.

```bash
tomba status --target "tomba.io"
```

<details>
<summary>Demo: tomba status</summary>

![tomba status](svg/status.svg)

</details>

### Phone Finder

Search for phone numbers given an email, domain, or LinkedIn URL.

```bash
tomba phone-finder --email "wade@zapier.com"
tomba phone-finder --domain "tomba.io"
tomba phone-finder --linkedin "https://www.linkedin.com/in/alex-maccaw-ab592978"
tomba phone-finder --domain "stripe.com" --full
```

<details>
<summary>Demo: tomba phone-finder</summary>

![tomba phone-finder](svg/phone-finder.svg)

</details>

### Phone Validator

Validate phone numbers.

```bash
tomba phone-validator --phone "+14155552671"
tomba phone-validator --phone "4155552671" --country-code US
```

<details>
<summary>Demo: tomba phone-validator</summary>

![tomba phone-validator](svg/phone-validator.svg)

</details>

### Reveal (Company Search)

Search for companies using natural language or structured filters.

```bash
tomba reveal --query "SaaS startups in Europe"
tomba reveal --country US,UK --industry Technology
tomba reveal --country US --size 101-500,501-1000 --page 2
```

<details>
<summary>Demo: tomba reveal</summary>

![tomba reveal](svg/reveal.svg)

</details>

### Similar Domains

Retrieve domains similar to a specific domain.

```bash
tomba similar --target "tomba.io"
```

<details>
<summary>Demo: tomba similar</summary>

![tomba similar](svg/similar.svg)

</details>

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

<details>
<summary>Demo: tomba whoami</summary>

![tomba whoami](svg/whoami.svg)

</details>

### Version

Print version number and build information.

```bash
tomba version
```

<details>
<summary>Demo: tomba version</summary>

![tomba version](svg/version.svg)

</details>

---

### AI Chat

Open AI chat with Tomba tools — auto-initializes your project on first run. Uses OpenAI function calling to chain Tomba API tools together.

```bash
export OPENAI_API_KEY=sk-...
tomba chat
```

Example session:

```
✓ Authenticated as info@tomba.io

Tomba AI Chat — type your question, 'quit' to exit
Example: Find the VP Sales at Stripe and get their email

You: Find the VP Sales at Stripe and get their email
  workflow: domain_search -> email_finder -> email_verifier

## Results
**Email found:** d.singleton@stripe.com
**Verified:** deliverable (score: 95)
```

Bulk mode via chat:

```bash
tomba chat --file contacts.csv --type enrich
tomba chat --model gpt-4o-mini
```

<details>
<summary>Demo: tomba chat --help</summary>

![tomba chat](svg/chat.svg)

</details>

---

### Bulk CSV Processing

Process a CSV file in bulk — auto-maps columns or use custom mapping.

```bash
tomba bulk --file contacts.csv --type enrich
```

```
✓ Authenticated as info@tomba.io

  Columns: email, first_name, last_name, company
  Rows: 1000

  ✓ Using column 'email' for email

  Processing: ██████████████████████████████ 100% (1000/1000)

✓ 847/1000 emails found (84.7% hit rate)
✓ Results saved to contacts_enriched.csv
```

Supported operation types:

```bash
tomba bulk --file contacts.csv --type enrich
tomba bulk --file leads.csv --type verify --column "Email Address"
tomba bulk --file prospects.csv --type finder --domain-col company --first-col first --last-col last
tomba bulk --file domains.csv --type search --column domain -o results.csv
```

<details>
<summary>Demo: tomba bulk --help</summary>

![tomba bulk](svg/bulk.svg)

</details>

---

### ClawHub Skills

ClawHub is a built-in skill system that chains multiple Tomba API calls into reusable workflows. Run complex multi-step operations with a single command.

```bash
tomba skill list
```

<details>
<summary>Demo: tomba skill list — view all available skills</summary>

![tomba skill list](svg/skill-list.svg)

</details>

#### Built-in Skills

| Skill            | Workflow                                       | Description                                      |
| ---------------- | ---------------------------------------------- | ------------------------------------------------ |
| `find-email`     | finder → verify                                | Find & verify someone's email from name + domain |
| `company-emails` | count → search                                 | Search all emails + department breakdown         |
| `enrich-verify`  | enrich → verify → sources                      | Full email intelligence report                   |
| `company-intel`  | status → count → search → technology → similar | Complete company dossier                         |
| `lead-gen`       | search                                         | Find decision makers by department               |
| `linkedin-intel` | linkedin → enrich → verify                     | Full intel from LinkedIn profile                 |
| `domain-check`   | status → count                                 | Quick domain health check                        |
| `author-verify`  | author → verify                                | Find & verify article author email               |

#### Run Skills

```bash
tomba skill run find-email --target stripe.com first_name=David last_name=Singleton
tomba skill run company-intel --target tomba.io
tomba skill run enrich-verify --target b.mohamed@tomba.io
tomba skill run lead-gen --target stripe.com department=sales
tomba skill run linkedin-intel --target "https://www.linkedin.com/in/mohamed-ben-rebia"
tomba skill run domain-check --target tomba.io
```

<details>
<summary>Demo: tomba skill run domain-check — quick domain health check</summary>

![tomba skill run domain-check](svg/skill-domain-check.svg)

</details>

<details>
<summary>Demo: tomba skill run find-email — find and verify an email</summary>

![tomba skill run find-email](svg/skill-find-email.svg)

</details>

<details>
<summary>Demo: tomba skill run company-intel — full company dossier</summary>

![tomba skill run company-intel](svg/skill-company-intel.svg)

</details>

<details>
<summary>Demo: tomba skill run enrich-verify — email intelligence report</summary>

![tomba skill run enrich-verify](svg/skill-enrich-verify.svg)

</details>

#### Skill Details

```bash
tomba skill info company-intel
```

<details>
<summary>Demo: tomba skill info — workflow tree view</summary>

![tomba skill info](svg/skill-info.svg)

</details>

#### Create & Share Custom Skills

```bash
# Interactive skill builder
tomba skill create

# Import a skill from YAML
tomba skill import my-workflow.yaml

# Export a skill to share
tomba skill export find-email > find-email.yaml

# Remove a custom skill
tomba skill remove my-custom-skill
```

<details>
<summary>Demo: tomba skill export — export skill as YAML</summary>

![tomba skill export](svg/skill-export.svg)

</details>

<details>
<summary>Custom skill YAML format</summary>

```yaml
name: my-skill
description: My custom workflow
author: user
version: 1.0.0
tags: [custom, email]
inputs:
    - name: domain
      description: Target domain
      type: domain
      required: true
steps:
    - id: search
      name: Search emails
      action: search
      params:
          domain: "{{input.domain}}"
          limit: "20"
    - id: verify_first
      name: Verify first result
      action: verify
      params:
          email: "{{steps.search.data.emails.0.email}}"
output:
    format: summary
    title: My Results
    fields: [email, score, verification]
```

</details>

<details>
<summary>Demo: tomba skill --help</summary>

![tomba skill help](svg/skill-help.svg)

</details>

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
| bulk            | Process a CSV file in bulk — auto-maps columns or use custom mapping.                     |
| chat            | Open AI chat with Tomba tools — auto-initializes your project on first run.               |
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
| skill           | ClawHub — manage and run reusable Tomba workflows.                                        |
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
| `-j`      | `--json`   | output JSON format.                                                |
| `-k`      | `--key`    | Tomba API KEY.                                                     |
| `-o`      | `--output` | Save the results to file.                                          |
| `-p`      | `--port`   | Sets the port on which the HTTP server should bind. (default 3000) |
| `-s`      | `--secret` | Tomba API SECRET.                                                  |
| `-t`      | `--target` | TARGET SPECIFICATION Can pass email, Domain, URL, Linkedin URL.    |
| `-y`      | `--yaml`   | output YAML format.                                                |

### Output Formats

By default, results are displayed as human-readable text. Use flags to switch formats:

```bash
# Default: human-readable text output
tomba search --target "tomba.io"

# JSON output
tomba search --target "tomba.io" --json

# YAML output
tomba search --target "tomba.io" --yaml

# Save to file
tomba search --target "tomba.io" -o results.json
```

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
