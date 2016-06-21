## TakAnalysis ##

This repository is dedicated to code that should help with analysis of Tak games.

## Installation ##

Once this is uploaded to github, you should be able to install it with:

```
go get github.com/nqeron/takanalysis
```

## Usage ##

At the moment, I only have one main program: takanalysis, run this from command line with:

```
takanalysis [file] [-d depth]
```

This will open a ptn file for analysis. Once opened, this can be directed with various commands:

- print -- prints out the entire PTN
- go [move] {black} -- moves the board to the given move in the PTN, use black to move to black's move

- n / next -- moves to the next ply
- p / prev -- moves to the previous ply
- ai -- calls up the ai analysis for the given position
- exit -- ends the program
- Alternatively navigate the plys with w/a/s/d:
  -- a/d - forward back one ply
  -- w/s - forward back one full turn

To start exploring alternative moves, just enter moves one at a time using standard notation (e.g. a5, 3c3>21). To undo moves, use 'u' or 'undo'. To stop exploring and exit out of the current chain, use 'q' or 'quit'

## Analysis ##

Currently, the analysis run locates changes in positional value that go significantly above or below the median. When the analysis is run, it will create a new file with the "-analysis.ptn" suffix, marked up with "!"s and "?"s

## TO-DOs ##

- Testing / Refinement
  - Explore outputs on various games to help improve
  - More detailed comments / mark-up?
- GUI
