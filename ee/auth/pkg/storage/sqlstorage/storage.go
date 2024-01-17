package sqlstorage

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/auth/pkg/oidc"
	"github.com/formancehq/auth/pkg/storage"
	"github.com/zitadel/oidc/v2/pkg/op"
)

func mapSqlError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return storage.ErrNotFound
	}
	return err
}

var _ oidc.Storage = (*Storage)(nil)

type Storage struct {
	db *bun.DB
}

func (s *Storage) CreateUser(ctx context.Context, user *auth.User) error {
	_, err := s.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (s *Storage) FindUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	user := &auth.User{}
	if err := s.db.NewSelect().Model(user).Where("email = ?", email).Scan(ctx); err != nil {
		return nil, err
	}
	return user, nil
}

func New(db *bun.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// AuthRequestByID implements the op.Storage interface
// it will be called after the Login UI redirects back to the OIDC endpoint
func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	request := &auth.AuthRequest{}
	if err := s.db.NewSelect().Model(request).Where("id = ?", id).Limit(1).Scan(ctx); err != nil {
		return nil, err
	}
	return request, nil
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	request := &auth.AuthRequest{}
	if err := s.db.NewSelect().Model(request).Where("code = ?", code).Limit(1).Scan(ctx); err != nil {
		return nil, err
	}
	return request, nil
}

func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	_, err := s.db.NewUpdate().
		Model(&auth.AuthRequest{}).
		Set("code = ?", code).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (s *Storage) SaveClient(ctx context.Context, client *auth.Client) error {
	_, err := s.db.NewInsert().Model(client).Exec(ctx)
	return err
}

func (s *Storage) SaveAuthRequest(ctx context.Context, request *auth.AuthRequest) error {
	_, err := s.db.NewInsert().Model(request).Exec(ctx)
	return err
}

func (s *Storage) FindAuthRequest(ctx context.Context, id string) (*auth.AuthRequest, error) {
	ret := &auth.AuthRequest{}
	err := s.db.NewSelect().
		Model(ret).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) FindAuthRequestByCode(ctx context.Context, code string) (*auth.AuthRequest, error) {
	ret := &auth.AuthRequest{}
	err := s.db.NewSelect().
		Model(ret).
		Where("code = ?", code).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) UpdateAuthRequest(ctx context.Context, request *auth.AuthRequest) error {
	_, err := s.db.NewUpdate().
		Model(request).
		Where("id = ?", request.ID).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) UpdateAuthRequestCode(ctx context.Context, id string, code string) error {
	_, err := s.db.NewUpdate().
		Model(&auth.AuthRequest{}).
		Where("id = ?", id).
		Set("code = ?", code).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	_, err := s.db.NewDelete().
		Model(&auth.AuthRequest{}).
		Where("id = ?", id).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) SaveRefreshToken(ctx context.Context, token *auth.RefreshToken) error {
	_, err := s.db.NewInsert().Model(token).Exec(ctx)
	return err
}

func (s *Storage) FindRefreshToken(ctx context.Context, token string) (*auth.RefreshToken, error) {
	ret := &auth.RefreshToken{}
	err := s.db.NewSelect().
		Model(ret).
		Where("id = ?", token).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := s.db.NewDelete().
		Model(&auth.RefreshToken{}).
		Where("id = ?", token).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) SaveAccessToken(ctx context.Context, token *auth.AccessToken) error {
	_, err := s.db.NewInsert().Model(token).Exec(ctx)
	return err
}

func (s *Storage) FindAccessToken(ctx context.Context, token string) (*auth.AccessToken, error) {
	ret := &auth.AccessToken{}
	err := s.db.NewSelect().
		Model(ret).
		Where("id = ?", token).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) DeleteAccessToken(ctx context.Context, token string) error {
	_, err := s.db.NewDelete().
		Model(&auth.AccessToken{}).
		Where("id = ?", token).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) DeleteAccessTokensForUserAndClient(ctx context.Context, userID string, clientID string) error {
	_, err := s.db.NewDelete().
		Model(&auth.AccessToken{}).
		Where("user_id = ? and application_id = ?", userID, clientID).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) DeleteAccessTokensByRefreshToken(ctx context.Context, token string) error {
	_, err := s.db.NewDelete().
		Model(&auth.AccessToken{}).
		Where("refresh_token_id = ?", token).
		Exec(ctx)
	return mapSqlError(err)
}

func (s *Storage) FindUser(ctx context.Context, id string) (*auth.User, error) {
	ret := &auth.User{}
	err := s.db.NewSelect().
		Model(ret).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) FindClient(ctx context.Context, id string) (*auth.Client, error) {
	ret := &auth.Client{}
	err := s.db.NewSelect().
		Model(ret).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}

func (s *Storage) SaveUser(ctx context.Context, user *auth.User) error {
	_, err := s.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (s *Storage) FindUserBySubject(ctx context.Context, subject string) (*auth.User, error) {
	ret := &auth.User{}
	err := s.db.NewSelect().
		Model(ret).
		Where("subject = ?", subject).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, mapSqlError(err)
	}
	return ret, nil
}
