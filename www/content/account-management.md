---
title: Managing user accounts
weight: 4
menu: true
hideFromIndex: true
---

Ruruku comes with built-in user management.
Users are created using the CLI and are given permissions in the sytem.
When a Ruruku server starts, it can emit a _root token_ which is used to create the first users in the system.

## Starting a server and adding users
The `ruruku serve` command has two flags which emit a root token: `--root-token-file <filename>` and `--root-token-stdout`.
The former stores the root token in a file, the later prints it to the stdout.

To use the root token, pass it to subsequent ruruku calls, either using `--token` or by setting the `RURUKU_TOKEN` environment varibale.
For example:
```
ruruku serve --root-token-file /var/run/ruruku.root
export RURUKU_TOKEN=$(cat /var/run/ruruku.root)
```

The root token gives you all permissions within ruruku. At this point we'll be using them to create the actual users.
To do this, you can either create them one-by-one using [`ruruku user add`](../cli/ruruku_user_add) and [`ruruku user grant`](../cli/ruruku_user_grant),
or all at once from a YAML file. The [`examples/` folder](https://github.com/32leaves/ruruku/tree/master/examples) of the Ruruku repository has
an example of how such a YAML file looks like. Let's use this file to create a bunch of users:
```
ruruku user add -f examples/users.yaml
```

## Logging in on the command line
To login as a user on the command line, use [`ruruku user login <username>`](../cli/ruruku_user_login). After entering the correct password,
ruruku is going to print the token as an environment variable (can be changed using `-o`). This way, the login
command can be used as
```
eval $(ruruku login user <username>)
```
Notice: tokens have a limited lifespan, you may need to reauthenticate every now and then.

## Granting permissions
Ruruku has a set of permissions users can have. Each permission gives the user a particular ability, e.g. participating
in a test session, or creating other users. See below for a list of permissions:

| Permission         | Description                                                                                   |
|--------------------|-----------------------------------------------------------------------------------------------|
| user.add           | Enables a user to create users.                                                               |
| user.delete        | Enables a user to delete users.                                                               |
| user.grant         | Allows a user to grant permissions to users.                                                  |
| user.chpwd         | Allows a user to change the password of other users. All users can change their own password. |
| user.list          | Allows a user to retrieve a list of all users in the system, including their permissions.     |
| session.start      | Allows a user to start a new test session,                                                    |
| session.close      | Allows a user to close a test session.                                                        |
| session.view       | Allows a user to view a test session, including all test case details and current results.    |
| session.contribute | Allows a user to join a session, claim test cases and to contribute test results.             |
| session.modify     | Allows a user to add/remove/modify the testcases and annotations of a session (if marked modifiable and it's open) |

To grant permissions use [`ruruku user grant`](../cli/ruruku_user_grant). You must have the `user.grant` permission yourself to be able to do that.

## Change password
All users can change their own password using [`ruruku user chpwd`](../cli/ruruku_user_chpwd). To change the password of another user, you need the
`user.chpwd` permission, but can use just the same command.