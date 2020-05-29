package storage

import (
	"github.com/jinzhu/gorm"
	"strings"
	// sqlite dialect for the gorm package
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	//	bot "github.com/pconcepcion/telegram_dice_bot/bot"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{DisableLevelTruncation: true, FullTimestamp: true, PadLevelText: true}
	log.Level = logrus.DebugLevel

}

// SQLiteStorage handles the storage on a SQLite database
type SQLiteStorage struct {
	db               *gorm.DB
	dbPath           string
	accessConnection string
}

// Connect opens and configures the sqlite DB using the recived configuration string
func Connect(accessConnection string) *SQLiteStorage {
	var err error
	var s *SQLiteStorage
	// TODO handle the schema and validate it
	s = &SQLiteStorage{accessConnection: accessConnection}
	log.Debug("AccessConnection: ", accessConnection)
	s.dbPath = strings.TrimPrefix(accessConnection, "sqlite://")
	log.Debug("DB Path: ", s.dbPath)
	// TODO validate the path
	s.db, err = gorm.Open("sqlite3", s.dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	log.Infof("Connected to the DB: %v on %s", s.db, s.dbPath)

	// Migrate the schema
	s.db.AutoMigrate(&Roll{})
	s.db.AutoMigrate(&Session{})
	s.db.AutoMigrate(&Player{})
	s.db.AutoMigrate(&Character{})

	// Create Indexes
	s.db.Model(&Player{}).AddIndex("idx_player_name", "name")
	s.db.Model(&Character{}).AddIndex("idx_character_name", "name")

	// Create
	//db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	//var product Product
	//db.First(&product, 1) // find product with id 1
	//db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	//db.Delete(&product)
	return s
}

// Close closes the SQLite database
func (s *SQLiteStorage) Close() {
	s.db.Close()
}
