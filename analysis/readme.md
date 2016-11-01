#TakAnalysis#

### Main analysis program ###

This will open a ptn file for analysis. Once opened, this can be directed with various commands:

- print -- prints out the entire PTN
- go [move] {black} -- moves the board to the given move in the PTN, use black to move to black's move

- n / next -- moves to the next ply
- p / prev -- moves to the previous ply
- ai -- calls up the ai analysis for the given position
- tps -- will print out the TPS of the current position
- exit -- ends the program
- Alternatively navigate the plys with w/a/s/d:
  -- a/d - forward back one ply
  -- w/s - forward back one full turn

To start exploring alternative moves, just enter moves one at a time using standard notation (e.g. a5, 3c3>21). To undo moves, use 'u' or 'undo'. To stop exploring and exit out of the current chain, use 'q' or 'quit'

## Analysis ##

Currently, the analysis run locates changes in positional value that go significantly above or below the median. When the analysis is run, it will create a new file with the "-analysis.ptn" suffix, marked up with "!"s and "?"s
