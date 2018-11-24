![ruruku](logo.png)

Suppose you and your team have a list of testcases you want to execute prior to a new release.
Because things are moving quickly many of those tests aren't yet automated, so testing becomes a team effort.
Ruruku helps coordinate this team effort by offering a single, low friction contact point where testers can sign up, claim test cases and provide feedback.

Ruruku offers a YAML-based testcase description that is meant to live next to your code.
When the time has come to go through the tests, run `ruruku start testcases.yaml` to spawn a webserver that allows others to participate in the test.
During the test run the server persists a test-report YAML file which can be used as pre-deployment gate or kept for reference.

**Beware: this is a side project and is by no means ready for actual use just yet. There will be dragons.**

## Getting started
At the moment the best way to run ruruku is in a Gitpod workspace, but there are other means, too:
- **Gitpod:** You can either open use our [demo repository](https://gitpod.io/#github.com/32leaves/ruruku-demo) which also serves a good starter for your own projects, or jump right in with the [development workspace](https://gitpod.io#https://github.com/32leaves/ruruku) of ruruku which runs a full build.
- **On your local machine:** Ruruku runs on [OSX, Linux and Windows](https://github.com/32leaves/ruruku/releases). To share a ruruku session with others, [Serveo](https://serveo.net) comes in handy, which exposes local servers to the internet. This way you can run ruruku on your local machine and share it with others. To get started download ruruku, run clone this repo and run `ruruku start example-suite.yaml && ssh -R 80:localhost:8080 serveo.net`.

### Create a testsuite
To create your own testsuite run `ruruku init` which will guide you through the process.
If you want to a converter that takes an existing testcase description and produces a ruruku one, make sure to look at `ruruku init --help` and `ruruku init testcase --help`.

## Development
[![Open in Gitpod](http://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/32leaves/ruruku)
[![Build Status](https://travis-ci.org/32leaves/ruruku.svg?branch=master)](https://travis-ci.org/32leaves/ruruku)

## FAQ

### What's with the name?
Ruruku is Maori and means *to draw together with a cord, bind together, lash, coordinate.*
It is pronounced just the way it's written - checkout the [Maori dictionary](http://maoridictionary.co.nz/search?idiom=&phrase=&proverb=&loan=&histLoanWords=&keywords=ruruku) for an audio sample.

### Why are you building this?
1. I wanted a fun side project that integrates Go and React in a single project.
2. Gitpod allows for a new kind of tools which no longer require complex hosting so that they can be available on the Internet. Your workspace becomes your hosting platform. I wanted to explore this concept in a real-world use-case.
3. When testing Gitpod we still have a handful of manual testcases. I hope that ruruku will be handy for testing those.
