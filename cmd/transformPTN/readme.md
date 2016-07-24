# transformPTN

### Transforms PTN into symmetrically agnostic form

### TRANSPoSE

#### Tak Rotationally Agnostic Notation System Positioned on Side Edge

## About

For every Tak game out there, there are a number (15) other identical positions, just in different rotations / mirroring.

This code will transform 15 of those into one standard, positioned on the edge (a1), so that tak games that are identical will be easier to compare.

## Method

To begin, the first move is positioned as close to a1 as possible.

Then, the positions are searched for the first move that defines the chirality of the board. This move will be positioned in the bottom right corner of the board.

Everything else is translated accordingly.

### Example:

Consider the following game:
![http://goo.gl/msazFp][http://goo.gl/msazFp]

After the appropriate transforms, the game will look like this:
![http://goo.gl/jMmaQ6][http://goo.gl/jMmaQ6]


## Usage

### [transformPTN][transformPTN doc]
  A program that will [transformPTN][transformPTN doc]  into a symmetrically unique PTN
  (Can also be run on a directory)

```
  transformPTN [-rv] [file ...] {-out ...}
```
- -r=true - recurse over a directory
- -v=true - verbose mode
- -out= sets output file / directory, otherwise uses default
