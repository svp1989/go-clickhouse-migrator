package command

import (
	"fmt"
)

const helpPath = "./help/main.txt"

func (c *Command) Help() error {
	help, err := c.migrator.Help(helpPath)
	if err != nil {
		return err
	}

	fmt.Print(help)

	return nil
}
