package cmd

import (
	"fmt"
	"db-dumper/db/mysql"
	"github.com/spf13/cobra"
)

var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Backup MySQL database",
	Run: func(cmd *cobra.Command, args []string) {
		backup := mysql.NewMySQLBackup()
		if err := backup.Backup(hostname, username, password, port, database, output); err != nil {
			fmt.Printf("Error backing up MySQL database: %v\n", err)
			return
		}
		fmt.Println("MySQL backup completed successfully!")
	},
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
}
