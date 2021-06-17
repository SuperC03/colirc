package data

// Represents a Client Connected to This Server or Elsewhere
type Client struct {
	Nickname string
	Username string
	Realname string
	Hostname string
}

// Stores All Clients Connected to THIS Server
var ConnectedClients = make([]*Client, 0)

func AddClient(client *Client) {
	ConnectedClients = append(ConnectedClients, client)
}

func FindClientByUsername(username string) *Client {
	for _, u := range ConnectedClients {
		if u.Username == username {
			return u
		}
	}
	return nil
}

func FindClientByNickname(nickname string) *Client {
	for _, u := range ConnectedClients {
		if u.Nickname == nickname {
			return u
		}
	}
	return nil
}

func DeleteClientByUsername(username string) bool {
	cs := make([]*Client, len(ConnectedClients))
	f := false
	for _, c := range ConnectedClients {
		if c.Username != username {
			cs = append(cs, c)
			continue
		}
		f = true
	}
	ConnectedClients = cs
	return f
}
