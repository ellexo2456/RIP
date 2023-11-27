package repository

import (
	"RIpPeakBack/internal/app/ds"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (repository *Repository) GetAlpinistByID(id int) (*ds.Alpinist, error) {
	alpinist := &ds.Alpinist{}

	err := repository.db.Preload("Expeditions").First(alpinist, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return alpinist, nil
}

func (repository *Repository) GetActiveAlpinists() (*[]ds.Alpinist, error) {
	alpinists := &[]ds.Alpinist{}
	err := repository.db.Find(alpinists, "status = ?", "действует").Error

	if err != nil {
		return nil, err
	}

	return alpinists, nil
}

func (repository *Repository) FilterByCountry(country string) (*[]ds.Alpinist, error) {
	alpinists := &[]ds.Alpinist{}
	err := repository.db.Find(alpinists, "country = ?", country).Error
	if err != nil {
		return nil, err
	}

	return alpinists, nil
}

func (repository *Repository) AddAlpinist(alpinist ds.Alpinist) (uint, error) {
	result := repository.db.Create(&alpinist)

	if err := result.Error; err != nil {
		return 0, err
	}
	return alpinist.ID, nil
}

func (repository *Repository) UpdateAlpinist(alpinist ds.Alpinist) error {
	result := repository.db.Save(&alpinist)

	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (repository *Repository) AddExpedition(expedition ds.Expedition) (uint, error) {
	var existingExpedition ds.Expedition
	result := repository.db.Where("status = ?", expedition.Status).First(&existingExpedition)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := repository.db.Create(&expedition).Error
		if err != nil {
			return 0, err
		}
		return expedition.ID, nil
	} else {
		err := repository.db.Model(&existingExpedition).Association("Alpinists").Append(expedition.Alpinists)
		if err != nil {
			return 0, err
		}
		return existingExpedition.ID, nil
	}

}

func (repository *Repository) GetExpeditionById(id uint) (*ds.Expedition, error) {
	expedition := &ds.Expedition{}

	err := repository.db.First(expedition, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return expedition, nil
}

func (repository *Repository) UpdateExpedition(expedition ds.Expedition) error {
	result := repository.db.Save(&expedition)

	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (repository *Repository) FilterByStatus(status string) (*[]ds.Expedition, error) {
	expedition := &[]ds.Expedition{}
	err := repository.db.Preload("Usr", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, login")
	}).Find(expedition, "status = ?", status).Error
	if err != nil {
		return nil, err
	}

	return expedition, nil
}

func (repository *Repository) FilterByFormedTime(startTime string, endTime string) (*[]ds.Expedition, error) {
	expedition := &[]ds.Expedition{}
	err := repository.db.Preload("Usr", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, login")
	}).Find(expedition, "status != 'удалено' AND formed_at BETWEEN ? AND ?", startTime, endTime).Error
	if err != nil {
		return nil, err
	}

	return expedition, nil
}

func (repository *Repository) FilterByFormedTimeAndStatus(startTime string, endTime string, status string) (*[]ds.Expedition, error) {
	expedition := &[]ds.Expedition{}
	err := repository.db.Preload("Usr", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, login")
	}).Find(expedition, "status = ? AND formed_at BETWEEN ? AND ?", status, startTime, endTime).Error
	if err != nil {
		return nil, err
	}

	return expedition, nil
}

func (repository *Repository) GetExpeditions() (*[]ds.Expedition, error) {
	expeditions := &[]ds.Expedition{}

	err := repository.db.Preload("Usr", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, login")
	}).Find(expeditions, "status != 'удалено'").Error
	if err != nil {
		return nil, err
	}

	return expeditions, nil
}

func (repository *Repository) UpdateStatus(expedition ds.Expedition) error {
	//result := repository.db.Table("expeditions").Where("id = ?", expedition.ID).Updates(ds.Expedition{Status: expedition.Status, FormedAt: expedition.FormedAt, ClosedAt: expedition.ClosedAt})
	result := repository.db.Table("expeditions").Where("id = ?", expedition.ID).Updates(expedition)

	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetExpeditionByID(expeditionID int) (*ds.Expedition, error) {
	expedition := &ds.Expedition{}

	err := repository.db.Preload("Alpinists").First(expedition, "id = ?", expeditionID).Error
	if err != nil {
		return nil, err
	}

	return expedition, nil
}

func (repository *Repository) DeleteExpedition(expedition ds.Expedition) error {
	//for _, alpinist := range expedition.Alpinists {
	//	err := repository.db.Model(&expedition).Association("Alpinists").Delete(&alpinist)
	//	if err != nil {
	//		return err
	//	}
	//}

	err := repository.db.Updates(ds.Expedition{ID: expedition.ID, Status: expedition.Status, ClosedAt: expedition.ClosedAt}).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) GetDraft(userID int) (ds.Expedition, error) {
	var expedition ds.Expedition
	err := repository.db.First(&expedition, "user_id = ? AND status = ?", userID, ds.StatusDraft).Error

	if err != nil {
		return ds.Expedition{}, err
	}

	return expedition, nil
}
