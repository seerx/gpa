package engine

import "time"

func (e *Engine) SetLogSQL(log bool) *Engine {
	e.GetProvider().SetLogSQL(log)
	return e
}

func (e *Engine) SetMaxIdleConns(n int) *Engine {
	e.db.SetMaxIdleConns(n)
	return e
}

func (e *Engine) SetMaxOpenConns(n int) *Engine {
	e.db.SetMaxOpenConns(n)
	return e
}

func (e *Engine) SetConnMaxIdleTime(d time.Duration) *Engine {
	e.db.SetConnMaxIdleTime(d)
	return e
}

func (e *Engine) SetConnMaxLifetime(d time.Duration) *Engine {
	e.db.SetConnMaxLifetime(d)
	return e
}
