package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrExists    = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

// represent the list of hosts to run the port scan on
type HostsList struct {
	Hosts []string
}

// search - searches for host in list
func (h *HostsList) search(host string) (bool, int) {
	sort.Strings(h.Hosts)

	i := sort.SearchStrings(h.Hosts, host)
	if i < len(h.Hosts) && h.Hosts[i] == host {
		return true, i
	}

	return false, -1
}

// Add - add a host to the list
func (h *HostsList) Add(host string) error {
	if found, _ := h.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}

	h.Hosts = append(h.Hosts, host)
	return nil
}

// Remove - deletes a host from the list
func (h *HostsList) Remove(host string) error {
	if found, i := h.search(host); found {
		h.Hosts = append(h.Hosts[:i], h.Hosts[i+1:]...)
		return nil
	}

	return fmt.Errorf("%w: %s", ErrNotExists, host)
}

// Load - obtain hosts from the hosts file
func (h *HostsList) Load(hostsFile string) error {
	file, err := os.Open(hostsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	// read contents of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		h.Hosts = append(h.Hosts, scanner.Text())
	}

	return nil
}

// Save - saves hosts to a Hosts file
func (h *HostsList) Save(hostsFile string) error {
	output := ""

	for _, host := range h.Hosts {
		output += fmt.Sprintln(host)
	}

	return os.WriteFile(hostsFile, []byte(output), 0644)
}
