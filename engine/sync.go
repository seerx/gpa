package engine

import "context"

func (e *Engine) Sync(beans ...interface{}) error {
	_, err := e.dialect.GetTables(e.provider.Executor(), context.Background())
	if err != nil {
		e.logger.Error(err, "get tables from database")
		return err
	}

	return nil
}
