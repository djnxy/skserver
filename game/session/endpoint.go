package session

import (
	"context"
	frame "nxy/testsocket/agent/service"

	"github.com/go-kit/kit/endpoint"
)

func makeTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(frame.Frame)
		switch req.Cmd {
		case "login":
			if len(req.Sessionid) == 1 {
				s.Login(req.Sessionid[0])
			}
			users, _ := s.ListUsers()
			if len(users) > 0 {
				req.Sessionid = users
				req.Data = users
			}
			return req, nil
		case "logout":
			if len(req.Sessionid) == 1 {
				s.Logout(req.Sessionid[0])
			}
			users, _ := s.ListUsers()
			if len(users) > 0 {
				req.Sessionid = users
				req.Data = users
			}
			return req, nil
		case "listusers":
			users, _ := s.ListUsers()
			req.Data = users
			return req, nil
		default:
			return nil, nil
		}
	}
}

//discard
func makeLoginTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(frame.Frame)
		if len(req.Sessionid) == 1 {
			s.Login(req.Sessionid[0])
		}
		return req, nil
	}
}

//discard
func makeListUsersTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(frame.Frame)
		users, _ := s.ListUsers()
		req.Data = users
		return req, nil
	}
}

//discard
func makeLogoutTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(frame.Frame)
		if len(req.Sessionid) == 1 {
			s.Logout(req.Sessionid[0])
		}
		return req, nil
	}
}
