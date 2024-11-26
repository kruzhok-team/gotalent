package talentv2

import (
	"context"

	"github.com/ogen-go/ogen/ogenerrors"
)

type ClientCredentialsSource interface {
	Credentials(context.Context) (clientID, clientSecret string, err error)
}

type TalentTokenSource interface {
	Token(context.Context) (string, error)
}

type clientCredentials struct {
	id     string
	secret string
}

// Credentials implements ClientCredentialsSource.
func (c *clientCredentials) Credentials(context.Context) (clientID string, clientSecret string, err error) {
	return c.id, c.secret, nil
}

func ClientCredentialsSecurity(clientID, clientSecret string) SecuritySource {
	return Security(&clientCredentials{
		id:     clientID,
		secret: clientSecret,
	}, nil)
}

func Security(ccsrc ClientCredentialsSource, tauth TalentTokenSource) SecuritySource {
	return &securitySource{
		ccsrc: ccsrc,
		tauth: tauth,
	}
}

type securitySource struct {
	ccsrc ClientCredentialsSource
	tauth TalentTokenSource
}

// ClientCredentials implements SecuritySource.
func (s *securitySource) ClientCredentials(ctx context.Context, operationName string) (ClientCredentials, error) {
	var cc ClientCredentials
	if s.ccsrc == nil {
		return cc, ogenerrors.ErrSkipClientSecurity
	}
	uname, pwd, err := s.ccsrc.Credentials(ctx)
	if err != nil {
		return cc, err
	}
	cc.Username = uname
	cc.Password = pwd
	return cc, nil
}

// TalentOAuth implements SecuritySource.
func (s *securitySource) TalentOAuth(ctx context.Context, operationName string) (TalentOAuth, error) {
	var auth TalentOAuth
	if s.tauth == nil {
		return auth, ogenerrors.ErrSkipClientSecurity
	}
	token, err := s.tauth.Token(ctx)
	if err != nil {
		return auth, err
	}
	auth.Token = token
	return auth, nil
}

var _ SecuritySource = (*securitySource)(nil)
