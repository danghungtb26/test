package directive

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"

	generatedAdmin "git.ctisoftware.vn/back-end/base/src/graph/generated/admin"
	network "git.ctisoftware.vn/back-end/base/src/network"
	"git.ctisoftware.vn/numerology/proto-lib/golang/authenticator"
	"git.ctisoftware.vn/numerology/proto-lib/grpc_client"
)

const (
	AccountTypeRoot      = 1
	AccountTypeSuperUser = 2
	AccountTypeUser      = 3
)

var AdminDirective = generatedAdmin.DirectiveRoot{
	RequiredAuthAdmin: func(ctx context.Context, obj interface{}, next graphql.Resolver, action []*string, actionAdmin *string, check_ip *bool) (res interface{}, err error) {
		if !network.HasToken(ctx) {
			return nil, fmt.Errorf("unauthorized")
		}

		tokenStr := network.Token(ctx)
		result, err := grpc_client.AuthenticatorClient().TokenVerify(ctx, &authenticator.TokenVerifyRequest{JwtToken: tokenStr})
		if err != nil || result == nil {
			return nil, err
		}

		if result.AccountType != AccountTypeRoot {
			return nil, fmt.Errorf("permission deny")
		}

		ctx = context.WithValue(ctx, "workspace_id", result.WorkspaceId)
		ctx = context.WithValue(ctx, "account_id", result.User.Account.Id)
		ctx = context.WithValue(ctx, "user_id", result.User.Id)
		ctx = context.WithValue(ctx, "email", result.User.Account.Email)
		ctx = context.WithValue(ctx, "sub_workspace_id", result.SubWorkspaceId)

		return next(ctx)
	},
	RequiredAuthSuperUser: func(ctx context.Context, obj interface{}, next graphql.Resolver, action *string) (res interface{}, err error) {
		if !network.HasToken(ctx) {
			return nil, fmt.Errorf("unauthorized")
		}

		tokenStr := network.Token(ctx)
		result, err := grpc_client.AuthenticatorClient().TokenVerify(ctx, &authenticator.TokenVerifyRequest{JwtToken: tokenStr})
		if err != nil || result == nil {
			return nil, err
		}

		if result.AccountType != AccountTypeSuperUser && result.AccountType != AccountTypeRoot {
			return nil, fmt.Errorf("permission deny")
		}

		ctx = context.WithValue(ctx, "workspace_id", result.WorkspaceId)
		ctx = context.WithValue(ctx, "account_id", result.User.Account.Id)
		ctx = context.WithValue(ctx, "user_id", result.User.Id)
		ctx = context.WithValue(ctx, "email", result.User.Account.Email)

		return next(ctx)
	},
}
