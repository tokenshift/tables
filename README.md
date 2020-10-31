# Tables

Command-line app for RPG-style random tables.

## Usage

Roll a random encounter that takes place while travelling through the woods:

```
$ tables random_encounter.csv -f Terrain=woods
[1d4 + 2 = 4] angry porcupines summoned by an evil wizard in a cave hidden in the underbrush in the woods
```

### Arguments

* `<filename>` (Required)  
  The JSON file with the table definition (see below).
* `-f {column=value[,value...]} | --filter {column=value[,value...]` (Optional)  
  Optional column filters. Only take results that match.
  Multiple values can be listed, separated by commas.
* `-n {number} | --number {number}` (Optional)  
  Number of selections to make. Defaults to 1.
* `-o {format} | --output {format}`  (Optional)
  Output format; `simple`, `csv`, or `table`. Defaults to `simple`.

## Table Files

Individual tables that you can roll on are defined as CSV files. For example:

```csv
Combatants [Independent],"Boss [Independent,Percent=25]",Location,[Independent],Terrain
a group of gnolls,summoned by an evil wizard,in a copse of trees,in the,woods
a large bear,being corralled by a chaotic evil druid,in a cave hidden in the underbrush,,woods
[1d4+2] angry porcupines,,in a hideout carved into the trunk of a massive oak tree,,woods
,,hidden behind a dune,,desert
```

By default, a single roll is made, selecting an entire row. However, you can
mark columns as `[Independent]` to roll a value for that column separately.

You can also add a probability for rolling a specific column; e.g.
`[Percentage=25]` will produce a value for that column 25% of the time.

Any of your row or column values can include die rolls delimited by square
brackets (e.g. `[1d6+2]`). The square brackets will be replaced by the result of
the die roll.

If you want the literal square brackets to appear in the value and not be
replaced by a die roll, escape them using a backslash.

For valid die roll specifiers, see ./dice/README.md.