package end

import "hcc/harp/lib/logger"

func loggerEnd() {
	_ = logger.FpLog.Close()
}
