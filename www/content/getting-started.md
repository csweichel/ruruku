---
title: Getting Started
weight: 2
menu: true
hideFromIndex: true
---

Ruruku ships as a single executable that contains the ruruku server and command-line interface.
The former stores the data, serves the web-based UI and offers a gRPC API. You can talk to this
server using the UI or a command-line interface. The latter can also be used to create testplans
and add testcases.

## Download Ruruku
Ruruku can be downloaded from the [release page](https://github.com/32leaves/ruruku/releases). It runs on Linux,
OSX and Windows. For concenience sake, it's best to place the `ruruku` binary somewhere in your `PATH` (that's
not a requirement though).

If you want to run ruruku in Docker, there's also a [Docker image](https://hub.docker.com/r/csweichel/ruruku/) available.

## Create a testplan
You want to run a series of tests (otherwise you probably would not be here). This series of tests is what
ruruku calls a _testplan_. Such a testplan consists of a number _testcases_. Each testcase has a unique ID,
a name, description and a series of steps that testers need to execute. You can use Markdown in for the
description and steps. Testcases can be also be grouped - the web UI allows users to sort by
groups.

To create a testplan run [`ruruku plan`](../cli/ruruku_plan) and the CLI will guide you through the process.
If you want to convert an existing list of testcase descriptions, have a closer look at the flags
of [`ruruku plan`](../cli/ruruku_plan), especially `-y` which disables manual input.

Testplans are meant to live next to your code, i.e. it's best to check them into your repository.

## Start a test session
The easiest way to start a test session is `ruruku start testplan.yaml`. This will launch a new ruruku server
on your machine and start the session itself. This is great if you want to go through the testplan yourself or
just share it with others on the same network. Services like [Serveo](https://serveo.net) let you share that
local server across the internet. `ruruku start` will also create a user called `admin` (password: `admin`),
which has all permissions; use that one to get started quickly.

More commonly, you'll have a central ruruku instance hosted somewhere (e.g. using the [Docker image](https://hub.docker.com/r/csweichel/ruruku/)).
In that case run `ruruku session start --server your-ruruku-server.com:1234 --plan testplan.yaml` to start the session.
If you don't provide a name for the session (using `--name`), ruruku will come up with one.

## Participate in a session
You can join a test session using the web-based UI ruruku provides. Go to the server where the session was started,
choose the session you want to join, enter your name and you're good to go.

During a session, participants can _claim_ testcases by which they indicate that they would like to execute
that particular test. This way all tester can coordinate their efforts. Once you've claimed a testcase, you
can also contribute back the result of your test. The UI will automatically update for all participants showing
the latest state of the test.

The UI is great for human testers. If you have machines joining the session (e.g. automated UI tests), the
command-line interface makes for an easier integration. You can perform just the same operations that are
available in the web using the CLI.

## When you're done testing
Once you're done testing, you can close the session using the command-line interface with [`ruruku session close`](../cli/ruruku_session_close).
This marks the end of the test and marks the session as readonly.

To integrate the test result as part of your release or deployment process, [`ruruku session status`](../cli/ruruku_session_status) comes in handy.
Make sure to check out the `-o` flag with which you download the test result as JSON file.

