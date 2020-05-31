Definition and description of the available commads
--------------------------------------------------- 

The bot can rol both single dices `/d<sides>` for the most common dices or dice expressions  `/de <expression> <reason>`
The expressions follow mostly the [Dice Notation](https://en.wikipedia.org/wiki/Dice_notation) 

## Single Die commands

* `d2`: Roll a die with 2 sides
* `d4`: Roll a die with 4 sides
* `d6`: Roll a die with 6 sides
* `d8`: Roll a die with 8 sides
* `d10`: Roll a die with 10 sides
* `d12`: Roll a die with 12 sides
* `d20`: Roll a die with 20 sides
* `d100`: Roll a die with 100 sides

## Dice expressions

All dice expressions have the format `/de <expression> (<reason>)` where `<expression>` is a valid dice expression (see examples below) and `<reason>` it's a textual explanation of the purpose of the dice roll.

While the `<expression>` it's a required argument, the `<reason>` can be omited 

The number of dices it's the number before the `d`, just after the `d` it's the nubmer of sides of the dice, and then there are the optional modifiers and extra parameters for those modifiers

### Multiple dices 

* `/de 3d6`: Roll 3 dices of 6 sides and sum the total

### Keep (`k`)

* `/de 4d6k3`: Roll 4d6 and keep the top (k) 3 results

### Keep Lower (`kl`)

* `/de 4d6kl3`: Roll 4d6 and keep the lower (kl) 3 results

### Explode (`e`)

* `/de 5d10e`: Roll 5d10 and explode (e) roll again and add each 10

### Success (`s`)

* `/de 4d10s8`: Roll 4d10 and count as a success (s) each >= 8 

### Exploding Success (`es`)

* `/de 6d10es8`: Roll 4d10 and explode(e) each 10 and count the success (s) >= 8 

### Exploding Success (`r`)

* `/de 4d6r2`: Roll 4d6 and re-roll (r) any result  < 2

## Session Handling

Sessions are expected to last one game session and will be used to handle who was playing, record the dice rolls and group together a set of dice rolls

* `/start_session <name>`: Start a session named `<name>`
* `/rename_session <name>`: Rename the current session with a to `<name>`
* `/end_session`: End the current Session
