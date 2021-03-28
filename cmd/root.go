package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ecator/gofile/server"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command
var addr string
var port int32
var webDir string
var dataDir string
var expireSec int32

func init() {
	rootCmd = &cobra.Command{
		Use:   "gofile",
		Short: "share files anonymously",
		RunE:  run,
	}
	rootCmd.Flags().StringVarP(&addr, "addr", "a", "localhost", "listen address")
	rootCmd.Flags().Int32VarP(&port, "port", "p", 1323, "listen port")
	rootCmd.Flags().StringVarP(&webDir, "webdir", "w", "web/dist", "web directory")
	rootCmd.Flags().StringVarP(&dataDir, "datadir", "d", "datadir", "data directory, must not exsit")
	rootCmd.Flags().Int32VarP(&expireSec, "expire", "e", 3600, "file expire seconds")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	if fileInfo, err := os.Stat(webDir); err != nil {
		return errors.New(fmt.Sprintf("webdir(%s) must exsit", webDir))
	} else if !fileInfo.IsDir() {
		return errors.New(fmt.Sprintf("webdir(%s) must be a directory", webDir))
	}
	if fileInfo, err := os.Stat(dataDir); err == nil {
		if !fileInfo.IsDir() {
			return errors.New(fmt.Sprintf("datadir(%s) must be a directory", dataDir))
		}
		fileInfos, _ := ioutil.ReadDir(dataDir)
		if len(fileInfos) > 0 {
			return errors.New(fmt.Sprintf("datadir(%s) must be a empty directory", dataDir))
		}
	}
	if err := os.MkdirAll(dataDir, os.FileMode(int(0700))); err != nil {
		return err
	}
	si := server.ServerInfo{
		ListenAddr: addr,
		ListenPort: int(port),
		WebDir:     webDir,
		DataDir:    dataDir,
		ExpireSec:  int(expireSec),
	}
	server.StartServer(si)
	return nil
}
