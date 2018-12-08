![ruruku](logo.png)

You and your team have a list of testcases you want to execute prior to a new release.
Because things are moving quickly many of those tests aren't yet automated, so testing becomes a team effort.
Ruruku helps coordinate this team effort by offering a single, low friction contact point where testers can sign up, claim test cases and provide feedback.

Ruruku offers a YAML-based testcase description that is meant to live next to your code.
When the time has come to go through the tests, run `ruruku session start --plan testcases.yaml` to start a test session on your ruruku installation.
You can also use `ruruku start testcases.yaml` to spawn a local webserver that allows others to participate in the test.

**Beware: this is a side project and it's early days. Here be dragons.**

## Getting started
Ruruku runs on OSX, Linux and Windows. You can either get going on your local machine, in a Gitpod or run this in a Docker container.
### On your local machine
To get started [download ruruku](https://github.com/32leaves/ruruku/releases).
You'll need a set of tests that you want to run through.
You can either create one yourself (`ruruku plan`) or use an [example](https://raw.githubusercontent.com/32leaves/ruruku/master/testplan-example.yaml).
Use `ruruku start testplan.yaml` to start the API serer and test session.

To share that session with others, [Serveo](https://serveo.net) comes in handy, which exposes local servers to the internet. This way you can run ruruku on your local machine and share it with others.

### Gitpod
Gitpod is an online IDE that offers instant workspaces in the cloud (full disclosure: this is what I work on during the day).
It makes trying things like ruruku a breeze. Ruruku itself was/is developed exclusively in Gitpod - I never once had to clone the repo locally.

You can either open use our [demo repository](https://gitpod.io/#github.com/32leaves/ruruku-demo) which also serves a good starter for your own projects, or jump right in with the [development workspace](https://gitpod.io#https://github.com/32leaves/ruruku) of ruruku which runs a full build.

### Hosting ruruku (Docker/Kubernetes)
Ruruku has a central server which hosts the Web UI (for tests) provides a gRPC based API for the command-line tools.
The ruruku server starts with `ruruku serve`. There also is a Docker image available for each release.

Note that the ruruku server needs a place where to store the data. By default that's in `/var/ruruku`.
To run a Docker container that makes this data persistent, you can use a volume:
```
docker run -p 8080:8080 -p 1234:1234 -v /path/on/my/machine:/var/ruruku csweichel/ruruku:latest
```

### Create a testsuite
To create your own testsuite run `ruruku plan` which will guide you through the process.
If you want to a converter that takes an existing testcase description and produces a ruruku one, make sure to look at `ruruku plan --help` and `ruruku plan add --help`.

## Development
[![Open in Gitpod](http://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/32leaves/ruruku)
[![Build Status](https://travis-ci.org/32leaves/ruruku.svg?branch=master)](https://travis-ci.org/32leaves/ruruku)
[![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## FAQ

### What's with the name?
Ruruku is Maori and means *to draw together with a cord, bind together, lash, coordinate.*
It is pronounced just the way it's written - checkout the [Maori dictionary](http://maoridictionary.co.nz/search?idiom=&phrase=&proverb=&loan=&histLoanWords=&keywords=ruruku) for an audio sample.

### Why are you building this?
1. I wanted a fun side project that integrates Go and React in a single project.
2. Gitpod allows for a new kind of tools which no longer require complex hosting so that they can be available on the Internet. Your workspace becomes your hosting platform. I wanted to explore this concept in a real-world use-case.
3. When testing Gitpod we still have a handful of manual testcases. I hope that ruruku will be handy for testing those.
