package talentv2

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/ogen-go/ogen/ogenerrors"
)

// Source of credentials for client_credentials OAuth grant.
type ClientCredentialsSource interface {
	Credentials(context.Context) (clientID, clientSecret string, err error)
}

// Source of Talent OAuth access_token.
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

// ClientCredentialsSecurity returns [SecuritySource] applicable only for client_credentials grant.
func ClientCredentialsSecurity(clientID, clientSecret string) SecuritySource {
	return Security(&clientCredentials{
		id:     clientID,
		secret: clientSecret,
	}, nil)
}

// Security returns [SecuritySource] which can be used as access_token source and/or client_credentials source.
// Both sources can be nil. If source is nil, appropriate security requirement will be ignored.
// If both of them nil, this SecuritySource can be used only for unauthenticated calls.
func Security(ccsrc ClientCredentialsSource, tauth TalentTokenSource) SecuritySource {
	return &securitySource{
		ccsrc: ccsrc,
		tauth: tauth,
	}
}

type ctxKey int

const (
	contextToken ctxKey = iota
)

// ContextTokenSource implements [TalentTokenSource] which uses token from context.
// Set token before making request with [ContextWithToken] function.
func ContextTokenSource() TalentTokenSource {
	return &contextTokenSource{}
}

type contextTokenSource struct{}

var ErrContextMissingToken = errors.New("context missing token")

// Token implements TalentTokenSource.
func (c *contextTokenSource) Token(ctx context.Context) (string, error) {
	if v := ctx.Value(contextToken); v != nil {
		return v.(string), nil
	}
	return "", ErrContextMissingToken
}

func ContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextToken, token)
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
