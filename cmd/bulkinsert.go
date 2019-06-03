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
	"github.com/Fjolnir-Dvorak/manageAMQ/amq"
	"github.com/Fjolnir-Dvorak/manageAMQ/queue"
	"github.com/Fjolnir-Dvorak/manageAMQ/ui"
	"github.com/pkg/errors"
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
	var fileList = &queue.FileList{}
	if len(amqFiles) != 0 {
		files, err := queue.BuildFromFileList(amqFiles)
		if err != nil {
			//fmt.Println(err)
			return
		}
		fileList.Append(files)
	}
	if amqDir != "" {
		files, err := queue.BuildFromDirectory(amqDir)
		if err != nil {
			//fmt.Println(err)
			return
		}
		fileList.Append(files)
	}

	runner := queue.NewRunner(fileList)

	runner.Waiter.Add(1)
	go func(run *queue.ActiveRunner) {
		defer runner.Waiter.Done()
		if run.Paused {
			<-run.Channel
		}
		err := amq.Connect(amqHost, amqUsername, amqPassword, amqPort)
		if err != nil {
			//fmt.Println(err)
			return
		}
		defer amq.Disconnect()

		for fileList.HasNextFile() {
			currentFile, err := fileList.GetNextFile()
			if err != nil {
				return
			}
			err = helperInsertSingleFile(currentFile, run)
			if err != nil {
				//fmt.Println("ERROR: something happened while reading file. aborting...")
				//fmt.Printf( "ERROR: %s", err)
				return
			}
		}
		return
	}(runner)

	if amqUi {
		ui.StartUI(runner)
	} else {
		runner.Paused = false
		runner.Channel <- struct{}{}
		runner.Waiter.Wait()
	}
}

func helperInsertSingleFile(fileContainer *queue.SingleFile, runner *queue.ActiveRunner) (err error) {
	file, err := os.OpenFile(fileContainer.FullPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		//fmt.Println("ERROR: Could not open file")
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	size, err := amq.GetEnqueuedCount(amqQueue, amqServicePort, amqUsername, amqPassword)
	if err != nil {
		size = 0
	}

	fileContainer.ReadingPosition = 1
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				//fmt.Println("Finished")
				break
			}
			//fmt.Println("ERROR: unexpected behaviour while reading file line by line")
			return err
		}

		for size >= amqQueuedLimit {
			for runner.Paused {
				<-runner.Channel
			}
			if !runner.Running {
				return errors.New("Exited")
			}
			time.Sleep(2 * time.Second)
			size, err = amq.GetEnqueuedCount(amqQueue, amqServicePort, amqUsername, amqPassword)
		}
		for runner.Paused {
			<-runner.Channel
		}
		if !runner.Running {
			return errors.New("Exited")
		}
		err = amq.SendMessage(amqQueue, line)
		if err != nil {
			//fmt.Printf("Could not write message into queue.\nMessage:\n%s\n\nAborting...\n", line)
			return err
		}
		size++
		fileContainer.ReadingPosition++
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
	bulkinsertCmd.Flags().MarkDeprecated()
}
