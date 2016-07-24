# TakAnalysis

This repository is dedicated to code that should help with analysis of Tak games.

#Installation

Once this is uploaded to github, you should be able to install it with:

```
go get github.com/nqeron/takanalysis
```

#Programs #

## [takanalysis][takanalysis doc]
The main analysis program, [takanalysis][takanalysis doc], can be run from command line with:

```
takanalysis [file] [-d depth]
```

[takanalysis doc]: https://github.com/nqeron/TakAnalysis/analysis/readme.md

## [transformPTN][transformPTN doc]
  A program that will [transformPTN][transformPTN doc]  into a symmetrically unique PTN
  (Can also be run on a directory)

```
  transformPTN [-rv] [file ...] {-out ...}
```
- -r=true - recurse over a directory
- -v=true - verbose mode
- -out= sets output file / directory, otherwise uses default

[transformPTN doc]:https://github.com/nqeron/TakAnalysis/cmd/transformPTN/readme.md

## TO-DOs ##

- Testing / Refinement
  - Explore outputs on various games to help improve
  - More detailed comments / mark-up?
- GUI
