package postgresql

import (
	"gorm.io/gorm"
	"services.shared/apperror"
	"services.shared/saga"
	"services.shared/saga/msg"
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
	ID      string `json:"id"`
	Message []byte `json:"message" sql:"type:json"`
}

func (r *repoImpl) CreateProcessedMessage(message msg.Message) error {
	m, err := msg.Marshal(message)
	if err != nil {
		return apperror.WithLog(err, "marshal message to json")
	}
	processedMessage := &ProcessedMessage{
		ID:      message.ID(),
		Message: m,
	}
	result := r.db.Create(processedMessage)
	if result.Error != nil {
		return apperror.WithLog(err, "create processed message using gorm")
	}
	return nil
}

func (r *repoImpl) GetProcessedMessageByID(id string) msg.Message {
	processedMessage := new(ProcessedMessage)
	result := r.db.Where("id = ?", id).First(processedMessage)
	if result.Error != nil {
		return nil
	}
	messageRaw := processedMessage.Message
	message, err := msg.Unmarshal(messageRaw)
	if err != nil {
		return nil
	}
	return message
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
