package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
)

// ColumnInfo 用于扫描SHOW COLUMNS的结果
type ColumnInfo struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default interface{}
	Extra   string
}

var MysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "执行MySQL数据库相关操作",
	Long:  `使用mysql命令可以执行查看数据库、表格、字段等操作，并可以执行自定义SQL命令。`,
	Run: func(cmd *cobra.Command, args []string) {
		opt, _ := cmd.Flags().GetString("opt")
		table, _ := cmd.Flags().GetString("table")
		address, _ := cmd.Flags().GetString("address")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		database, _ := cmd.Flags().GetString("database")
		raw, _ := cmd.Flags().GetString("raw")
		db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=%s",
			username, password, address, database, "Asia%2FShanghai"))
		NoErr(err)
		switch opt {
		case "list_tables":
			var tables []string
			db.Raw("SHOW TABLES").Pluck("Tables_in_"+table, &tables)
			fmt.Println(tables)
		case "list_databases":
			var databases []string
			db.Raw("SHOW DATABASES").Pluck("Database", &databases)
			fmt.Println(databases)
		case "list_columns":
			if table == "" {
				fmt.Println("请提供表名")
				return
			}
			var columns []struct {
				Field   string
				Type    string
				Null    string
				Key     string
				Default sql.NullString // 使用 sql.NullString 以正确处理 NULL 值
				Extra   string
			}
			// 使用Scan而不是Pluck来接收所有列信息
			result := db.Raw("SHOW COLUMNS FROM " + table).Scan(&columns)
			if result.Error != nil {
				fmt.Printf("查询出错: %v\n", result.Error)
				return
			}

			for _, col := range columns {
				fmt.Printf("字段名: %s, 类型: %s, 允许空值: %s, 键: %s, 默认值: %v, 额外信息: %s\n",
					col.Field, col.Type, col.Null, col.Key, col.Default, col.Extra)
			}
			// ...
		case "raw":
			db.Exec(raw)
		default:
			fmt.Println("未知的操作，请检查opt参数")
		}
	},
}

func init() {
	MysqlCmd.Flags().StringP("table", "t", "", "指定操作的表名")
	MysqlCmd.Flags().StringP("opt", "o", "", "指定操作类型：list_tables, list_databases, list_columns, raw")
	MysqlCmd.Flags().StringP("address", "a", "localhost:3306", "指定MySQL服务器的地址")
	MysqlCmd.Flags().StringP("username", "u", "root", "指定MySQL服务器的用户名")
	MysqlCmd.Flags().StringP("password", "p", "root", "指定MySQL服务器的密码")
	MysqlCmd.Flags().StringP("raw", "r", "", "指定要执行的原始SQL命令")
	MysqlCmd.Flags().StringP("database", "", "go_ori", "指定要执行的原始SQL命令")
}
