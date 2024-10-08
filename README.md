# ganesya_tools

## 
```bash 
go mod tidy
```
go sumへチェックサムが生成される
つまるところ、go.sumは改竄されたモジュールを使用することを防ぐことができるためのファイルです。



## 永続的なフラグとローカルフラグの違いについて
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "app",
    Short: "アプリケーションの説明",
}

var subCmd = &cobra.Command{
    Use:   "sub",
    Short: "サブコマンドの説明",
    Run: func(cmd *cobra.Command, args []string) {
        // サブコマンドの実行時の処理
        fmt.Println("Subcommand executed")
    },
}

func main() {
    // 永続的なフラグの定義
    rootCmd.PersistentFlags().String("config", "", "設定ファイルのパス")

    // ローカルフラグの定義
    subCmd.Flags().BoolP("verbose", "v", false, "詳細出力を有効にする")

    // コマンドの追加
    rootCmd.AddCommand(subCmd)

    // コマンドの実行
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        return
    }
}
```
実装例
```bash
app --config=config.yaml #Persistent Flag)
app sub --config=config.yama #local Flag
```
実行例


# フラグの受け取り方
BoolVarP(&filepath,"path","p",false,"これは説明文でしゅ")