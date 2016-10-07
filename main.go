package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	flag "github.com/ogier/pflag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	Version = "No version specified"
)

const (
	exitCodeOk             int = 0
	exitCodeError          int = 1
	exitCodeFlagParseError     = 10 + iota
	exitCodeAWSError
)

const helpString = `Usage:
  dynamo-env [-i] [--table dynamo_table] [name=value ...] [utility [argument ...]]

Flags:
  -h, --help    Print this help message
  -r, --region  The AWS region the table is in
  -t, --table   The name of the DynamoDB table
      --version Print the version number
`

type Item struct {
	Name  string
	Value string
}

func (i Item) String() string {
	return fmt.Sprintf("%s=%s", i.Name, i.Value)
}

var (
	f = flag.NewFlagSet("flags", flag.ContinueOnError)

	// options
	helpFlag    = f.BoolP("help", "h", false, "Show help")
	regionFlag  = f.StringP("region", "r", "us-east-1", "The AWS region")
	tableFlag   = f.StringP("table", "t", "", "The Dynamo table")
	versionFlag = f.Bool("version", false, "Print the version")

	// flags
	ignoreFlag = f.BoolP("ignore", "i", false, "Ignore the inherited environment")
)

func main() {
	var (
		environ []string
		err     error
		items   []Item
	)
	keys := make(map[string]int)

	if err = f.Parse(os.Args[1:]); err != nil {
		fmt.Println(err.Error())
		os.Exit(exitCodeFlagParseError)
	}

	if *helpFlag == true {
		fmt.Print(helpString)
		os.Exit(exitCodeOk)
	}

	if *versionFlag == true {
		fmt.Println(Version)
		os.Exit(exitCodeOk)
	}

	if *tableFlag == "" {
		fmt.Print(helpString)
		os.Exit(exitCodeFlagParseError)
	}

	// copy the shell environment
	if *ignoreFlag == false {
		environ = os.Environ()

		// parse existing environment
		for i, e := range environ {
			chk := strings.Split(e, "=")
			keys[chk[0]] = i
		}
	}

	// pull items from dynamo
	items, err = getDynamoItems(*tableFlag)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(exitCodeError)
	}

	// parse items in args
	cmdkv, commands := parseArguments(f.Args())
	items = append(items, cmdkv...)

	// join items to environ
	for _, item := range items {
		if *ignoreFlag == false {
			if _, exists := keys[item.Name]; exists {
				environ[keys[item.Name]] = item.String()
				continue
			}
		}
		environ = append(environ, item.String())
	}

	// print environ
	if len(commands) == 0 {
		for _, item := range environ {
			fmt.Println(item)
		}
		os.Exit(exitCodeOk)
	}

	// capture interrupts
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// run subshell with environ
	if err := runCommand(environ, commands); err != nil {
		os.Exit(exitCodeError)
	}
	os.Exit(exitCodeOk)
}

// return defined key pairs
func parseArguments(args []string) ([]Item, []string) {
	var (
		pairs    []Item
		commands []string
	)
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			parts := strings.SplitN(arg, "=", 2)
			pairs = append(pairs, Item{parts[0], parts[1]})
		} else {
			commands = append(commands, arg)
		}
	}
	return pairs, commands
}

// run the provided command with environment
func runCommand(env []string, args []string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	cmd.Start()
	return cmd.Wait()
}

func getDynamoItems(table string) ([]Item, error) {
	var items []Item

	sess, err := session.NewSession(&aws.Config{Region: regionFlag})
	if err != nil {
		return items, err
	}
	svc := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(table),
	}
	resp, err := svc.Scan(params)
	if err != nil {
		return items, err
	}
	for _, item := range resp.Items {
		items = append(items, Item{Name: *item["Name"].S, Value: *item["Value"].S})
	}
	return items, err
}
