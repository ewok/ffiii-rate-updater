# FFIII Rate Updater

A command-line tool to update exchange rates in Firefly III.

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Usage](#usage)
5. [Configuration](#configuration)
6. [License](#license)
7. [Contact](#contact)

## Introduction

`ffiii-rate-updater` is a tool designed to fetch exchange rates for specified currencies and update them in Firefly III via its API. The tool is implemented in Golang and utilizes the Cobra library for command-line interaction.

The tool levarages [ Free Cuurency Exchange Rates API ](https://github.com/fawazahmed0/exchange-api) to fetch the latest exchange rates, and updates Firefly III accordingly.

## Features

- Fetches exchange rates for multiple currencies.
- Updates Firefly III with the latest exchange rates.
- Utilizes a fallback configuration for robust fetching.(not implemented yet)

## Installation

To build and install `ffiii-rate-updater`, ensure you have Go installed, then run:

```sh
go build
```

This will generate an executable named `ffiii-rate-updater`.

## Usage

### CLI
To use the tool, execute the following command in your terminal:

```sh
./ffiii-rate-updater update --config /path/to/config.yaml
```

Replace `/path/to/config.yaml` with the path to your configuration file.

Default paths for the config file are:
- Current directory: `./config.yaml`
- User home directory: `~/.ffiii-rate-updater.yaml`
- User config directory: `~/.config/ffiii-rate-updater/config`

### From Docker or docker-compose

TBD

## Configuration

The tool requires a configuration YAML file that specifies:

- `currencies`: A list of currency codes to fetch rates for.
- Firefly III API credentials, including:
  - `firefly.api_key`: Your Firefly III API key. [How to get an API key](https://docs.firefly-iii.org/how-to/firefly-iii/features/api/#personal-access-tokens)
  - `firefly.api_url`: The base URL for the Firefly III API.

Example configuration (config_example.yaml):

```yaml
firefly:
  api_key: YOUR_API_KEY
  api_url: https://your-firefly-instance.com/api/
currencies:
  - USD
  - EUR
  - JPY
```


## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for more details.

## Contact

Developed by Artur Taranchiev. Contact via email at [artur.taranchiev@gmail.com](mailto:artur.taranchiev@gmail.com).
