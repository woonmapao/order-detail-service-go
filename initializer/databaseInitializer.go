package initializer

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBInitializer() {

	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	fmt.Println(`
	 ***********************
	*  Services Started !!  *
	 ***********************
	 

 ██████╗ ██████╗ ██████╗ ███████╗██████╗    
 ██╔═══██╗██╔══██╗██╔══██╗██╔════╝██╔══██╗   
 ██║   ██║██████╔╝██║  ██║█████╗  ██████╔╝   
 ██║   ██║██╔══██╗██║  ██║██╔══╝  ██╔══██╗   
 ╚██████╔╝██║  ██║██████╔╝███████╗██║  ██║   
  ╚═════╝ ╚═╝  ╚═╝╚═════╝ ╚══════╝╚═╝  ╚═╝   
											 
 ██████╗ ███████╗████████╗ █████╗ ██╗██╗     
 ██╔══██╗██╔════╝╚══██╔══╝██╔══██╗██║██║     
 ██║  ██║█████╗     ██║   ███████║██║██║     
 ██║  ██║██╔══╝     ██║   ██╔══██║██║██║     
 ██████╔╝███████╗   ██║   ██║  ██║██║███████╗
 ╚═════╝ ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝
 
							  
	 
	`)

	var now time.Time
	db.Raw("SELECT NOW()").Scan(&now)

	fmt.Println(now)

	DB = db
}
