package models

import (
	"fmt"
	"strings"
	"time"
)

type Order struct {
	Id      uint      `db:"id"`
	Content string    `db:"content"`
	Created time.Time `db:"created"`
	AgentId uint      `db:"agent_id"`
}

func (o Order) String() string {
	return fmt.Sprintf("id: %d, content: %s, created: %s, agentId: %d", o.Id, o.Content, o.Created, o.AgentId)
}

func NewOrder(id uint, content string, created time.Time, agentId uint) *Order {
	return &Order{Id: id, Content: content, Created: created, AgentId: agentId}
}

type ConnectionString struct {
	Hostname string
	Database string
	Port     uint
	Username string
	Password string
}

func NewConnectionString(hostname string, database string, port uint, username string, password string) *ConnectionString {
	conn := ConnectionString{
		Hostname: hostname,
		Database: database,
		Port:     port,
		Username: username,
		Password: password,
	}

	return &conn
}

func (c ConnectionString) String() string {
	// "postgresql://[user[:password]@][netloc][:port][/dbname][?param1=value1&...]"

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		c.Username,
		c.Password,
		c.Hostname,
		c.Port,
		c.Database,
	)
}

type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}
