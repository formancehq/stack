package v2

import (
	"github.com/formancehq/stack/libs/go-libs/time"
	"io"
	"net/http"
	"strconv"


	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/ledger/internal/storage/ledgerstore"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
)

func getPITOOTFilter(r *http.Request) (*ledgerstore.PITFilter, error) {
	pitString := r.URL.Query().Get("endTime")
	ootString := r.URL.Query().Get("startTime")

	pit := time.Now()
	oot := time.Time{}

	if pitString != "" {
		var err error
		pit, err = time.ParseTime(pitString)
		if err != nil {
			return nil, err
		}
	}

	if ootString != "" {
		var err error
		oot, err = time.ParseTime(ootString)
		if err != nil {
			return nil, err
		}
	}

	return &ledgerstore.PITFilter{
		PIT: &pit,
		OOT: &oot,
	}, nil
}

func getPITFilter(r *http.Request) (*ledgerstore.PITFilter, error) {
	pitString := r.URL.Query().Get("pit")

	pit := time.Now()

	if pitString != "" {
		var err error
		pit, err = time.ParseTime(pitString)
		if err != nil {
			return nil, err
		}
	}

	return &ledgerstore.PITFilter{
		PIT: &pit,
	}, nil
}

func getPITFilterWithVolumes(r *http.Request) (*ledgerstore.PITFilterWithVolumes, error) {
	pit, err := getPITFilter(r)
	if err != nil {
		return nil, err
	}
	return &ledgerstore.PITFilterWithVolumes{
		PITFilter:              *pit,
		ExpandVolumes:          collectionutils.Contains(r.URL.Query()["expand"], "volumes"),
		ExpandEffectiveVolumes: collectionutils.Contains(r.URL.Query()["expand"], "effectiveVolumes"),
	}, nil
}

func getFiltersForVolumes(r *http.Request) (*ledgerstore.FiltersForVolumes, error) {
	pit, err := getPITOOTFilter(r)
	if err != nil {
		return nil, err
	}

	useInsertionDate := sharedapi.QueryParamBool(r, "insertionDate")
	groupLvl := 0

	groupLvlStr := r.URL.Query().Get("groupBy")
	if groupLvlStr != "" {
		groupLvlInt, err := strconv.Atoi(groupLvlStr)
		if err != nil {
			return nil, err
		}
		if groupLvlInt > 0 {
			groupLvl = groupLvlInt
		}
	}
	return &ledgerstore.FiltersForVolumes{
		PITFilter:        *pit,
		UseInsertionDate: useInsertionDate,
		GroupLvl:         uint(groupLvl),
	}, nil
}

func getQueryBuilder(r *http.Request) (query.Builder, error) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		return query.ParseJSON(string(data))
	}
	return nil, nil
}

func getPaginatedQueryOptionsOfPITFilterWithVolumes(r *http.Request) (*ledgerstore.PaginatedQueryOptions[ledgerstore.PITFilterWithVolumes], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	pitFilter, err := getPITFilterWithVolumes(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := bunpaginate.GetPageSize(r)
	if err != nil {
		return nil, err
	}

	return pointer.For(ledgerstore.NewPaginatedQueryOptions(*pitFilter).
		WithQueryBuilder(qb).
		WithPageSize(pageSize)), nil
}

func getPaginatedQueryOptionsOfFiltersForVolumes(r *http.Request) (*ledgerstore.PaginatedQueryOptions[ledgerstore.FiltersForVolumes], error) {
	qb, err := getQueryBuilder(r)
	if err != nil {
		return nil, err
	}

	filtersForVolumes, err := getFiltersForVolumes(r)
	if err != nil {
		return nil, err
	}

	pageSize, err := bunpaginate.GetPageSize(r)
	if err != nil {
		return nil, err
	}

	return pointer.For(ledgerstore.NewPaginatedQueryOptions(*filtersForVolumes).
		WithPageSize(pageSize).
		WithQueryBuilder(qb)), nil
}
