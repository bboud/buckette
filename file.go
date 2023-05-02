package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"strings"
	"time"
)

const (
	FileStoreDir   = "./buckette-data/files/"
	RecordStoreDir = "./buckette-data/records/"
	TmpDir         = "./buckette-data/tmp/"
)

type FileContent []byte

type File struct {
	FileName      string
	UUID          [32]byte
	URL           string
	Size          int64
	ContentType   string
	Uploaded      time.Time
	UserUploaded  string
	DownloadCount int

	tmpHash  string
	uuidName string
}

func (f *File) copyToFileSystem(reader io.Reader) error {
	record, err := os.OpenFile(TmpDir+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+TmpDir+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return err
	}
	defer record.Close()

	file, err := os.OpenFile(TmpDir+"DAT_"+f.tmpHash, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		LogWarning(
			"Unable to create file "+TmpDir+"DAT_"+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		LogWarning(
			"Unable to copy data to file "+TmpDir+"DAT_"+f.tmpHash,
			"Handling upload part for "+f.tmpHash,
			err,
		)
		return err
	}

	//Record writer
	rData, err := json.Marshal(f)
	if err != nil {
		LogFatal(
			"Unable to marshal json for "+f.tmpHash,
			"Writing the record to the record store",
			err)
	}

	record.Write(rData)

	return nil
}

func (f *File) cleanTmp(success bool) error {

	if success {
		err := os.Rename(TmpDir+f.tmpHash, RecordStoreDir+f.uuidName)
		if err != nil {
			LogFatal(
				"Unable to move temporary file to file store",
				"Moving file from temp to final storage for "+f.tmpHash,
				err)
			return err
		}

		// We want to store the files using their hash for faster lookup on disk
		err = os.Rename(TmpDir+"DAT_"+f.tmpHash, FileStoreDir+f.uuidName)
		if err != nil {
			LogFatal(
				"Unable to move temporary file to file store",
				"Moving file from temp to final storage for "+f.tmpHash,
				err)
			return err
		}
		return nil
	} else {
		if err := os.Remove(TmpDir + f.tmpHash); err != nil {
			if !os.IsNotExist(err) {
				LogWarning(
					"Unable to clean up after "+f.tmpHash,
					"Cleaning the temporary file directory",
					err)
				return err
			}
		}
		if err := os.Remove(TmpDir + "DAT_" + f.tmpHash); err != nil {
			if !os.IsNotExist(err) {
				LogWarning(
					"Unable to clean up after "+f.tmpHash,
					"Cleaning the temporary file directory",
					err)
				return err
			}
		}
		return nil
	}
}

func encodeToString(v []byte) string {
	vName := base64.RawStdEncoding.EncodeToString(v)
	vName = strings.ReplaceAll(vName, "/", "*")
	return vName
}
