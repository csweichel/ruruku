package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        lne := scanner.Text()

        if strings.Contains(lne, "`json:") {
            lne = addYamlTag(lne)
        }

        fmt.Println(lne)
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
}

func addYamlTag(line string) string {
    ts := strings.Index(line, "`")
    te := strings.LastIndex(line, "`")
    jsonTag := line[ts+1:te]
    yamlTag := strings.Replace(jsonTag, "json", "yaml", 1)
    newTag := fmt.Sprintf("%s %s", jsonTag, yamlTag)
    return strings.Replace(line, jsonTag, newTag, 1)
}
