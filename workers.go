package main

import (
	"log"
)

func (wrkr Worker) ttlDeleter() {
	err := wrkr.DB.DeleteTtl(wrkr.CTX)
	if err != nil {
		log.Println("TTLDeleter |Failed to delete TTL urls | Err: ", err.Error())
		return
	}
	log.Println("TTLDeleter | Deleted URLS")

}
