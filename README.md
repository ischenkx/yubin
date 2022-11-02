## Mailer
A service for email delivery

### Supports

 - Subscriptions
 - Scheduled publications
 - HTML templates
 - Checking whether an email was viewed
 - Reports

### Overview

#### Users
Each user is represented by their ID, Email, Name,
Surname and other meta information.
Users are the targets which will receive the letters.

#### Topics
Topics are named sets of users

#### Sources
Source is the email address from which the emails are sent

#### Templates
Templates help build dynamic emails

#### Reports
Reports can contain some helpful information on
the results of the email delivery

### How to run?

#### Prerequisites

1. Install Go and Docker Compose

#### ViewStat

ViewStat is a plugin that checks whether an email was viewed by a recipient.

How to set up?
1. Get some Redis server (probably you can use a free plan on RedisLabs)
2. Run ```cmd/viewstat/main.go``` somewhere (Heroku?)
3. Provide the redis client information and a link to the viewstat server in config

#### Mailer

1. Fix config (provide correct database information)
2. Install Go and Docker (Compose)
3. Run `bash ./docker-boot.sh`

Now you can interact with Mailer via a REST api (check out docs in ```internal/api/rest```)!