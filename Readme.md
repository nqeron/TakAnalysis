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

[takanalysis doc]:analysis/readme.md

## [analyze][analyze doc]
Single use command to analyze and annotate a given file or directory.

```
  analyze [-rv] [file ...] {-out ...}
```

## [transformPTN][transformPTN doc]
  A program that will [transformPTN][transformPTN doc]  into a symmetrically unique PTN
  (Can also be run on a directory)

```
  transformPTN [-rv] [file ...] {-out ...}
```

## [getTaks][getTaks doc]
  A program that will get the Taks from a given TPS, tagged in a PTN or stored in .tps. Will also work recursively over a directory

```
  getTaks [-rv] [file ...] {-out ...}
```

[analyze doc]:analyze/readme.md
[transformPTN doc]:cmd/transformPTN/readme.md
[getTaks doc]:cmd/getTaks/readme.md

## TO-DOs ##

- Testing / Refinement
  - Explore outputs on various games to help improve
  - More detailed comments / mark-up?
    - annotate tak, tinuë (& missed tinuë)
- GUI
