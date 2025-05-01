package thera

type Thera struct {
	db  *DB
	cfg Config
}

func NewThera(db *DB, cfg Config) *Thera {
	d := &Thera{
		db:  db,
		cfg: cfg,
	}
	return d
}

func (th *Thera) Start() error {
	return nil
}
