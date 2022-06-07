package config

import "fmt"

type File struct {
	Application Application `yaml:"application"`
	Database    Database    `yaml:"database"`
	Kafka       Kafka       `yaml:"kafka"`
}

type Application struct {
	Name     string `yaml:"name"`
	TestData bool   `yaml:"testData"`
}

type Broker struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
}

type Kafka struct {
	Brokers          BrokersArray     `yaml:"brokers"`
	IssueOrderTopics IssueOrderTopics `yaml:"issueOrderTopics"`
}

type IssueOrderTopics struct {
	IssueOrder        string `yaml:"issueOrder"`
	UndoIssueOrder    string `yaml:"undoIssueOrder"`
	RemoveOrder       string `yaml:"removeOrder"`
	MarkOrderIssued   string `yaml:"markOrderIssued"`
	ConfirmIssueOrder string `yaml:"confirmIssueOrder"`
}

type Database struct {
	DBMS     string `yaml:"dbms"`
	DB       string `yaml:"db"`
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type BrokersArray []Broker

func (d *Database) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.DB)
}

func (brokersArray BrokersArray) String() []string {
	brokersStr := make([]string, len(brokersArray))

	for i, broker := range brokersArray {
		brokersStr[i] = fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	}

	return brokersStr
}
