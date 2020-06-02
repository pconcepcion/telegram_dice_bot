package storage

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

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

//BaseModel is a base to replace gorm base model to use UUIDs instead of ints
type BaseModel struct {
	ID        uuid.UUID `valid:"uuidv4" gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// BeforeCreate will set a UUID rather than numeric ID for the primary key
func (bm *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.New()
	return scope.SetColumn("ID", uuid)
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
	// Add logs for the schema changes
	s.db.LogMode(true)
	// Migrate the schema
	s.db.AutoMigrate(&Roll{})
	s.db.AutoMigrate(&Session{})
	s.db.AutoMigrate(&Player{})
	s.db.AutoMigrate(&Character{})

	// Create Indexes
	// Player
	s.db.Model(&Player{}).AddIndex("idx_player_name", "name")
	s.db.Model(&Player{}).AddIndex("idx_player_username", "username")
	// Session
	s.db.Model(&Session{}).AddIndex("idx_session_name", "name")
	// Character
	s.db.Model(&Character{}).AddIndex("idx_character_name", "name")
	// Roll
	s.db.Model(&Roll{}).AddIndex("idx_roll_description", "description")
	s.db.Model(&Roll{}).AddIndex("idx_roll_expression", "expression")
	s.db.Model(&Roll{}).AddIndex("idx_roll_total", "total")

	// Stop generation of logs
	s.db.LogMode(false)

	return s
}

// Close closes the SQLite database
func (s *SQLiteStorage) Close() {
	s.db.Close()
}
