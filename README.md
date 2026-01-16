# Pincher CLI
<p>
  <img width="1536" height="298" alt="pincher-cli-logo" src="https://github.com/user-attachments/assets/221e4863-6428-4dd1-ae3e-a7eb7ec7a646" />
</p>

Take control of your budget through the command line with an integrated REPL offering a user-friendly
experience when interacting with any server hosting the [Pincher API][pincher-api].

[pincher-api]: https://github.com/YouWantToPinch/pincher-api

![pincher-api](https://github.com/user-attachments/assets/e7640da4-7735-4388-a66f-676c99ffb262)

The Pincher CLI experience! A user checks balances, logs a transaction, and notes the change.

## Quick Start

You can install the CLI with:

```
go install github.com/YouWantToPinch/pincher-cli@latest
```

Then run it with `pincher-cli`.

When first starting the CLI, you can run `help` to see available commands, or `help -v` to see all possible commands
that may be available to you after creating a user and logging in.

You may edit your configuration for the CLI with the `config edit` command.
