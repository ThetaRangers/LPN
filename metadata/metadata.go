package metadata

import (
	"SDCC/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type KeyMetadata struct {
	Key  string `gorm:"primaryKey"`
	Size uint64
}

type Db struct {
	Db *gorm.DB
}

func (m Db) PutKey(key, master string, size uint64) {
	m.Db.Transaction(func(tx *gorm.DB) error {
		meta := KeyMetadata{Key: key, Size: size}
		var res KeyMetadata
		err := tx.Model(&KeyMetadata{}).Where("key = ?", key).First(&res).Error
		if err == gorm.ErrRecordNotFound {
			tx.Create(&meta)
		} else if err != nil {
			return err
		} else {
			tx.Model(&KeyMetadata{}).Where("key = ?", key).Updates(meta)
		}

		return nil
	})
}

func (m Db) AppendKey(key string, size uint64) {
	m.Db.Transaction(func(tx *gorm.DB) error {
		meta := KeyMetadata{Key: key, Size: size}
		var res KeyMetadata
		err := tx.Model(&KeyMetadata{}).Where("key = ?", key).First(&res).Error
		if err == gorm.ErrRecordNotFound {
			tx.Create(&meta)
		} else if err != nil {
			return err
		} else {
			meta.Size = meta.Size + res.Size
			tx.Model(&KeyMetadata{}).Where("key = ?", key).Updates(meta)
		}

		return nil
	})
}

func (m Db) Del(key string) {
	err := m.Db.Where("key = ?", key).Delete(&KeyMetadata{}).Error
	if err == gorm.ErrRecordNotFound {
		return
	}
}

func (m Db) ValidateKey(key string) {
	m.Db.Model(&KeyMetadata{}).Where("key = ?", key).Update("valid", true)
}

func (m Db) GetKey(key string) KeyMetadata {
	var meta KeyMetadata
	err := m.Db.Where("key = ?", key).First(&meta).Error
	if err == gorm.ErrRecordNotFound {
		return meta
	}

	return meta
}

func (m Db) NodesKeys(master string) []KeyMetadata {
	var keys []KeyMetadata
	m.Db.Model(&KeyMetadata{}).Where("master = ?", master).Find(&keys)

	return keys
}

func (m Db) GetAllKeys() []KeyMetadata {
	var keys []KeyMetadata
	m.Db.Model(&KeyMetadata{}).Find(&keys)

	return keys
}

func (m Db) EvaluateOffloading(keyList []string) []KeyMetadata {
	var keys []KeyMetadata
	res := make([]KeyMetadata, 0)

	m.Db.Model(&KeyMetadata{}).Order("size desc").Find(&keys)

	for _, k := range keys {
		if utils.Contains(keyList, k.Key) {
			res = append(res, k)
		}
	}

	return res
}

func (m Db) GetUsedSpace() uint64 {
	var space uint64

	err := m.Db.Model(&KeyMetadata{}).Select("sum(size)").First(&space).Error
	if err != nil {
		return 0
	}

	return space
}

func GetSizeDb() Db {
	db, err := gorm.Open(sqlite.Open("size.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("failed to connect database")
	}

	mt := Db{Db: db}
	db.AutoMigrate(&KeyMetadata{})

	return mt
}
