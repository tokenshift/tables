package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/olekukonko/tablewriter"
)

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	var args struct {
		Filename string            `kong:"arg,required,type='path',help='The CSV file with the table definition.'"`
		Filters  map[string]string `kong:"short='f',name='filter',help='Optional column filter(s). Only take results that match.'"`
		Number   int               `kong:"short='n',default=1,help='Number of rolls/selections to make. Defaults to 1.'"`
		Output   string            `kong:"short='o',enum='simple,table,csv',default='simple',help='Output format. Simple, tabular, or CSV.'"`
	}

	ctx := kong.Parse(&args, kong.UsageOnError())
	if args.Filename == "" {
		ctx.FatalIfErrorf(fmt.Errorf("<filename> is required"))
	}

	table, err := LoadFile(args.Filename)
	ctx.FatalIfErrorf(err)

	filtered := table.Filter(args.Filters)

	switch args.Output {
	case "simple":
		filtered.DisplaySimpleRows(args.Number)
	case "csv":
		filtered.DisplayCSVRows(args.Number)
	case "table":
		filtered.DisplayTableRows(args.Number)
	default:
		filtered.DisplaySimpleRows(args.Number)
	}
}

func (table Table) DisplaySimpleRows(count int) {
	for n := 0; n < count; n++ {
		row := table.RandomRow()
		row = row.ApplyDieRolls(Rand)
		fmt.Println(row)
	}
}

func (table Table) DisplayCSVRows(count int) {
	writer := csv.NewWriter(os.Stdout)

	writer.Write(table.Columns.Names())

	for n := 0; n < count; n++ {
		row := table.RandomRow()
		row = row.ApplyDieRolls(Rand)
		writer.Write(row)
	}

	writer.Flush()
}

func (table Table) DisplayTableRows(count int) {
	writer := tablewriter.NewWriter(os.Stdout)
	writer.SetHeader(table.Columns.Names())

	for n := 0; n < count; n++ {
		row := table.RandomRow()
		row = row.ApplyDieRolls(Rand)
		writer.Append(row)
	}

	writer.Render()
}
