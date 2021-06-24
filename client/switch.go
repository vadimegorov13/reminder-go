// Command switch
// Processes the commands from the CLI

package client

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Custome type for ids
type idsFlag []string

func (list idsFlag) String() string {
	return strings.Join(list, ", ")
}

func (list *idsFlag) Set(v string) error {
	*list = append(*list, v)
	return nil
}

type BackendHTTPClient interface {
	Create(title, message string, duration time.Duration) ([]byte, error)
	Edit(id, title, message string, duration time.Duration) ([]byte, error)
	Fetch(id []string) ([]byte, error)
	Delete(id []string) error
	Health(host string) bool
}

type Switch struct {
	client        BackendHTTPClient
	backendAPIURL string
	commands      map[string]func() func(string) error
}

func NewSwitch(uri string) Switch {
	httpClient := NewHTTPClient(uri)

	s := Switch{
		client:        httpClient,
		backendAPIURL: uri,
	}

	s.commands = map[string]func() func(string) error{
		"create": s.create,
		"edit":   s.edit,
		"fetch":  s.fetch,
		"delete": s.delete,
		"health": s.health,
	}

	return s
}

func (s Switch) Switch() error {
	cmdName := os.Args[1]

	cmd, ok := s.commands[cmdName]

	if !ok {
		return fmt.Errorf("Invalid command '%s'\n", cmdName)
	}

	return cmd()(cmdName)
}

func (s Switch) Help() {
	var help string
	for name := range s.commands {
		help += name + "\t --help\n"
	}
	fmt.Printf("Usage of '%s'\n <command> [<args>]\n%s", os.Args[0], help)
}

func (s Switch) create() func(string) error {
	return func(cmd string) error {
		createCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		title, message, duration := s.reminderFlags(createCmd)

		if err := s.checkArgs(3); err != nil {
			return err
		}

		if err := s.parseCmd(createCmd); err != nil {
			return err
		}

		res, err := s.client.Create(*title, *message, *duration)

		if err != nil {
			return wrapError("Could not create reminder", err)
		}

		fmt.Printf("Reminder created successfully:\n%s\n", string(res))

		return nil
	}
}

func (s Switch) edit() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		editCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		editCmd.Var(&ids, "id", "The ID of the reminder to edit")
		title, message, duration := s.reminderFlags(editCmd)

		if err := s.checkArgs(2); err != nil {
			return err
		}

		if err := s.parseCmd(editCmd); err != nil {
			return err
		}

		lastID := ids[len(ids)-1]
		res, err := s.client.Edit(lastID, *title, *message, *duration)

		if err != nil {
			return wrapError("Could not edit reminder", err)
		}

		fmt.Printf("Reminder edited successfully:\n%s\n", string(res))

		return nil
	}
}

func (s Switch) fetch() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		fetchCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		fetchCmd.Var(&ids, "id", "List of the reminder IDs to fetch")

		if err := s.checkArgs(1); err != nil {
			return err
		}

		if err := s.parseCmd(fetchCmd); err != nil {
			return err
		}

		res, err := s.client.Fetch(ids)

		if err != nil {
			return wrapError("Could not fetch reminder", err)
		}

		fmt.Printf("Reminder fetched successfully:\n%s\n", string(res))

		return nil
	}
}

func (s Switch) delete() func(string) error {
	return func(cmd string) error {
		ids := idsFlag{}
		deleteCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		deleteCmd.Var(&ids, "id", "List of the reminder IDs to delete")

		if err := s.checkArgs(1); err != nil {
			return err
		}

		if err := s.parseCmd(deleteCmd); err != nil {
			return err
		}

		err := s.client.Delete(ids)

		if err != nil {
			return wrapError("Could not delete reminder", err)
		}

		fmt.Printf("Reminder deleted successfully:\n%v\n", ids)

		return nil
	}
}

func (s Switch) health() func(string) error {
	return func(cmd string) error {
		var host string
		healthCmd := flag.NewFlagSet(cmd, flag.ExitOnError)
		healthCmd.StringVar(&host, "host", s.backendAPIURL, "Host to ping for health")

		if err := s.parseCmd(healthCmd); err != nil {
			return err
		}

		if !s.client.Health(host) {
			fmt.Printf("Host %s is down\n", host)
		} else {
			fmt.Printf("Host %s is up and running\n", host)
		}

		return nil
	}
}

func (s Switch) reminderFlags(f *flag.FlagSet) (*string, *string, *time.Duration) {
	title, message, duration := "", "", time.Duration(0)

	f.StringVar(&title, "title", "", "Reminder title")
	f.StringVar(&title, "t", "", "Reminder title")
	f.StringVar(&message, "message", "", "Reminder message")
	f.StringVar(&message, "m", "", "Reminder message")
	f.DurationVar(&duration, "duration", 0, "Reminder time")
	f.DurationVar(&duration, "d", 0, "Reminder time")

	return &title, &message, &duration
}

func (s Switch) checkArgs(minArgs int) error {
	if len(os.Args) == 3 && os.Args[2] == "--help" {
		return nil
	}
	if len(os.Args)-2 < minArgs {
		fmt.Printf(
			"Incorrect use of %s\n%s %s --help\n",
			os.Args[1], os.Args[0], os.Args[1],
		)
		return fmt.Errorf(
			"%s expects at least %d arg(s), %d provided\n",
			os.Args[1], minArgs, len(os.Args)-2,
		)
	}
	return nil
}

func (s Switch) parseCmd(cmd *flag.FlagSet) error {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		return wrapError("Could not parse '"+cmd.Name()+"' command flags", err)
	}
	return nil
}
