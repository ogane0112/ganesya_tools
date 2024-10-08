package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
	"path/filepath"

    "github.com/joho/godotenv"
    "github.com/spf13/cobra"
)

var save bool
var filePath string

// notionCmd represents the notion command
var notionCmd = &cobra.Command{
    Use:   "notion",
    Short: "ローカルのMarkdownファイルをNotionに保存またはNotionの変更をローカルに反映します",
    Long: `指定したMarkdownファイルの内容をNotionデータベースに保存するか、
Notionの変更内容をローカルファイルに反映します。
.envファイルからAPIトークンとデータベースIDを読み込みます。`,
    Run: func(cmd *cobra.Command, args []string) {
        // .envファイルをロード
        err := godotenv.Load()
        if err != nil {
            log.Fatal("Error loading .env file")
        }

        // 環境変数を取得
        notionToken := os.Getenv("NOTION_TOKEN")
        databaseID := os.Getenv("NOTION_DATABASE_ID")

        if notionToken == "" || databaseID == "" {
            log.Fatal("NOTION_TOKEN or NOTION_DATABASE_ID is not set in environment variables.")
        }

        if save {
            // Check if filePath is provided
            if filePath == "" {
                log.Fatal("File path must be specified with -p flag.")
            }

            // Markdownファイルを読み込む
            content, err := ioutil.ReadFile(filePath)
			//ファイル名のみを取得する処理
			fileName := filepath.Base(filePath)
            if err != nil {
                log.Fatalf("Error reading file %s: %v", filePath, err)
            }

            // Notion APIリクエスト用のペイロード作成
            payload := map[string]interface{}{
                "parent": map[string]string{
                    "database_id": databaseID,
                },
                "properties": map[string]interface{}{
                    "title": []map[string]interface{}{
                        {
                            "text": map[string]string{
                                "content": fileName, // Inline page title
                            },
                        },
                    },
                },
                "children": []map[string]interface{}{
                    {
                        "object": "block",
                        "type":   "paragraph",
                        "paragraph": map[string]interface{}{
                            "rich_text": []map[string]interface{}{
                                {
                                    "type": "text",
                                    "text": map[string]string{
                                        "content": string(content), // Markdownの内容を保存
                                    },
                                },
                            },
                        },
                    },
                },
            }

            payloadBytes, err := json.Marshal(payload)
            if err != nil {
                log.Fatal("Error marshaling JSON:", err)
            }

            url := fmt.Sprintf("https://api.notion.com/v1/pages")
            req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
            if err != nil {
                log.Fatal("Error creating request:", err)
            }

            req.Header.Set("Authorization", "Bearer "+notionToken)
            req.Header.Set("Content-Type", "application/json")
            req.Header.Set("Notion-Version", "2022-06-28")

            client := &http.Client{}
            resp, err := client.Do(req)
            if err != nil {
                log.Fatal("Error making request:", err)
            }
            defer resp.Body.Close()

            if resp.StatusCode != http.StatusOK {
                body, _ := ioutil.ReadAll(resp.Body)
                log.Fatalf("Failed to create page: %s\nResponse: %s", resp.Status, string(body))
            }

            fmt.Println("Markdown file successfully saved to Notion.")
        } else if len(args) > 0 && args[0] == "pull" {
            // Implement pull functionality here
            fmt.Println("Pulling changes from Notion...")
            
            // Example pull logic (should be replaced with actual implementation):
            // 1. Fetch data from Notion.
            // 2. Update local Markdown files with fetched data.
            
        } else {
            fmt.Println("Invalid command or missing flag. Use -s to save or 'pull' to update.")
        }
    },
}

func init() {
    rootCmd.AddCommand(notionCmd)

    // Define flags and configuration settings.
    notionCmd.Flags().BoolVarP(&save, "save", "s", false, "Save local Markdown file to Notion")
    notionCmd.Flags().StringVarP(&filePath, "path", "p", "", "Path to the Markdown file to save")
}

