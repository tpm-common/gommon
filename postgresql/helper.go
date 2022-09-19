package postgresql

import "context"

func (pg PostgreSQL) QueryAll(ctx context.Context, dst interface{}, sql string, args ...interface{}) error {
	rows, err := pg.Query(ctx, sql, args...)
	if err == nil {
		defer rows.Close()

		err = pg.ScanAll(dst, rows)
	}

	if err != nil {
		return err
	}

	return nil
}

func (pg PostgreSQL) QueryOne(ctx context.Context, dst interface{}, sql string, args ...interface{}) error {
	rows, err := pg.Query(ctx, sql, args...)
	if err == nil {
		defer rows.Close()

		err = pg.ScanOne(dst, rows)
	}

	if err != nil {
		return err
	}

	return nil
}
