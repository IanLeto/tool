package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

var (
	url        string
	dbType     string
	host       string
	user       string
	dbPassword string
	database   string
	table      string
	size       int
)

var DBCmd = &cobra.Command{
	Use: "db",
	Run: func(cmd *cobra.Command, args []string) {
		// 构建数据库连接字符串
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, host, database)
		db, err := sql.Open(dbType, dsn)
		if err != nil {
			log.Fatalf("无法连接到数据库: %v", err)
		}
		defer db.Close()

		// 验证数据库连接
		err = db.Ping()
		if err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}
		log.Printf("成功连接到数据库 %s", database)

		// 模拟向数据库的指定表中写入数据的逻辑
		for i := 0; i < size; i++ {
			// 模拟数据插入，假设表的结构为(id, name, created_at)
			query := fmt.Sprintf("INSERT INTO %s (name, created_at) VALUES (?, ?)", table)
			_, err := db.Exec(query, fmt.Sprintf("Project-%d", i+1), time.Now())
			if err != nil {
				log.Printf("插入数据失败: %v", err)
			} else {
				log.Printf("成功插入第 %d 条数据", i+1)
			}
		}
	},
}

func init() {
	// 添加命令行参数
	DBCmd.Flags().StringVarP(&dbType, "type", "t", "mysql", "数据库类型 (例如: mysql)")
	DBCmd.Flags().StringVarP(&host, "host", "H", "localhost:3306", "数据库主机和端口")
	DBCmd.Flags().StringVarP(&user, "user", "u", "root", "数据库用户名")
	DBCmd.Flags().StringVarP(&dbPassword, "password", "p", "", "数据库密码")
	DBCmd.Flags().StringVarP(&database, "database", "d", "", "数据库名称")
	DBCmd.Flags().StringVarP(&table, "table", "T", "", "写入数据的表名")
	DBCmd.Flags().IntVarP(&size, "size", "s", 100, "插入数据的条数")
}
