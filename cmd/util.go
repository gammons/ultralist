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
			fmt.Println("Could not parse ID: '%v'", id)
			os.Exit(0)
			return nil
		}

		ids = append(ids, intID)
	}

	return ids
}
