![ruruku](logo.png)

A simple manual test coordinator written in Go and TypeScript/React.

**Beware: this is a spare-time project and by no means ready for actual use just yet. There will be dragons.**

Ruruku offers text-based test-case description that is meant to live next to your code.
When the time has come to go through the tests, run `ruruku start testcases.yaml` to spawn a webserver that allows others to participate in the test.
During the test run the server persists a test-report YAML file which can be used as pre-deployment gate or kept for reference.

## Development
[![Open in Gitpod](http://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io#https://github.com/32leaves/ruruku)
[![Build Status](https://travis-ci.org/32leaves/ruruku.svg?branch=master)](https://travis-ci.org/32leaves/ruruku)

## FAQ

### What's with the name?
Ruruku is Maori and means *to draw together with a cord, bind together, lash, coordinate.*
It is pronounced just the way it's written - checkout the [Maori dictionary](http://maoridictionary.co.nz/search?idiom=&phrase=&proverb=&loan=&histLoanWords=&keywords=ruruku) for an audio sample.

### Why are you building this?
1. I wanted a fun spare-time project that integrates Go and React in a single project.
2. Gitpod allows for a new kind of tools which no longer require complex hosting so that they can be available on the Internet. Your workspace becomes your hosting platform. I wanted to explore this concept in a real-world use-case.
3. When testing Gitpod we still have a handful of manual test-cases. I hope that ruruku will be handy for testing those.
