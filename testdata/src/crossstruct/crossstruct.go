package crossstruct

type DB struct{}

func (d *DB) query() {} // want `function "DB.query" is called by "Handler.handle" but declared before it \(stepdown rule\)`

type Handler struct {
	db *DB
}

func (h *Handler) handle() {
	h.db.query()
}
