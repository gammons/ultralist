package cmd

import (
	"fmt"
	"os"
	"strconv"
)

func argsToIDs(args []string) []int {
	var ids []int

	for _, id := range args {
		intID, err := strconv.Atoi(id)

		if err != nil {
			fmt.Printf("Could not parse ID: '%v'\n", id)
			os.Exit(1)
			return nil
		}

		ids = append(ids, intID)
	}

	return ids
}
