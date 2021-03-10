package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"gopkg.in/yaml.v2"
)

type mailer struct {
	From           string `yaml:"from"`
	RecipientEmail string `yaml:"recipient_email"`
	RecipientName  string `yaml:"recipient_name"`
	SmtpServer     string `yaml:"smtp_server"`
	SmtpLogin      string `yaml:"smtp_login"`
	SmtpPassword   string `yaml:"smtp_password"`
}

func (m *mailer) Init(configFile string) *mailer {
	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Printf("yamlFile.get err #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, m)

	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return m
}

func (m *mailer) send(subject string, content string) {
	to := []string{m.RecipientEmail}
	auth := sasl.NewPlainClient("", m.SmtpLogin, m.SmtpPassword)
	msg := strings.NewReader("To: " + m.RecipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + content + "\r\n")

	err := smtp.SendMail(m.SmtpServer, auth, m.From, to, msg)

	if err != nil {
		log.Fatal(err)
	}
}

func parseTodo(todoFile string) []string {
	todoContents, err := ioutil.ReadFile(todoFile)

	if err != nil {
		log.Fatalf("Failed to read %v: %v", todoFile, err)
	}

	lines := strings.Split(string(todoContents), "\n")
	var incompleteTodos []string

	for _, line := range lines {
		if strings.Contains(line, "[ ]") {
			todo := " * " + strings.Split(line, "] ")[1]
			incompleteTodos = append(incompleteTodos, todo)
		}
	}

	return incompleteTodos
}

func composeMessage(todos string) string {
	message := "Here are the things that need to be done\n\n" +
		todos + "\n\nMake sure to mark complete tasks as done." +
		"\n\nRegards,\nReminderd"

	return message
}

func main() {
	configFile := flag.String("config", "config.yaml", "location of the config file")
	todoFile := flag.String("todo", "todo.org", "an org file containing a todo list")
	flag.Parse()

	var autoMailer mailer
	autoMailer.Init(*configFile)

	todoList := parseTodo(*todoFile)
	todoString := strings.Join(todoList, "\n")
	message := composeMessage(todoString)

	autoMailer.send("Daily reminders", message)
}
