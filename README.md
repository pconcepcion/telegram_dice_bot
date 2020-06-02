
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

### Register the bot on BotFather

To use the bot you need to register the bot on [BotFather](https://core.telegram.org/bots#6-botfather) by talking with the @BotFather bot and issuing the command `/newbot`

Use the `/setcommands` to tell @BotFather the available commands. This is a suggestion for the message to @BotFather:

```
d2 - Roll a d2
d4 - Roll a d4
d6 - Roll a d6
d8 - Roll a d8
d10 - Roll a d10
d12 - Roll a d12
d20 - Roll a d20
d100 - Roll a d100
de - Roll the given dice expression
start_session - Start a session with a name
end_session - End the current Session
```

Also is useful to set ask @BotFather to set a description (`/setdescription`). This is a suggestion for the description message to @BotFatther

```
The bot can roll single dices `/d<sides>` or dice expressions `/de <expression>`
(see [Dice Notation](https://en.wikipedia.org/wiki/Dice_notation)) 
Examples:
`/d2`
`/d100`
`/de 4d6k3` Roll 4d6: keep(k) the top 3 results
`/de 4d6kl3` Roll 4d6: keep the lower(kl) 3 results
`/de 5d10e` Roll 5d10: explode(e) roll again and add each 10
`/de 4d10s8` Roll 4d10: count as success(s) each >= 8
`/de 6d10es8` Roll 4d10: explode(e) each 10 and count # success(s) >= 8
`/de 4d6r2` Roll 4d6: re-roll(r) any result  < 2
``` 

Once you have the telegram api token, it can be set on the default config file  `$HOME/.telegram_dice_bot.yaml` (accets JSON, TOML, YAML, HCL, or Java properties file format): 

_Currently the only supported storage it's SQLite_

### Sample `.telegram_dice_bot.yaml`  

```yaml
--- # Telegam dice bot configuration 
api_token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
storage: sqlite://telegram_dice_bot.sqlite
```

or with the `TDB_API_TOKEN` environment variable

    export TDB_API_TOKEN="<YOUR API TOKEN HERE>"

## Docker

The bot can run on docker, a [`Docerfile`](Dockerfile) and [`docker-compose.yaml`](docker-compose.yaml) file are provided the targets to buile the docker images incudoded on the  [`Makefile`](Makefile)

So to build the dockerimage you just need to run `make docker` or `make docker-compose` and then you can start the image with `docker-compose up` (this will create persistend [docker volume](https://docs.docker.com/storage/volumes/) named `tbd-sqlite-volume`)

## Testing
Improve tests

## References
* [Statically compiling Go programs](https://www.arp242.net/static-go.html)
* [GORM](https://gorm.io) Go ORM 
  * [GORM Docs](https://gorm.io/docs)
  * [https://medium.com/@the.hasham.ali/how-to-use-uuid-key-type-with-gorm-cc00d4ec7100](https://medium.com/@the.hasham.ali/how-to-use-uuid-key-type-with-gorm-cc00d4ec7100)

## License

This code is released under the [BSD-3 Clause License](http://opensource.org/licenses/BSD-3-Clause)
