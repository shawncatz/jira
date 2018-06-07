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
? Project for issue? TICKET
? Choose an issue type: Bug
? Choose a sprint: Backlog
? Title for issue? Title of issue
? Please enter a description [Enter to launch editor]
```

Some of the fields allow you to select from a list 
(thanks to [AlecAivazis/survey](https://github.com/AlecAivazis/survey)):

```bash
? Choose an issue type:  [Use arrows to move, type to filter]
❯ Bug
  Task
  Story
```

## Setup

Create a configuration file in `$HOME/.jira.yaml`

```yaml
jira_base: <base URL>
jira_user: <email>
jira_pass: <api key>
jira_project: PROJECT
jira_types: # First in the list is default
  - Bug
  - Task
  - Story
jira_sprints: # Automatically adds Backlog to the list
  - Sprint1
  - Sprint2
  - Sprint3
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

## Thanks

* [AlecAivazis/survey](https://github.com/AlecAivazis/survey)
* [spf13/cobra](https://github.com/spf13/cobra)
* [spf13/viper](https://github.com/spf13/viper)
* [andygrunwald/go-jira](https://github.com/andygrunwald/go-jira)
