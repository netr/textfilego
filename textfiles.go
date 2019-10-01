package textfiles

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/ini.v1"
)

// Files contains all files in the files directory with their rows and pointers
type Files struct {
	ini       *ini.File
	TextFiles map[string]*textFile
	lock      *sync.Mutex
	path      string
}

type textFile struct {
	Rows    []string
	pointer int
}

// Init must be called to load your ini file
func (f *Files) Init(path string) error {

	f.path = "textfiles.ini"
	f.initializeLock()

	if !f.doesIniExist() {
		f.createEmptyIniFile()
		f.Init(path)
		return nil
	}

	f.loadIniFile()
	err := f.putAllFilesInStorage(path)
	if err != nil {
		return err
	}

	return nil
}

// Next returns the next item in the text file mapping
func (f *Files) Next(filename string) string {
	f.lock.Lock()
	defer f.lock.Unlock()

	row := f.getCurrentLine(filename)
	return row
}

// Count returns the total number of lines in the text file
func (f *Files) Count(filename string) int {
	if _, ok := f.TextFiles[filename]; !ok {
		return 0
	}
	return len(f.TextFiles[filename].Rows)
}

// Shuffle rows in specific file
func (f *Files) Shuffle(filename string) {
	for i := len(f.TextFiles[filename].Rows) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		f.TextFiles[filename].Rows[i], f.TextFiles[filename].Rows[j] = f.TextFiles[filename].Rows[j], f.TextFiles[filename].Rows[i]
	}
}

// Reverse rows in specific file
func (f *Files) Reverse(filename string) {
	for i := len(f.TextFiles[filename].Rows)/2 - 1; i >= 0; i-- {
		opp := len(f.TextFiles[filename].Rows) - 1 - i
		f.TextFiles[filename].Rows[i], f.TextFiles[filename].Rows[opp] = f.TextFiles[filename].Rows[opp], f.TextFiles[filename].Rows[i]
	}
}

func (f *Files) getCurrentLine(filename string) string {
	if _, ok := f.TextFiles[filename]; !ok {
		return ""
	}
	line := f.TextFiles[filename].Rows[f.TextFiles[filename].pointer]
	f.incrementPointer(filename)
	return line
}

func (f *Files) newTextFile(filename string) *textFile {
	return &textFile{
		pointer: f.getPointerFromIni(filename),
	}
}

func (f *Files) putAllFilesInStorage(path string) error {

	if f.ini == nil {
		return errors.New("Must have valid INI file (Try running Init() first)")
	}

	f.TextFiles = make(map[string]*textFile)

	// grab all files in folder
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// rotate through all files
	for _, fi := range files {
		if !isTextFile(fi.Name()) {
			continue
		}
		f.putTextFileLinesInStorage(path, fi.Name())
	}

	return nil

}

func (f *Files) putTextFileLinesInStorage(path, filename string) {

	// extract filename and create a new textfile struct
	raw := getRawFileName(filename)
	f.TextFiles[raw] = f.newTextFile(raw)

	file, err := os.Open(path + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.TextFiles[raw].Rows = append(f.TextFiles[raw].Rows, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (f *Files) incrementPointer(filename string) {
	f.TextFiles[filename].pointer++
	if f.TextFiles[filename].pointer >= len(f.TextFiles[filename].Rows) {
		f.TextFiles[filename].pointer = 0
	}

	fmt.Printf("%s - pointer: %v -- rows: %v\n\n", filename, f.TextFiles[filename].pointer, len(f.TextFiles[filename].Rows))

	f.storePointer(filename)
}

func (f *Files) storePointer(filename string) {
	f.ini.Section("files").Key(filename).SetValue(fmt.Sprintf("%v", f.TextFiles[filename].pointer))
	f.ini.SaveTo(f.path)
}

func (f *Files) getPointerFromIni(filename string) int {
	pointer := f.ini.Section("files").Key(filename).String()
	if pointer == "" {
		f.ini.Section("files").Key(filename).SetValue("0")
		return 0
	}

	p, err := strconv.ParseInt(pointer, 10, 64)
	if err != nil {
		return 0
	}

	return int(p)
}

func (f *Files) initializeLock() {
	f.lock = &sync.Mutex{}
}

func (f *Files) createEmptyIniFile() {
	emptyFile, err := os.Create(f.path)
	if err != nil {
		log.Fatal(err)
	}
	emptyFile.Close()
}

func (f *Files) doesIniExist() bool {
	_, err := os.Stat(f.path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *Files) loadIniFile() {
	cfg, err := ini.Load(f.path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	f.ini = cfg
}

func isTextFile(filename string) bool {
	if !strings.Contains(filename, ".txt") {
		return false
	}
	return true
}

func getRawFileName(filename string) string {
	return strings.Replace(filename, ".txt", "", -1)
}