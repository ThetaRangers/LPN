package metadata

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

type Keys struct {
	Key string `gorm:"primaryKey"`
}

type DbKeys struct {
	Db *gorm.DB
}

func (k DbKeys) PutKey(key string) {
	k.Db.Transaction(func(tx *gorm.DB) error {
		meta := Keys{Key: key}
		var res Keys
		err := tx.Model(&Keys{}).Where("key = ?", key).First(&res).Error
		if err == gorm.ErrRecordNotFound {
			tx.Create(&meta)
		} else if err != nil {
			return err
		} else {
			tx.Model(&Keys{}).Where("key = ?", key).Updates(meta)
		}

		return nil
	})
}

func (k DbKeys) DelKey(key string) {
	err := k.Db.Where("key = ?", key).Delete(&Keys{}).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
}

func (k DbKeys) GetKeys() []string {
	var keyStrings = make([]string, 0)
	var keys []Keys
	k.Db.Model(&Keys{}).Find(&keys)

	for _, k := range keys {
		keyStrings = append(keyStrings, k.Key)
	}

	return keyStrings
}

func GetKeyDb() DbKeys {
	db, err := gorm.Open(sqlite.Open("key.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}

	mt := DbKeys{Db: db}
	db.AutoMigrate(&Keys{})

	return mt
}
