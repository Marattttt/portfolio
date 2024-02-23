package storageutils

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/Marattttt/portfolio/portfolio_back/internal/applog"
	"github.com/jackc/pgx/v5"
)

// Scans and returns all rows
//
// Panics if the row data cannot be scanned into type T
func PGScanMany[T any](ctx context.Context, l applog.Logger, rows *pgx.Rows) []T {
	vals := make([]T, 0)

	for (*rows).Next() {
		var val T
		if err := (*rows).Scan(&val); err != nil {
			(*rows).Close()
			l.Error(
				ctx,
				applog.DB,
				"while scanning data into "+fmt.Sprintf("Could not decode data from row into %s", reflect.TypeOf(val)),
				err,
				slog.Any("columns received", getColumnNames(rows)))

			panic(reflect.TypeOf(val))
		}

		vals = append(vals, val)
	}

	return vals
}

// Closes the rows passed
func getColumnNames(rows *pgx.Rows) []string {
	(*rows).Close()
	fields := (*rows).FieldDescriptions()
	names := make([]string, len(fields))

	for _, f := range fields {
		names = append(names, f.Name)
	}
	return names
}
