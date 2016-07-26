### Analyzes a single file or a directory of files and outputs an annotated ptn.

## About
  This is a command line way to analyze a given ptn, that doesn't require the standard console. It can also work (slowly) through a directory.

  All the files output are in standard ptn format, so may be used in any standard PTN viewer / browser.

## Method

Overall, this uses the methods defined in the analysis library of this repository. For more details about these functions, please see the analysis readme: [[takanalysis][takanalysis doc]]

[takanalysis doc]: https://github.com/nqeron/TakAnalysis/analysis/readme.md

##Signature
'''
analyze [---] file
'''

## Flags

- -out=[fileName] -- custom set the out put file / directory
- -depth=[int] -- set the maximum depth for the analysis and tinue search
- -debug -- print debug information
- -limit=[int] -- set the maximum time limit for searches
- -sensitivity=[int] -- set the annotation level to show
- -r -- loop over a directory of ptns
- -v -- show progress / non-essential text
