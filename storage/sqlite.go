package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
	//	bot "github.com/pconcepcion/telegram_dice_bot/bot"
)

var log = logrus.New()

func init() {
	log.Formatter = &logrus.TextFormatter{DisableLevelTruncation: true, FullTimestamp: true, PadLevelText: true}
	log.Level = logrus.DebugLevel

}

type Roll struct {
	gorm.Model
	DiceExpression string
	Sesssion       Session
	Results        string
	FinalResult    uint32
}

var (
	db *gorm.DB
)

func main() {
	db, err := gorm.Open("sqlite3", "telagram_dice_bot.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Roll{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Player{})
	db.AutoMigrate(&Character{})

	// Create Indexes
	db.Model(&Player{}).AddIndex("idx_player_name", "name")
	db.Model(&Character{}).AddIndex("idx_character_name", "name")

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
}
