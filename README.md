
## Introduction
Telegram dice expressions bot written in Go

## Source
[Github repo](https://github.com/pconcepcion/telegram_dice_bot.git)

## Dependencies

* [Cobra](https://github.com/spf13/cobra): A Commander for modern Go CLI interactions
  * [Pflag](https://github.com/spf13/pflag): Drop-in replacement for Go's flag package, implementing POSIX/GNU-style --flags.
  * [Viper](https://github.com/spf13/viper): Go configuration with fangs.
* [Logrus](https://github.com/Sirupsen/logrus): Structured, pluggable logging for Go.
  * [Prefixed Log Formatter](github.com/x-cray/logrus-prefixed-formatter): Logrus Prefixed Log Formatter
* [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api/):Golang bindings for the Telegram Bot API
* [RPG Dice library](http://pconcepcion.github.io/dice/): Dice roll generators written in go

## Configuration 

The bot currently only needs one configuration parameter, the telegram api token, it can be set on the default config file  `$HOME/.telegram_dice_bot.yaml` (accets JSON, TOML, YAML, HCL, or Java properties file format): 

    --- 
    api_token: <YOUR API TOKEN HERE> 

or with the `TDB_API_TOKEN` environment variable

    export TDB_API_TOKEN="<YOUR API TOKEN HERE>"
 
## Testing
Improve tests

## License
* This code is released under the [BSD-3 Clause License](http://opensource.org/licenses/BSD-3-Clause)
