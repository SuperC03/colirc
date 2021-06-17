package data

import (
	"fmt"
	"strings"
)

type Message struct {
	Tags    map[string]string // Optional Tags
	Source  string            // Optional Source
	Command string            // Command
	Params  []string          // Command Paramters
}

// Occurs if Input is Formatted Incorrectly
type SyntaxError struct{}

func (e *SyntaxError) Error() string {
	return "Badly Formated Input"
}

// Accepts Raw String and Returns Decoded Message
func UnmarshalMessage(input string) (*Message, error) {
	msg := &Message{}
	parts := strings.Split(input, " ")
	// Atleast a Command Should be Present
	if len(parts) < 1 {
		return nil, &SyntaxError{}
	}
	// Check for Tags
	if parts[0][0] == '@' {
		msg.Tags = make(map[string]string)
		for _, pair := range strings.Split(parts[0][1:], ";") {
			split := strings.Split(pair, "=")
			if len(split) != 2 {
				return nil, &SyntaxError{}
			}
			msg.Tags[split[0]] = split[1]
		}
	}
	// Check for Source
	if parts[0][0] == '@' && parts[1][0] == ':' {
		msg.Source = parts[1][1:]
	} else if parts[0][0] == ':' {
		msg.Source = parts[0][1:]
	}
	// Parse Command...Oh Gosh, Here we Go
	var cmd string
	if parts[0][0] == '@' && parts[1][0] == ':' { // Both Tags and Source are Present
		cmd = parts[2]
	} else if parts[0][0] == '@' || parts[0][0] == ':' { // Only Tags OR Source are Present
		cmd = parts[1]
	} else { // ¯\_(ツ)_/¯
		cmd = parts[0]
	}
	msg.Command = strings.ToUpper(cmd)
	// Put Everything Else in Params
	if parts[0][0] == '@' && parts[1][0] == ':' && len(parts) > 3 {
		msg.Params = parts[3:]
	} else if parts[0][0] == '@' && parts[1][0] == ':' { // No Params
		msg.Params = make([]string, 0)
	} else if (parts[0][0] == '@' || parts[0][0] == ':') && len(parts) > 2 {
		msg.Params = parts[2:]
	} else if parts[0][0] == '@' || parts[0][0] == ':' {
		msg.Params = make([]string, 0)
	} else if msg.Command != "" && len(parts) > 1 {
		msg.Params = parts[1:]
	} else if msg.Command != "" {
		msg.Params = make([]string, 0)
	} else {
		return nil, &SyntaxError{}
	}
	// Concactinate ":" Into Multi Spaced Entry
	for i, seg := range msg.Params {
		if seg[0] == ':' {
			temp := make([]string, len(msg.Params))
			copy(temp, msg.Params)
			msg.Params = msg.Params[:i]
			msg.Params = append(msg.Params, strings.Join(temp[i:], " "))
		}
	}
	return msg, nil
}

func MarshalMessage(msg *Message) (string, error) {
	var output string

	// Handle Tags
	if len(msg.Tags) > 0 {
		output += "@"
		counter := 0
		for k, v := range msg.Tags {
			output += fmt.Sprintf("%s=%s", k, v)
			if counter < len(msg.Tags)-1 {
				output += ";"
			}
			counter += 1
		}
		output += " "
	}

	// Handle Source
	if msg.Source != "" {
		output += ":" + msg.Source + " "
	}

	// Add Command
	output += msg.Command

	// Add Parameters
	for _, param := range msg.Params {
		output += " " + param
	}

	// Add Closing Characters
	output += "\r\n"

	return output, nil
}
