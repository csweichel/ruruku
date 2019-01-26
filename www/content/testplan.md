---
title: Plans, sessions and sets
weight: 3
menu: true
hideFromIndex: true
---

In Ruruku a testcase is the basic "unit of test". A testcase is a single test that can fail or pass.
It has a description and a set of steps that must be performed by a tester to execute the testcase. Both fields support Markdown.
Additionally, test cases can have annotations which are basically key/value pairs used for conveying additional metadata or for selecing tests (see test sets below).

### Test Plan
A testplan is a collection of test cases. Testplans are YAML files that are meant to live next to the code they test.
To begin testing Ruruku will upload the testplan to a Ruruku instance, thus creating a test session (see below).

Next to the testcases, a testplan contains an ID, a description and a list of testsets.
A testplan can be created using `ruruku plan`.

## Test Session
A test session is the execution of a testplan - an instance of a testplan, basically. In a running session, testers can claim testcases
and contribute test results. If a session has been marked _modifiable_, the testcases can be altered in a session.
Test sessions are started using `ruruku session start`.

### Test Set
A test set is a selection of testcases of a testplan. Each test set has an _expression_ which selects testcases based on their annotations, ID or name.
Ruruku uses [`goevaluate`](https://github.com/Knetic/govaluate) for evaluating the expressions - see their [manual](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) for details.

The annotations of a test case are directly available in an expression, as are their name as `_name`, ID as `_id` and Group as `_group`.
See below for a few examples.

* Select testcases in the group `foobar`: `_group == "foobar"`
* Select testcases based on a `testlevel` annotation: `level == "basic"`
* Select testcases that require more than five users (user requirement expression as `requiredUsers` annotation): `level == "basic" && requiredUsers > 5`