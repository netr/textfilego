## TextFileGo

Text File Line by Line Rotation Library with INI State Management

---

Easy to use library to rotate through lines in a text files. You can pick up wherever you left off using INI files to store the line number pointer. Must use \*.txt files for now. No tests as of now. Can expand it/clean it up if anyone wants me to.

Works by loading all files in a directory into a struct.

```
var txts *textfiles.Files

// initialize your text files from directory
txts = &textfiles.Files{}
err := txts.Init("/my-text-files/")
if err != nil {
    log.Fatal(err)
    return
}

// will get next line for /my-text-files/filename.txt
fmt.Printf("test: %s\n", txts.Next("filename"))
```
