# Jira CLI

Basic golang-based `Jira` CLI. Allows you to `create`, `show`, and `browse` to issues

## Show

Get a quick summary of a ticket.

> jira show TICKET-1234

```bash
TICKET-1234    : Title of issue
Type           : Bug
Priority       : Blocker
Assigned       : person
Description    :
Something was borked

person2 (person2@email.com)
Fixed!

person3 (person3@email.com)
Link to document
https://docs.google.com/document/d/...

https://yourcompany.atlassian.net/browse/TICKET-1234
```

## Browse

This will open the link to the ticket with your default browser (using
`open` on `macOS`)

> jira browse TICKET-1234

```bash
opening: https://anyperk.atlassian.net/browse/TICKET-1234
```

## Create

> jira create

```bash

```

## Setup

Create a configuration file in `$HOME/.jira.yaml`

```yaml
jira_base: <base URL>
jira_user: <email>
jira_pass: <api key>
```

* `base URL`: The URL to your Atlassian Cloud JIRA account 
   (`https://yourname.atlassian.net`)
* `email`: The email of your account
* `api key`: Your Atlassian API token, see [here](https://confluence.atlassian.com/cloud/api-tokens-938839638.html) for more info

## Installation

The easiest and recommended way to setup the tool is to download a prebuilt 
[release](https://github.com/shawncatz/jira/releases) from Github.

### Install from source

If you have the `go` toolchain configured on your computer, you can install the tool 
with the following command:

> go get github.com/shawncatz/jira

You can update the installed version with:

> go get -u github.com/shawncatz/jira
