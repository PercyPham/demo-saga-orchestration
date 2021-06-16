package postgresql

import (
	"github.com/percypham/saga-go"
	"gorm.io/gorm"
	"services.shared/apperror"
)

func (r *repoImpl) CreateSaga(sagaInstance *saga.Saga) error {
	result := r.db.Create(sagaInstance)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "create Saga in db using gorm")
	}
	return nil
}

func (r *repoImpl) FindSagaByID(id string) *saga.Saga {
	sagaInstance := new(saga.Saga)
	result := r.db.Where("id = ?", id).First(sagaInstance)
	if result.Error != nil {
		return nil
	}
	return sagaInstance
}

func (r *repoImpl) UpdateSaga(sagaInstance *saga.Saga) error {
	result := r.db.Save(sagaInstance)
	if result.Error != nil {
		return apperror.WithLog(result.Error, "save Saga in db using gorm")
	}
	return nil
}

type ProcessedMessage struct {
	ID string `json:"id" gorm:"primaryKey"`
}

func (r *repoImpl) CheckIfMessageProcessed(id string) bool {
	pGorm := new(ProcessedMessage)
	result := r.db.Where("id = ?", id).First(pGorm)
	return result.Error == nil
}

func (r *repoImpl) RecordMessageAsProcessed(id string) error {
	result := r.db.Create(&ProcessedMessage{id})
	if result.Error != nil {
		return apperror.WithLog(result.Error, "create ProcessedMessage "+id+" in db using gorm")
	}
	return nil
}

func (r *repoImpl) BeginTransaction() saga.Transaction {
	txDB := r.db.Begin()
	txRepo := &repoImpl{txDB}
	return &repoTransaction{txDB, txRepo}
}

type repoTransaction struct {
	txDB *gorm.DB
	repo saga.Repo
}

func (tx *repoTransaction) Repo() saga.Repo {
	return tx.repo
}

func (tx *repoTransaction) RollbackTransaction() {
	tx.txDB.Rollback()
}

func (tx *repoTransaction) CommitTransaction() {
	tx.txDB.Commit()
}
