package files

import (
	"log"
	"time"
)

func (fm *FileManager) expiryWorker() {
	ticker := time.NewTicker(time.Second * 60)

	for {
		select {
		case <-fm.stop:
			ticker.Stop()
			return
		case <-ticker.C:
			expiredHashes, err := fm.db.GetExpiredFileHashes()
			if err != nil {
				log.Printf("error getting expired file hashes: %s\n", err)
				continue
			}

			for _, hash := range expiredHashes {
				err = deleteFile(hash)
				if err != nil {
					log.Printf("error deleting file from disk %s: %s\n", hash, err)
					continue
				}

				err = fm.db.RemoveHashFromFiles(hash)
				if err != nil {
					log.Printf("error deleting file from db %s: %s\n", hash, err)
					continue
				}
			}

			if len(expiredHashes) > 0 {
				log.Printf("deleted file with %d distinct hashes\n", len(expiredHashes))
			}
		}
	}
}
