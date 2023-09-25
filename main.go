package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var username string
var password string
var remote string
var port string
var files string

func init() {
	flag.StringVar(&username, "username", "", "Your username")
	flag.StringVar(&password, "password", "", "Your password")
	flag.StringVar(&remote, "remote", "sftp.contributor.adobestock.com", "The remote host. Usually this is 'sftp.contributor.adobestock.com'.")
	flag.StringVar(&port, "port", "22", "Port of the remote host. Usually this is port 22")
	flag.StringVar(&files, "files", "", "Comma separated list of files you'd like to upload")
}

func main() {
	flag.Parse()

	remote = cleanHostname(remote)
	filePaths := cleanFiles(files)

	if username == "" {
		log.Fatal("You need to provide a username: --username peter")
	}

	if password == "" {
		log.Fatal("You need to provide a password: --password secret")
	}

	conn, client := initiateSftpConnection(username, password, remote, port)
	defer conn.Close()
	defer client.Close()

	// uploading files to sftp
	for _, file := range filePaths {
		fmt.Println("File: ", file)
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
	fmt.Printf("copied bytes: %d\n", bytes)
}

// Takes a comma separated string and returns a cleaned up slice
func cleanFiles(files string) []string {
	if files == "" {
		return make([]string, 0)
	}
	var cleanFileList []string
	fileList := strings.Split(files, ",")
	for _, file := range fileList {
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
