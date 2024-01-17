package bunconnect

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xo/dburl"
)

func IAMOpener(loadOptions ...func(options *config.LoadOptions) error) OpenerFn {
	return func(driverName, dataSourceName string) (*sql.DB, error) {

		url, err := dburl.Parse(dataSourceName)
		if err != nil {
			return nil, err
		}

		cfg, err := config.LoadDefaultConfig(context.Background(), loadOptions...)
		if err != nil {
			return nil, err
		}

		authenticationToken, err := auth.BuildAuthToken(
			context.Background(), url.Host, cfg.Region, url.User.Username(), cfg.Credentials)
		if err != nil {
			return nil, err
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			url.Hostname(), url.Port(), url.User.Username(), authenticationToken, url.Path[1:],
		)
		for key, strings := range url.Query() {
			for _, value := range strings {
				dsn = fmt.Sprintf("%s %s=%s", dsn, key, value)
			}
		}

		return sql.Open(driverName, dsn)
	}
}
