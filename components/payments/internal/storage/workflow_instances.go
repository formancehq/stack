package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type instance struct {
	bun.BaseModel `bun:"table:workflows_instances"`

	// Mandatory fields
	ID          string             `bun:"id,pk,type:text,notnull"`
	ScheduleID  string             `bun:"schedule_id,pk,type:text,notnull"`
	ConnectorID models.ConnectorID `bun:"connector_id,pk,type:character varying,notnull"`
	CreatedAt   time.Time          `bun:"created_at,type:timestamp without time zone,notnull"`
	UpdatedAt   time.Time          `bun:"updated_at,type:timestamp without time zone,notnull"`

	// Optional fields with default
	// c.f. https://bun.uptrace.dev/guide/models.html#default
	Terminated bool `bun:"terminated,type:boolean,notnull,nullzero,default:false"`

	// Optional fields
	// c.f.: https://bun.uptrace.dev/guide/models.html#nulls
	TerminatedAt *time.Time `bun:"terminated_at,type:timestamp without time zone,nullzero"`
	Error        *string    `bun:"error,type:text,nullzero"`
}

func (s *store) InstancesUpsert(ctx context.Context, instance models.Instance) error {
	toInsert := fromInstanceModel(instance)

	_, err := s.db.NewInsert().
		Model(&toInsert).
		On("CONFLICT (id, schedule_id, connector_id) DO NOTHING").
		Exec(ctx)

	return e("failed to insert new instance", err)
}

func (s *store) InstancesUpdate(ctx context.Context, instance models.Instance) error {
	toUpdate := fromInstanceModel(instance)

	_, err := s.db.NewUpdate().
		Model(&toUpdate).
		Set("updated_at = ?", instance.UpdatedAt).
		Set("terminated = ?", instance.Terminated).
		Set("terminated_at = ?", instance.TerminatedAt).
		Set("error = ?", instance.Error).
		WherePK().
		Exec(ctx)

	return e("failed to update instance", err)
}

func (s *store) InstancesDeleteFromConnectorID(ctx context.Context, connectorID models.ConnectorID) error {
	_, err := s.db.NewDelete().
		Model((*instance)(nil)).
		Where("connector_id = ?", connectorID).
		Exec(ctx)

	return e("failed to delete instances", err)
}

type InstanceQuery struct{}

type ListInstancesQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[InstanceQuery]]

func NewListInstancesQuery(opts bunpaginate.PaginatedQueryOptions[InstanceQuery]) ListInstancesQuery {
	return ListInstancesQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) instancesQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "schedule_id",
			key == "connector_id":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'connector_id' column can only be used with $match")
			}
			return fmt.Sprintf("%s = ?", key), []any{value}, nil
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *store) InstancesList(ctx context.Context, q ListInstancesQuery) (*bunpaginate.Cursor[models.Instance], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.instancesQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[InstanceQuery], instance](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[InstanceQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			if where != "" {
				query = query.Where(where, args...)
			}

			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch instances", err)
	}

	instances := make([]models.Instance, 0, len(cursor.Data))
	for _, i := range cursor.Data {
		instances = append(instances, toInstanceModel(i))
	}

	return &bunpaginate.Cursor[models.Instance]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     instances,
	}, nil
}

func fromInstanceModel(from models.Instance) instance {
	return instance{
		ID:           from.ID,
		ScheduleID:   from.ScheduleID,
		ConnectorID:  from.ConnectorID,
		CreatedAt:    from.CreatedAt,
		UpdatedAt:    from.UpdatedAt,
		Terminated:   from.Terminated,
		TerminatedAt: from.TerminatedAt,
		Error:        from.Error,
	}
}

func toInstanceModel(from instance) models.Instance {
	return models.Instance{
		ID:           from.ID,
		ScheduleID:   from.ScheduleID,
		ConnectorID:  from.ConnectorID,
		CreatedAt:    from.CreatedAt,
		UpdatedAt:    from.UpdatedAt,
		Terminated:   from.Terminated,
		TerminatedAt: from.TerminatedAt,
		Error:        from.Error,
	}
}
