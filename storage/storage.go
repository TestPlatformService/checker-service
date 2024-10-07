package storage

import (
	"checker/storage/repo"
)

type Isorage interface {
	Check() repo.ICheckStorage
}
