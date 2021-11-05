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
	Size   uint64
}

type MetadataDb struct {
	Db *gorm.DB
}

func (m MetadataDb) PutKey(key, master string, size uint64) {
	m.Db.Transaction(func(tx *gorm.DB) error {
		meta := KeyMetadata{Key: key, Master: master, Valid: true, Size: size}
		var res interface{}
		err := tx.Model(&KeyMetadata{}).Where("key = ? AND master = ?", key, master).First(res).Error
		if err == gorm.ErrRecordNotFound {
			tx.Create(&meta)
		} else if err != nil {
			return err
		} else {
			tx.Model(&KeyMetadata{}).Where("key = ? AND master = ?", key, master).Updates(meta)
		}

		return nil
	})
}

func (m MetadataDb) ValidateKey(key string) {
	m.Db.Model(&KeyMetadata{}).Where("key = ?", key).Update("valid", true)
}

func (m MetadataDb) GetKey(key string) KeyMetadata {
	var meta KeyMetadata
	err := m.Db.Where("key = ?", key).First(&meta).Error
	if err == gorm.ErrRecordNotFound {
		return meta
	}

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

func (m MetadataDb) EvaluateOffloading(master string) []KeyMetadata {
	var keys []KeyMetadata
	m.Db.Model(&KeyMetadata{}).Where("master = ?", master).Order("size desc").Find(&keys)

	return keys
}

func (m MetadataDb) GetUsedSpace(master string) uint64 {
	var space uint64

	m.Db.Model(&KeyMetadata{}).Where("master = ?", master).Select("sum(size)").First(&space)

	return space
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

	dat.PutKey("a", "192.168.1.1", 20)
	dat.PutKey("a", "192.168.1.1", 40)
	dat.PutKey("b", "192.168.1.1", 60)
	dat.PutKey("c", "192.168.1.2", 60)
	dat.ValidateKey("a")

	fmt.Println(dat.GetUsedSpace("192.168.1.1"))
	fmt.Println("First: ", dat.GetKey("a").Valid)

	dat.Invalidate("192.168.1.1")
	fmt.Println("Second: ", dat.GetKey("a").Valid)

	fmt.Println("Invalid keys: ", dat.GetInvalidKeys("192.168.1.1"))
	fmt.Println("Nodes keys: ", dat.NodesKeys("192.168.1.1"))
	fmt.Println("Offloading: ", dat.EvaluateOffloading("192.168.1.1"))
}
