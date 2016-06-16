package setup

import (
	"os"
)

func setupLive() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Live(
		profile  			CHAR(50) PRIMARY KEY NOT NULL,
		head_url 			CHAR(50) NOT NULL,
		recheck_threshold   INT NOT NULL,
		ping_threshold		INT NOT NULL,
		head_threshold		INT NOT NULL,
		ping_address		CHAR(50) NOT NULL,
		ping_protocol		CHAR(10)
		);`)
	Database.Exec(`INSERT INTO Live VALUES (
	"default", "https://google.com", 30000, 4000, 4000,"8.8.8.8", "tcp");`)
}

func setupDisk() {
	Database.Exec(`CREATE TABLE IF NOT EXISTS Disk(
		profile				CHAR(50) PRIMARY KEY NOT NULL,
		space_w_limit		INT NOT NULL,
		space_d_limit		INT NOT NULL,
		inode_w_limit		INT NOT NULL,
		inode_d_limit		INT NOT NULL,
		recheck_threshold 	INT NOT NULL,
		disk				CHAR(500) NOT NULL
		);`)
	Database.Exec(`INSERT INTO Disk VALUES ("default", 2000, 1000, 2000, 1000, 5000,
			"/,` + os.Getenv("HOME") + `");`)
}

func setupDB() {
	setupLive()
	setupDisk()
	return
}
