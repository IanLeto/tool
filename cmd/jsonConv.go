package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"text/template"
)

var JsonCmd = &cobra.Command{
	Use: "jsonconv",
	Run: func(cmd *cobra.Command, args []string) {
		jsonStr := `{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_age": "1d",
            "max_size": "50gb"
          },
          "set_priority": {
            "priority": 100
          }
        }
      }
    }
  }
}`

		removeNewlines, _ := cmd.Flags().GetBool("remove-newlines")
		addNewlines, _ := cmd.Flags().GetBool("add-newlines")
		prettyFormat, _ := cmd.Flags().GetBool("pretty-format")
		escape, _ := cmd.Flags().GetBool("escape")
		unescape, _ := cmd.Flags().GetBool("unescape")

		if removeNewlines {
			jsonStr = strings.ReplaceAll(jsonStr, "\n", "")
			jsonStr = strings.ReplaceAll(jsonStr, " ", "")
		}

		if addNewlines && !removeNewlines {
			jsonStr = strings.ReplaceAll(jsonStr, ",", ",\n")
		}

		if prettyFormat {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, []byte(jsonStr), "", "\t"); err != nil {
				fmt.Println("JSON parse error: ", err)
				return
			}
			jsonStr = prettyJSON.String()
		}

		if escape {
			t := template.New("t")
			t = template.Must(t.Parse(jsonStr))
			var tpl bytes.Buffer
			if err := t.Execute(&tpl, nil); err != nil {
				fmt.Println("Template execution error: ", err)
				return
			}
			jsonStr = tpl.String()
		}

		if unescape {
			jsonStr = strings.ReplaceAll(jsonStr, "\\u003c", "<")
			jsonStr = strings.ReplaceAll(jsonStr, "\\u003e", ">")
			jsonStr = strings.ReplaceAll(jsonStr, "\\u0026", "&")
		}

		fmt.Println(jsonStr)
	},
}

func init() {
	JsonCmd.Flags().StringP("input", "i", "", "输入json")
	JsonCmd.Flags().IntP("size", "", 100, "文件大小")
	JsonCmd.Flags().BoolP("remove-newlines", "r", false, "去除换行符")
	JsonCmd.Flags().BoolP("add-newlines", "a", false, "添加换行符")
	JsonCmd.Flags().BoolP("pretty-format", "p", false, "美化格式")
	JsonCmd.Flags().BoolP("escape", "e", false, "添加转义")
	JsonCmd.Flags().BoolP("unescape", "u", false, "去除转义")
}
