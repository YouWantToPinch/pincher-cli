# Introduction

### Purpose
pincher-sdk-go is a library providing convenient access to the Pincher REST API from applications written in Go.

### Background
When beginning to work with the [revoltgo](https://github.com/sentinelb51/revoltgo) library for Stoat, I learned a lot from the design philosophy within that codebase. I decided to adapt some of its methodology to the existing client package I had built for---and had long wanted to decouple from---the [Pincher CLI](https://github.com/YouWantToPinch/pincher-cli). This package is the result.

# Getting started

## Installation
Within a working Go environment, run the following command to install the library:

```bash
go get github.com/YouWantToPinch/pincher-sdk-go
```

## Usage
First import the package:

```go
import "github.com/YouWantToPinch/pincher-sdk-/pinchergo"
```

Then, initialize a new **client** using one of the following two functions:

```go
client, err := pinchergo.NewClientWithDefaults()
```


```go
client, err := pinchergo.NewClient(
  "http://localhost:8080", // base URL
  time.Second * 10, // client timeout
  time.Minute * 5, // cache reap interval
  // whether to use refresh token to get new access tokens after 401 responses
  true,
  )
```

The `Client` type is what you will use to make your API calls. It also contains an internal, in-memory cache you can use, whose type is aptly named `Cache`.

The client is given dedicated methods for each budgetary resource defined by the Pincher API. One of these are budget accounts, which feature inline documentation for each of those aforementioned methods explaining their methodology; this methodology applies almost entirely to each of the other budgetary resources such as groups, categories, payees, etc., so the documentation translates.

You can therefore find the inline documentation under the `requests_accounts.go` file.

## License: BSD 3-Clause

pincher-sdk-go is licensed under the BSD 3-Clause License. What this means is that:

#### You are allowed to:

+ Modify the code, and distribute your own versions.
+ Use this library in personal, open-source, or commercial projects.
+ Include it in proprietary software, without making your project open-source.

#### You are not allowed to:

- Remove or alter the license and copyright notice.
- Use the names "pincher-sdk-go" or its contributors for endorsements without permission.

### Attribution

As previously mentioned, much of the design methodology baked into this package is owed to the [revoltgo](https://github.com/sentinelb51/revoltgo) project. It is licensed under the BSD 3-Clause license, the same license as this SDK.

*NOTE: this package is in no way involved with nor endorsed by the creators nor contributors of RevoltGo.*
