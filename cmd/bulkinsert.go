// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"fmt"
	"github.com/Fjolnir-Dvorak/manageAMQ/amq"
	"github.com/Fjolnir-Dvorak/manageAMQ/queue"
	"github.com/Fjolnir-Dvorak/manageAMQ/utils"
	"github.com/spf13/cobra"
	"io"
	"os"
	"time"
)


var (
	//amqQueue string    Used from insert.go
	amqFiles       []string
	verbose        bool
	interactive    bool
	amqDir         string
	amqQueuedLimit int
	amqUi          bool
)

// bulkinsertCmd represents the bulkinsert command
var bulkinsertCmd = &cobra.Command{
	Use:   "bulkinsert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: doInsertBulk,
}

func doInsertBulk(cmd *cobra.Command, args []string) {
	var fileList = queue.FileList{}
	if len(amqFiles) != 0 {
		files, err := queue.BuildFromFileList(amqFiles)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileList.Append(files)
	}
	if amqDir != "" {
		files, err := queue.BuildFromDirectory(amqDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileList.Append(files)
	}

	runner := queue.NewRunner(fileList)

	go func(run *queue.ActiveRunner) {
		for run.Running {
			if run.Paused {
				run.Lock()
				run.Cond.Wait()
				run.Unlock()
				continue
			}
			err := amq.Connect(amqHost, amqUsername, amqPassword, amqPort)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("... Connected")
			defer amq.Disconnect()

			for fileList.HasNextFile() {

				if interactive {
					reader := bufio.NewReader(os.Stdin)
					input, _, err := reader.ReadRune()
					if err != nil {
						return
					}
					switch input {
					case 0x000A, 'y', 'Y':
						fmt.Println("... reading")
						err = helperInsertSingleFile(filename)
						if err != nil {
							return
						}
					case 'n', 'N':
						fmt.Println("... continue")
						continue
					case 'a', 'A':
						fmt.Println("... Aborting")
						return
					default:
						fmt.Println("... That was no valid character. Perhaps you meant to say No...")
						fmt.Printf("... Following character was typed: %#U\n", input)
					}
				} else {
					err = helperInsertSingleFile(filename)
					if err != nil {
						fmt.Println("ERROR: something happened while reading file. aborting...")
						fmt.Printf( "ERROR: %s", err)
						return
					}
				}
			}
		}
		return
	}(&runner)
}

func helperInsertSingleFile(filename string) (err error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println("ERROR: Could not open file")
		return err
	}
	defer file.Close()
	lines, err := utils.CountLinesFromFile(file)
	reader := bufio.NewReader(file)
	fileinfo, err := file.Stat()
	fileinfo.Name()

	size, err := amq.GetEnqueuedCount(amqQueue, amqServicePort, amqUsername, amqPassword)
	if err != nil {
		fmt.Printf("ERROR: Could not get size of queue. Assuming it is empty.\n")
		fmt.Printf("%s\n", err)
		//Assume that the queue is empty and that there were no connection errors.
		size = 0
	}
	fmt.Printf("File %s; line %d from %d. %.2f%%\n", fileinfo.Name(), 0, lines, 0.0)
	fmt.Printf("Queue size: %d/%d\n", size, amqQueuedLimit)

	lineNumber := 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Finished")
				break
			}
			fmt.Println("ERROR: unexpected behaviour while reading file line by line")
			return err
		}
		fmt.Printf("I am here\n")
		time.Sleep(10 * time.Millisecond)

		for size >= amqQueuedLimit {
			if verbose {
				fmt.Println("filled queue to Max. Waiting to empty...")
			}
			time.Sleep(5 * time.Second)
			size, err = amq.GetEnqueuedCount(amqQueue, amqServicePort, amqUsername, amqPassword)
			if err != nil {
				fmt.Println("ERROR: unexpected behaviour while getting the size of the queue")
			}
			if verbose {
				fmt.Printf("Queue size: %d/%d\n", size, amqQueuedLimit)
			}
			fmt.Printf("File %s; line %d from %d. %.2f%%\n", fileinfo.Name(), lineNumber, lines, float64(lineNumber)/float64(lines))
		}

		if verbose || interactive {
			fmt.Printf("\n\n%d:\n%s\n", lineNumber, line)
		}
		if interactive {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Send message (yes, no, abort) (Y|n|a): ")
			input, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println("ERROR: Could not read character from STDIN.")
				return err
			}
			switch input {
			case 0x000A, 'y', 'Y':
				amq.SendMessage(amqQueue, line)
				fmt.Println("... Message send")
			case 'n', 'N':
				fmt.Println("... continue")
				continue
			case 'a', 'A':
				fmt.Println("... Aborting")
				return nil
			default:
				fmt.Println("... That was no valid character. Perhaps you meant to say No...")
				fmt.Printf("... Following character was typed: %#U\n", input)
			}

		} else {
			err = amq.SendMessage(amqQueue, line)
			if err != nil {
				fmt.Printf("Could not write message into queue.\nMessage:\n%s\n\nAborting...\n", line)
				return err
			}
		}
		size++
		lineNumber++
	}
	return nil
}

func init() {
	rootCmd.AddCommand(bulkinsertCmd)

	bulkinsertCmd.PersistentFlags().StringVarP(&amqQueue, "queue", "q", "",
		"queue to write the message into")
	bulkinsertCmd.PersistentFlags().StringArrayVarP(&amqFiles, "files", "f", amqFiles,
		"messages to send separated by line ending \\n")
	bulkinsertCmd.PersistentFlags().StringVarP(&amqDir, "directory", "d", "",
		"messages to send separated by line ending \\n")
	bulkinsertCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose logging. Prints the message which will be written into the queue.")
	bulkinsertCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false,
		"Asks for each message if it should be written. 'y' for Yes and 'n' for no.")
	bulkinsertCmd.PersistentFlags().IntVarP(&amqQueuedLimit, "max", "m", 10000,
		"Max limit of entries to queue at the same time.")
	bulkinsertCmd.PersistentFlags().BoolVar(&amqUi, "ui", false,
		"Starts a terminal ui with status informations")
}
