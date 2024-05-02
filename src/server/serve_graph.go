package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"git.ctisoftware.vn/back-end/base/src/graph/directive"
	generated_admin "git.ctisoftware.vn/back-end/base/src/graph/generated/admin"
	generated_user "git.ctisoftware.vn/back-end/base/src/graph/generated/user"
	resolver_admin "git.ctisoftware.vn/back-end/base/src/graph/resolver/admin"
	resolver_user "git.ctisoftware.vn/back-end/base/src/graph/resolver/user"
	"git.ctisoftware.vn/back-end/base/src/middleware"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
)

func ServeGraph(ctx context.Context, addr string) (err error) {
	defer log.Println("HTTP server stopped", err)

	r := chi.NewRouter()
	v1(r)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	srv := http.Server{
		Addr:    addr,
		Handler: r,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	errChan := make(chan error, 1)

	go func(ctx context.Context, errChan chan error) {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}(ctx, errChan)

	log.Printf("Listen and Serve Authenticator-Graph-Service API at: %s\n", addr)

	select {
	case <-ctx.Done():
		return nil
	case err = <-errChan:
		return err
	}
}

func v1(r chi.Router) {
	configAdmin := generated_admin.Config{Resolvers: &resolver_admin.Resolver{}}
	configAdmin.Directives = directive.AdminDirective

	configUser := generated_user.Config{Resolvers: &resolver_user.Resolver{}}
	configUser.Directives = directive.UserDirective

	srvAdmin := handler.NewDefaultServer(generated_admin.NewExecutableSchema(configAdmin))
	srvUser := handler.NewDefaultServer(generated_user.NewExecutableSchema(configUser))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AllowAll().Handler)
		r.With(middleware.Middleware()).Route("/graphql", func(r chi.Router) {
			r.Handle("/admin", srvAdmin)
			r.Handle("/user", srvUser)
		})
	})
}
