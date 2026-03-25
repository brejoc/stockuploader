package main

import (
	"fmt"
	"io"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	tomlConfig "github.com/brejoc/stockuploader/tomlConfig"
	"github.com/docopt/docopt-go"
	"github.com/mitchellh/ioprogress"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var Version = "development"
var usage string

const (
	defaultRemote = "sftp.contributor.adobestock.com"
	defaultPort   = "22"
)

type cliOutput int

const (
	NoOutput cliOutput = iota
	Text
)

func init() {
	usage = `stockuploader

	Usage:
	  stockuploader [options] <file>...
	  stockuploader (-v | --version)
	  stockuploader (-h | --help)

	Arguments:
	  <file>                      One or multiple files you'd like to upload

	Options:
	  -h --help                   Show this screen.
	  -v --version                Show version number.
	  -c --config=<path>         Load username, password, and optional remote/port from TOML (overrides auto-detection).
	  -u --username=<username>  Provide your Adobe username (overrides config).
	  -p --password=<password>  Provide your Adobe password (overrides config).
	  -r --remote=<host>        SFTP remote host (when omitted: config value, else Adobe Stock SFTP host).
	  -P --port=<port>          Remote port (when omitted: config value, else 22).`
}

func main() {
	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Fatal(err)
	}
	if arguments["--version"] == true {
		print("Version: " + Version)
		os.Exit(0)
	}

	cfgPath := resolveConfigPath(arguments)
	cfg, err := tomlConfig.Load(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	username := firstNonEmpty(strArg(arguments, "--username"), cfg.Username)
	password := firstNonEmpty(strArg(arguments, "--password"), cfg.Password)
	remote := firstNonEmpty(strArg(arguments, "--remote"), cfg.Remote)
	if remote == "" {
		remote = defaultRemote
	}
	remote = cleanHostname(remote)

	port := firstNonEmpty(strArg(arguments, "--port"), cfg.Port)
	if port == "" {
		port = defaultPort
	}

	if username == "" || password == "" {
		log.Fatal("username and password are required (use --username / --password and/or --config)")
	}

	files := cleanFiles(arguments["<file>"].([]string))

	conn, client := initiateSftpConnection(username, password, remote, port)
	defer conn.Close()
	defer client.Close()

	// uploading files to sftp
	for _, file := range files {
		copyFile(client, file, path.Base(file), Text)
	}
}

func strArg(arguments map[string]interface{}, key string) string {
	v, ok := arguments[key]
	if !ok || v == nil || v == false {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

// resolveConfigPath selects the TOML config file in this precedence order:
// 1) --config=<path> (exactly as provided)
// 2) ./stockuploader.toml
// 3) ~/.config/stockuploader.toml
// 4) ~/.stockuploader.toml
// 5) /etc/stockuploader.toml
func resolveConfigPath(arguments map[string]interface{}) string {
	if p := strArg(arguments, "--config"); p != "" {
		return p
	}
	return findFirstExistingDefaultConfig()
}

func findFirstExistingDefaultConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = ""
	}

	candidates := []string{
		"./stockuploader.toml",
	}
	if home != "" {
		candidates = append(candidates,
			filepath.Join(home, ".config", "stockuploader.toml"),
			filepath.Join(home, ".stockuploader.toml"),
		)
	}
	candidates = append(candidates, "/etc/stockuploader.toml")

	for _, p := range candidates {
		_, statErr := os.Stat(p)
		if statErr == nil {
			return p
		}
		if errors.Is(statErr, os.ErrNotExist) {
			continue
		}
		// Non "not-exist" errors (e.g. permission) should surface.
		return p
	}

	return ""
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
func copyFile(client *sftp.Client, source string, target string, output cliOutput) {

	// Create destination file
	dstFile, err := client.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// The source file
	srcFile, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}

	switch output {
	case Text:
		fmt.Println(source)
		sourceFileSize, err := srcFile.Stat()
		if err != nil {
			log.Fatal(err)
		}

		progressR := &ioprogress.Reader{
			Reader: srcFile,
			Size:   sourceFileSize.Size(),
		}

		// copy with progress
		_, err = io.Copy(dstFile, progressR)
		if err != nil {
			log.Fatal(err)
		}
	default:
		// copy silently
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			log.Fatal(err)
		}
	}
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
