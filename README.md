# Jira CLI

Basic golang-based `Jira` CLI. Allows you to `create`, `show`, and `browse` to issues

## Setup

Create a configuration file in `$HOME/.jira.yaml`

```yaml
jira_base: <base URL>
jira_user: <email>
jira_pass: <api key>
```

* `base URL`: The URL to your Atlassian Cloud JIRA account (`https://yourname.atlassian.net`)
* `email`: The email of your account
* `api key`: Your Atlassian API token, see [here](https://confluence.atlassian.com/cloud/api-tokens-938839638.html) for more info

## Installation

The easiest way to setup the tool is to download a prebuilt release from Github.

### Install from source

If you have the `go` toolchain configured on your computer, you can install the tool 
with the following command:

> go get github.com/shawncatz/jira

You can update the installed version with:

> go get -u github.com/shawncatz/jira
