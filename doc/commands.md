Definition and description of the available commads
--------------------------------------------------- 

The bot can rol both single dices `/d<sides>` for the most common dices or dice expressions  `/de <expression> <reason>`
The expressions follow mostly the [Dice Notation](https://en.wikipedia.org/wiki/Dice_notation) 

## Single Die commands

Simple, roll a single die and get the result, just predefined the most common dices for simpler access

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

Roll some dices, keep the top _X_ (_X_ must not be more than the number of dices rolled) and then sum the result of the dices kept

* `/de 4d6k3`: Roll 4d6 and keep the top 3 results

### Keep Lower (`kl`)

Roll some dices, keep the lower _X_ (_X_ must not be more than the number of dices rolled) and then sum the result of the dices kept

* `/de 4d6kl3`: Roll 4d6 and keep the lower  3 results

### Explode (`e`)

On exploding rolls the system rolls again the dices that get the maximum value of the dice, and then add all the results, including the re rolled dices. 

* `/de 5d10e`: Roll 5d10 and explode roll again and add each 10

### Success (`s`)

The _success_ rolls the system counts the dices with a result over the give threshold 

* `/de 4d10s8`: Roll 4d10 and count as a success each result of more 8 or more

### Exploding Success (`es`)

* `/de 6d10es8`: Roll 4d10 and explode (roll again) each 10 and count the number of success (results of 8 or more including the 10 that are rerolled)

### Exploding Success (`r`)

* `/de 4d6r2`: Roll 4d6 and re-roll any result lower 2 util you get 2 or more (there is a hard limit of 100 re-rolls) 

## Session Handling

Sessions are expected to last one game session and will be used to handle who was playing, record the dice rolls and group together a set of dice rolls

* `/start_session <name>`: Start a session named `<name>`
* `/rename_session <name>`: Rename the current session with a to `<name>`
* `/end_session`: End the current Session
