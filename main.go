package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var Version = "development"
var usage string

func init() {
	usage = `stockuploader

	Usage:
	  stockuploader --username peter --password secret <file>...
	  stockuploader --username peter --password secret [--remote hostname] <file>...
	  stockuploader --username peter --password secret [--remote hostname] [--port 22] <file>...
	  stockuploader (-v | --version)

	  stockuploader (-h | --help)
	
	Arguments:
	  <file>		One or multiple files you'd like to upload

	Options:
	  -h --help                   Show this screen.
	  -v --version                Show version number.
	  -u --username=<username>    Provide your Adobe username.
	  -p --password=<password>    Provide your Adobe password.
	  -r --remote=<host>          sFTP remote host [default: sftp.contributor.adobestock.com].
	  -P --port=<port>            Remote port [default: 22].`
}

func main() {
	arguments, _ := docopt.ParseDoc(usage)
	if arguments["--version"] == true {
		print("Version: " + Version)
		os.Exit(0)
	}

	remote := cleanHostname(arguments["--remote"].(string))
	username := arguments["--username"].(string)
	password := arguments["--password"].(string)
	port := arguments["--port"].(string)
	files := cleanFiles(arguments["<file>"].([]string))

	conn, client := initiateSftpConnection(username, password, remote, port)
	defer conn.Close()
	defer client.Close()

	// uploading files to sftp
	for _, file := range files {
		copyFile(client, file, path.Base(file))
	}
}

func initiateSftpConnection(user string, pass string, remote string, port string) (*ssh.Client, *sftp.Client) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Initiate ssh connection
	conn, err := ssh.Dial("tcp", remote+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}

	// Create sftp channel
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}

	return conn, client
}

// Copies file to remote host
func copyFile(client *sftp.Client, source string, target string) {
	fmt.Print(source)
	// Create destination file
	dstFile, err := client.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	//dstFile, err := client.Create(target)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// The source file
	srcFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	// Copy file
	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\t\t(successful with %d bytes)\n", bytes)
}

// Takes a comma separated string and returns a cleaned up slice
func cleanFiles(files []string) []string {
	var cleanFileList []string
	for _, file := range files {
		file = strings.TrimSpace(file)
		if file != "" {
			cleanFileList = append(cleanFileList, file)
		}
	}
	return cleanFileList
}

// Removed trailing "sftp://" and also trimms the string
func cleanHostname(remote string) string {
	if strings.HasPrefix(remote, "sftp://") {
		remote = strings.TrimSpace(strings.Replace(remote, "sftp://", "", 1))

	}
	// TODO: Check for port
	return remote
}
