package link

import (
	"go_purple/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	Database *db.Db
}

func NewLinkRepository(db *db.Db) *LinkRepository {
	return &LinkRepository{
		Database: db,
	}
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	res := repo.Database.Create(link)
	if res.Error != nil {
		return nil, res.Error
	}

	return link, nil
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	res := repo.Database.First(&link, "hash = ?", hash)

	if res.Error != nil {
		return nil, res.Error
	}

	return &link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	res := repo.Database.Clauses(clause.Returning{}).Updates(link)

	if res.Error != nil {
		return nil, res.Error
	}

	return link, nil
}

func (repo *LinkRepository) Delete(id uint) error {
	res := repo.Database.DB.Delete(&Link{}, id)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (repo *LinkRepository) GetById(id uint) (*Link, error) {
	var link Link
	res := repo.Database.DB.First(&link, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &link, nil
}

func (repo *LinkRepository) Count() int64 {
	var count int64
	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Count(&count)

	return count
}

func (repo *LinkRepository) GetAll(limit, offset uint) []Link {
	var links []Link

	repo.Database.
		Table("links").
		Where("deleted_at is null").
		Order("id asc").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&links)

	return links
}
