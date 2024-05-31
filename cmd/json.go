package cmd

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var JSONCmd = &cobra.Command{
	Use: "json",
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		validate, _ := cmd.Flags().GetBool("validate")
		format, _ := cmd.Flags().GetBool("format")
		minify, _ := cmd.Flags().GetBool("minify")
		escape, _ := cmd.Flags().GetBool("escape")
		unescape, _ := cmd.Flags().GetBool("unescape")
		get, _ := cmd.Flags().GetString("get")
		set, _ := cmd.Flags().GetString("set")
		remove, _ := cmd.Flags().GetString("remove")
		toYaml, _ := cmd.Flags().GetBool("to-yaml")
		toXml, _ := cmd.Flags().GetBool("to-xml")

		if validate {
			validateJSON(input)
		}
		if format {
			formatJSON(input)
		}
		if minify {
			minifyJSON(input)
		}
		if escape {
			escapeJSON(input)
		}
		if unescape {
			unescapeJSON(input)
		}
		if get != "" {
			getJSONValue(input, get)
		}
		if set != "" {
			setJSONValue(input, set)
		}
		if remove != "" {
			removeJSONField(input, remove)
		}
		if toYaml {
			convertToYAML(input)
		}
		if toXml {
			convertToXML(input)
		}
	},
}

func init() {
	JSONCmd.Flags().StringP("input", "i", "", "输入的JSON字符串")
	JSONCmd.Flags().BoolP("validate", "v", false, "验证JSON格式是否正确")
	JSONCmd.Flags().BoolP("format", "f", false, "格式化JSON字符串")
	JSONCmd.Flags().BoolP("minify", "m", false, "压缩JSON字符串")
	JSONCmd.Flags().BoolP("escape", "e", false, "对JSON字符串进行转义")
	JSONCmd.Flags().BoolP("unescape", "u", false, "去除JSON字符串的转义符号")
	JSONCmd.Flags().StringP("get", "g", "", "根据JSON路径获取指定字段的值")
	JSONCmd.Flags().StringP("set", "s", "", "根据JSON路径设置指定字段的值")
	JSONCmd.Flags().StringP("remove", "r", "", "根据JSON路径删除指定字段")
	JSONCmd.Flags().BoolP("to-yaml", "y", false, "将JSON转换为YAML格式")
	JSONCmd.Flags().BoolP("to-xml", "x", false, "将JSON转换为XML格式")
}

func validateJSON(input string) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
	} else {
		fmt.Println("JSON格式正确")
	}
}

func formatJSON(input string) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	formattedJSON, _ := json.MarshalIndent(jsonData, "", "  ")
	fmt.Println(string(formattedJSON))
}

func minifyJSON(input string) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	minifiedJSON, _ := json.Marshal(jsonData)
	fmt.Println(string(minifiedJSON))
}

func escapeJSON(input string) {
	escaped, _ := json.Marshal(input)
	fmt.Println(string(escaped))
}

func unescapeJSON(input string) {
	var unescaped string
	err := json.Unmarshal([]byte(input), &unescaped)
	if err != nil {
		fmt.Println("无效的JSON字符串")
		return
	}
	fmt.Println(unescaped)
}

func getJSONValue(input string, path string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	// 简单实现,仅支持单层路径
	value, exists := jsonData[path]
	if exists {
		fmt.Println(value)
	} else {
		fmt.Printf("字段 '%s' 不存在\n", path)
	}
}

func setJSONValue(input string, pathValue string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	// 简单实现,仅支持单层路径和字符串值
	path, value, found := strings.Cut(pathValue, "=")
	if !found {
		fmt.Println("无效的路径和值格式")
		return
	}
	jsonData[path] = value
	updatedJSON, _ := json.Marshal(jsonData)
	fmt.Println(string(updatedJSON))
}

func removeJSONField(input string, path string) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	// 简单实现,仅支持单层路径
	delete(jsonData, path)
	updatedJSON, _ := json.Marshal(jsonData)
	fmt.Println(string(updatedJSON))
}

func convertToYAML(input string) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	yamlData, _ := yaml.Marshal(jsonData)
	fmt.Println(string(yamlData))
}

func convertToXML(input string) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		fmt.Println("无效的JSON格式")
		return
	}
	xmlData, _ := xml.MarshalIndent(jsonData, "", "  ")
	fmt.Println(string(xmlData))
}
