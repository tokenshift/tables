# Dice

Golang library for parsing die roll specifications like "2d20kh1 + 3".

Valid die roll specifiers:

`1d20 + 2d6 + 5` - Roll and add together a d20, 2d6, and add 5.  
`2d20kh1` - Roll a d20 with advantage (keep the highest).  
`2d20kl1` - Roll a d20 with disadvantage (keep the lowest).  
`2d20k1` - Synonym for `kh1`.
`6` - The number 6 (don't roll anything).