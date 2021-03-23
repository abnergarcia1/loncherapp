package main

import (
	"context"
	"sync"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"
)

var manager *ProcessManager

type ProcessManager struct {
	wg sync.WaitGroup
	sync.Mutex
	Ctx context.Context
	DB  sql.StorageDB
}
