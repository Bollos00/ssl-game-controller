package rcon

import (
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslconn"
	"github.com/odeke-em/go-uuid"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
)

type TeamServer struct {
	ProcessTeamRequest func(teamName string, request TeamToController) error
	AllowedTeamNames   []string
	*Server
}

type TeamClient struct {
	*Client
}

func NewTeamServer(address string) (s *TeamServer) {
	s = new(TeamServer)
	s.ProcessTeamRequest = func(string, TeamToController) error { return nil }
	s.Server = NewServer(address)
	s.ConnectionHandler = s.handleClientConnection
	return
}

func (c *TeamClient) receiveRegistration(server *TeamServer) error {
	registration := TeamRegistration{}
	if err := sslconn.ReceiveMessage(c.Conn, &registration); err != nil {
		return err
	}

	if registration.TeamName == nil {
		return errors.New("No team name specified")
	}
	if !isAllowedTeamName(*registration.TeamName, server.AllowedTeamNames) {
		return errors.Errorf("Invalid team name: '%v'. Expecting one of these: %v", *registration.TeamName, server.AllowedTeamNames)
	}
	c.Id = *registration.TeamName
	if _, exists := server.Clients[c.Id]; exists {
		return errors.New("Team with given name already registered: " + c.Id)
	}
	c.PubKey = server.TrustedKeys[c.Id]
	if c.PubKey != nil {
		err := c.verifyRegistration(registration)
		if err != nil {
			return err
		}
	} else {
		c.Token = ""
	}

	c.reply(c.Ok())

	log.Printf("Team %v connected.", *registration.TeamName)

	return nil
}

func isAllowedTeamName(teamName string, allowed []string) bool {
	for _, name := range allowed {
		if name == teamName {
			return true
		}
	}
	return false
}

func (c *TeamClient) verifyRegistration(registration TeamRegistration) error {
	if registration.Signature == nil {
		return errors.New("Missing signature")
	}
	if registration.Signature.Token == nil || *registration.Signature.Token != c.Token {
		sendToken := ""
		if registration.Signature.Token != nil {
			sendToken = *registration.Signature.Token
		}
		return errors.Errorf("Client %v sent an invalid token: %v != %v", c.Id, sendToken, c.Token)
	}
	signature := registration.Signature.Pkcs1V15
	registration.Signature.Pkcs1V15 = []byte{}
	err := VerifySignature(c.PubKey, &registration, signature)
	registration.Signature.Pkcs1V15 = signature
	if err != nil {
		return errors.New("Invalid signature")
	}
	c.Token = uuid.New()
	return nil
}

func (c *TeamClient) verifyRequest(req TeamToController) error {
	if req.Signature == nil {
		return errors.New("Missing signature")
	}
	if req.Signature.Token == nil || *req.Signature.Token != c.Token {
		sendToken := ""
		if req.Signature.Token != nil {
			sendToken = *req.Signature.Token
		}
		return errors.Errorf("Invalid token: %v != %v", sendToken, c.Token)
	}
	signature := req.Signature.Pkcs1V15
	req.Signature.Pkcs1V15 = []byte{}
	err := VerifySignature(c.PubKey, &req, signature)
	req.Signature.Pkcs1V15 = signature
	if err != nil {
		return errors.Wrap(err, "Verification failed.")
	}
	c.Token = uuid.New()
	return nil
}

func (s *TeamServer) handleClientConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Could not close team client connection: %v", err)
		}
	}()

	client := TeamClient{Client: &Client{Conn: conn, Token: uuid.New()}}
	client.reply(client.Ok())

	err := client.receiveRegistration(s)
	if err != nil {
		client.reply(client.Reject(err.Error()))
		return
	}

	s.Clients[client.Id] = client.Client
	defer s.CloseConnection(client.Id)
	log.Printf("Client %v connected", client.Id)
	for _, observer := range s.ClientsChangedObservers {
		observer()
	}

	for {
		req := TeamToController{}
		if err := sslconn.ReceiveMessage(conn, &req); err != nil {
			if err == io.EOF {
				return
			}
			log.Print(err)
			continue
		}
		if client.PubKey != nil {
			if err := client.verifyRequest(req); err != nil {
				client.reply(client.Reject(err.Error()))
				continue
			}
		}
		if err := s.ProcessTeamRequest(client.Id, req); err != nil {
			client.reply(client.Reject(err.Error()))
		} else {
			client.reply(client.Ok())
		}
	}
}

func (s *TeamServer) SendRequest(teamName string, request ControllerToTeam) error {
	if client, ok := s.Clients[teamName]; ok {
		return client.SendRequest(request)
	}
	return errors.Errorf("Client '%v' not connected", teamName)
}

func (c *Client) SendRequest(request ControllerToTeam) error {
	return sslconn.SendMessage(c.Conn, &request)
}

func (c *Client) reply(reply ControllerReply) {
	msg := ControllerToTeam_ControllerReply{ControllerReply: &reply}
	response := ControllerToTeam{Msg: &msg}
	if err := sslconn.SendMessage(c.Conn, &response); err != nil {
		log.Print("Failed to send reply: ", err)
	}
}
