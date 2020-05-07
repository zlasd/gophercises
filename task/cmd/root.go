package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	db          *bolt.DB

	rootCmd = &cobra.Command{
		Use:  "task",
		Long: `task is a CLI for managing your TODOs.`,
	}

	addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new task to your TODO list",
		Run: func(cmd *cobra.Command, args []string) {
			db.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("TaskList"))
				if err != nil {
					return fmt.Errorf("create bucket: %s", err)
				}
				id, _ := b.NextSequence()
				task := strings.Join(args, " ")
				err = b.Put([]byte(strconv.Itoa(int(id))), []byte(task))
				if err != nil {
					log.Fatal(err)
					return err
				}
				fmt.Printf("Added \"%v\" to your task list.\n", task)
				return nil
			})
		},
	}

	doCmd = &cobra.Command{
		Use:   "do",
		Short: "Mark a task on your TODO list as complete",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("This is DO command.")
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("TaskList"))
				c := b.Cursor()

				index := 1
				target, err := strconv.Atoi(args[0])
				if err != nil {
					log.Fatal(err)
					return err
				}
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if index == target {
						b.Delete(k)
						fmt.Printf("You have completed the \"%v\" task.", string(v))
					}
					index++
				}

				return nil
			})
		},
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all of your incomplete tasks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("You have the following tasks:")
			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("TaskList"))
				c := b.Cursor()

				index := 1
				for k, v := c.First(); k != nil; k, v = c.Next() {
					fmt.Printf("%v. %v\n", index, string(v))
					index++
				}

				return nil
			})
		},
	}
)

// Execute execute CLI
func Execute() {
	path := "task.db"
	var err error
	db, err = bolt.Open(path, 0666, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	defer db.Close()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(listCmd)
}
