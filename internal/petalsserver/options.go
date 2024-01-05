package petalsserver

type RunOptions struct {
	PublicName    string   // public name to use for the server
	MaxDiskSize   int      // max disk size to use for the server
	NumBlocks     int      // number of blocks to use for the server
	AutoStart     bool     // if true, will start server automatically
	ModelName     string   // e.g. "petals-team/StableBeluga2"
	StopOnProcess []string // process names that will stop the server
	//Threshold     float64  // threshold for auto-starting server
}
