package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type KeyMetadata struct {
	Key    string `gorm:"primaryKey"`
	Master string
	Valid  bool
}

type MetadataDb struct {
	Db *gorm.DB
}

func (m MetadataDb) PutKey(key, master string) {
	m.Db.Transaction(func(tx *gorm.DB) error {
		meta := KeyMetadata{Key: key, Master: master, Valid: true}
		var res interface{}
		err := tx.Model(&KeyMetadata{}).Where("key = ? AND master = ?", key, master).Find(res).Error
		if err == gorm.ErrRecordNotFound {
			tx.Create(&meta)
		} else if err != nil {
			return err
		}

		return nil
	})
}

func (m MetadataDb) ValidateKey(key string) {
	m.Db.Model(&KeyMetadata{}).Where("key = ?", key).Update("valid", true)
}

func (m MetadataDb) GetKey(key string) KeyMetadata {
	var meta KeyMetadata
	m.Db.Where("key = ?", key).First(&meta)

	return meta
}

func (m MetadataDb) Invalidate(master string) {
	m.Db.Model(&KeyMetadata{}).Where("master = ?", master).Update("valid", false)
}

func (m MetadataDb) GetInvalidKeys(master string) []KeyMetadata {
	var keys []KeyMetadata
	m.Db.Model(&KeyMetadata{}).Where("master = ? AND valid = false", master).Find(&keys)

	return keys
}

func (m MetadataDb) NodesKeys(master string) []KeyMetadata {
	var keys []KeyMetadata
	m.Db.Model(&KeyMetadata{}).Where("master = ?", master).Find(&keys)

	return keys
}

func GetDb() MetadataDb {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	mt := MetadataDb{Db: db}
	db.AutoMigrate(&KeyMetadata{})

	return mt
}

func main() {
	dat := GetDb()

	dat.PutKey("a", "192.168.1.1")
	dat.PutKey("a", "192.168.1.1")
	dat.ValidateKey("a")
	fmt.Println("First: ", dat.GetKey("a").Valid)

	dat.Invalidate("192.168.1.1")
	fmt.Println("Second: ", dat.GetKey("a").Valid)

	fmt.Println("Invalid keys: ", dat.GetInvalidKeys("192.168.1.1"))
}
