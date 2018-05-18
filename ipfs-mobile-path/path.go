package mobilePath

func InitMobilePath() error {
	err := initExternStorageFilePath()
	if err != nil {
		return err
	}
	//add other path....
	return nil
}

func GetExternStorageFilePath() string {
	return getExternStorageFilePath()
}

func GetStorageCachePath() string {
	return getStorageCachePath()
}
