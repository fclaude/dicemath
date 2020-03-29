#!/bin/sh

go build .

scp dicemath fclaude@dicemath.recoded.cl:/home/fclaude/dicemath/dicemath.new
ssh fclaude@dicemath.recoded.cl "mv /home/fclaude/dicemath/dicemath /home/fclaude/dicemath/dicemath.old"
ssh fclaude@dicemath.recoded.cl "mv /home/fclaude/dicemath/dicemath.new /home/fclaude/dicemath/dicemath"
ssh fclaude@dicemath.recoded.cl "sudo service dicemath restart"
