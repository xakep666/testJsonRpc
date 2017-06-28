package dao

import "github.com/jinzhu/gorm"

func transactionalOperation(db *gorm.DB, action func (tx *gorm.DB) (*gorm.DB, error)) error {
	tx:=db.Begin()
	resp, err:=action(tx)
	if resp.Error!=nil {
		tx.Rollback()
		return resp.Error
	}
	if err!=nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
