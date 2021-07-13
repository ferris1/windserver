package discovery

import (
"errors"
)

var (
	DefaultDiscovery = NewDiscovery()

	// Not found error when GetService is called
	ErrNotFound = errors.New("service not found")
	// Watcher stopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)

// service discovery
type Discovery interface {
	SetUp(...Option) error
	Options() Options
	Register(*Node, ...RegisterOption) error
	Deregister(*Node, ...DeregisterOption) error
	GetService(string, ...GetOption) (Service, error)
	ListServices(...ListOption) (Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Service struct {
	Type        string            	`json:"type"`
	Version   	string            	`json:"version"`
	Nodes     	[]*Node				`json:"nodes"`
	Metadata  	map[string]string 	`json:"metadata"`
}

type Node struct {
	Id       string            		`json:"id"`
	Address  string            		`json:"address"`
	Type 	 string 				`json:"type"`
	Metadata map[string]string 		`json:"metadata"`
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

type DeregisterOption func(*DeregisterOptions)

type GetOption func(*GetOptions)

type ListOption func(*ListOptions)

// Register a service node. Additionally supply options such as TTL.
func Register(n *Node, opts ...RegisterOption) error {
	return DefaultDiscovery.Register(n, opts...)
}

// Deregister a service node
func Deregister(s *Node) error {
	return DefaultDiscovery.Deregister(s)
}

// Retrieve a service. A slice is returned since we separate Name/Version.
func GetService(sType string) ([]*Service, error) {
	return DefaultDiscovery.GetService(sType)
}

// List the services. Only returns service names
func ListServices() ([]*Service, error) {
	return DefaultDiscovery.ListServices()
}

// Watch returns a watcher which allows you to track updates to the registry.
func Watch(opts ...WatchOption) (Watcher, error) {
	return DefaultDiscovery.Watch(opts...)
}

func String() string {
	return DefaultDiscovery.String()
}
