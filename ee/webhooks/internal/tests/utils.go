package tests

import (
	"fmt"
	_ "fmt"
	_ "math/rand"
	"net/http"

	"github.com/formancehq/webhooks/internal/commons"
	"github.com/formancehq/webhooks/internal/migrations"
	"github.com/formancehq/webhooks/internal/services/httpclient"
	"github.com/formancehq/webhooks/internal/services/storage/postgres"

	"os"

	"database/sql"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/pgtesting"

	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
)


type FakeEnviroment struct {
	State *commons.State

}

const (
	endpointTest = "127.0.0.1:65445/test"
	payloadTest = "{'test':'true'}"
)

func StartPostgresServer(){
	if err := pgtesting.CreatePostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
}

func StopPostgresServer(){
	if err := pgtesting.DestroyPostgresServer(); err != nil {
		logging.Error(err)
		os.Exit(1)
	}	
}

func GetStoreProvider() (storage.PostgresStore, error) {
	
	postgres := storage.NewPostgresStoreProvider(nil)
	db, err := sql.Open("postgres", pgtesting.Server().GetDSN())
	if err != nil {
		logging.Error(err)
		os.Exit(1)
	}

	bunDB := bun.NewDB(db, pgdialect.New())

	err = migrations.Migrate(logging.TestingContext(), bunDB)

	if(err!= nil){
		return postgres,err
	}

	return storage.NewPostgresStoreProvider(bunDB), nil
}

func GetStoreProviderWithData(nbHooks int, nbAttempts int, t pgtesting.TestingT) (storage.PostgresStore, error) {
	
	provider, err := GetStoreProvider()
	if(err != nil){
		return provider, err
	}
	
	err = fillDBWithData(provider, nbHooks, nbAttempts)
	if(err != nil) {return provider, err}

	return provider, nil
	
}

func fillDBWithData(db storage.PostgresStore, nbHooks int, nbAttempts int) error {
	
	// hooks := make([]*commons.Hook, 0)
	// attempts := make([]*commons.Attempt, 0)

	// for i:=0; i < nbHooks ; i++ {
	// 	name := fmt.Sprintf("test-hook-%d", i)
	// 	hooks = append(hooks, commons.NewHook(name, 
	// 		[]string{fmt.Sprintf("test-event-%d", i)}, 
	// 		endpointTest, ""))
	// }

	// for i:=0; i < nbAttempts ; i++ {
	// 	hook := hooks[rand.Intn(nbHooks)]
	// 	id := hook.ID
	// 	name := hook.Name
	// 	attempts = append(attempts, commons.NewAttempt(
	// 		id, 
	// 		name, 
	// 		hook.Endpoint, 
	// 		endpointTest, 
	// 		payloadTest))
	// }

	// err := db.SaveHooks(hooks)
	// if(err!=nil){
	// 	return err
	// }

	// err = db.SaveAttempts(attempts)
	// if(err!=nil){
	// 	return err
	// }
	
	return nil
}


func NewHTTPServer(port int, routes ...[2]interface{}) *http.Server{
	mux := http.NewServeMux()

	for _, route := range routes {
        url, ok1 := route[0].(string)
        handler, ok2 := route[1].(http.HandlerFunc)
        if !ok1 || !ok2 {
            logging.Errorf("Invalid route format. Expected (string, http.HandlerFunc), got (%T, %T)", route[0], route[1])
        }
        mux.HandleFunc(url, handler)
    }
	
    server := &http.Server{
        Addr:    fmt.Sprintf(":%d", port),
        Handler: mux,
    }

    go func() {
       
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logging.Errorf("Could not listen on port %d: %v", port, err)
			os.Exit(1)
        }
		logging.Debugf("TestServer is running on %s", server.Addr)
    }()

    return server
}

func NewHTTPClient() *httpclient.DefaultHttpClient {
	client := httpclient.NewDefaultHttpClient(&http.Client{})
	return &client
}