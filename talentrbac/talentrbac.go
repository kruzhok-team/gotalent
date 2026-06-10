// Клиент для выполнения проверок в API TalentRBAC.
//
// Документация к API: http://talent.kruzhok.org/rbac/docs/
package talentrbac

import (
	"context"
	"fmt"
	"net/http"
)

//go:generate go run github.com/shagohead/gotools@v0.1.3 ogen -target ./ -package talentrbac -clean openapi.yaml

// Интерфейс клиента для проверок наличия доступа.
type PermissionChecker interface {
	// Получение идентификатора доступа по символьному ключу.
	PermissionID(ctx context.Context, permission string) (int32, error)

	// Проверка наличия доступа у пользователя Таланта.
	UserHasPermission(ctx context.Context, pid, uid int32) (bool, error)
}

const defaultServerURL = "http://talent-rbac-api:8080/rbac"

// NewExternalChecker создает PermissionChecker для запросов из локальной сети.
func NewLocalChecker(http *http.Client) PermissionChecker {
	client, err := NewClient(defaultServerURL, WithClient(http))
	if err != nil {
		panic(err)
	}
	return &permissionChecker{invoker: client}
}

// NewExternalChecker создает PermissionChecker для запросов из сети интернет.
// Реализация [http.RoundTripper] указанного http клиента должна добавлять авторизационный заголовок с токеном TalentOAuth.
func NewExternalChecker(http *http.Client, host string) PermissionChecker {
	client, err := NewClient(host, WithClient(http))
	if err != nil {
		panic(err)
	}
	return &permissionChecker{invoker: client}
}

type permissionChecker struct {
	invoker Invoker
}

// PermissionID implements PermissionChecker.
func (p *permissionChecker) PermissionID(ctx context.Context, permission string) (int32, error) {
	res, err := p.invoker.PermissionMeta(ctx, PermissionMetaParams{PermissionKey: permission})
	if err != nil {
		return 0, err
	}
	switch resp := res.(type) {
	case *PermissionMeta:
		return resp.ID, nil
	case *ObjectNotFound:
		err = fmt.Errorf("клиент вернул NotFound для доступа %q", permission)
	default:
		err = fmt.Errorf("неожиданный client.PermissionIDRes: %#v", resp)
	}
	return 0, err
}

// UserHasPermission implements PermissionChecker.
func (p *permissionChecker) UserHasPermission(ctx context.Context, pid, uid int32) (bool, error) {
	res, err := p.invoker.HasPermission(ctx, HasPermissionParams{
		PermissionID: pid,
		UserID:       uid,
	})
	if err != nil {
		return false, err
	}
	var has bool
	switch resp := res.(type) {
	case *HasPermissionOK:
		has = true
	case *HasPermissionForbidden:
		has = false
	default:
		err = fmt.Errorf("unexpected client.HasPermissionRes: %#v", resp)
	}
	return has, err
}
