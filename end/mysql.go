package end

import "hcc/harp/lib/mysql"

func mysqlEnd() {
	if mysql.Db != nil {
		_ = mysql.Db.Close()
	}
}
