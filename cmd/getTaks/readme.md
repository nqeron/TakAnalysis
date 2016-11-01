### Analyzes a single PTN (with TPS tag) or TPS or a directory of such files and outputs all available taks for the given position

## About
  This is a command line tool to determine all possible taks from a given position.

  It returns a JSON formatted response in the format:

  [{ TPS: string,
    Moves: [string]  
  }]

## Method

This takes the TPS and parses it into a structured representation (using taktician's code) and does a depth one search for all available moves.

These are then filtered and parsed to ensure they are Taks.

##Signature
'''
getTaks [---] file
'''

## Flags

- -out=[fileName] -- custom set the output file
- -debug -- print debug information
- -r -- loop over a directory
- -v -- show additional error / feedback
