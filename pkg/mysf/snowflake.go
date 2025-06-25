package mysf

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
)

var node *sf.Node

// Initialize Snowflake ID Generator.
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		zap.L().Error("Snowflake Init time.Parse failed:", zap.Error(err))
		return
	}

	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	zap.L().Info("Snowflake generator initialized successfully.")
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}
