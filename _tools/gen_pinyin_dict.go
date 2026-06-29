package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type cmdArgs struct {
	inputFile    string
	outputFile   string
	varName      string
	overlayFiles []string
}

type dictData struct {
	values map[string]string
	order  []string
}

func newDictData() *dictData {
	return &dictData{
		values: make(map[string]string),
		order:  make([]string, 0),
	}
}

func (d *dictData) set(hexCode string, pinyin string) {
	if _, ok := d.values[hexCode]; !ok {
		d.order = append(d.order, hexCode)
	}
	d.values[hexCode] = pinyin
}

func (d *dictData) merge(hexCode string, pinyin string) {
	d.set(hexCode, mergePinyinValues(d.values[hexCode], pinyin))
}

func genCode(data *dictData, outFile *os.File, varName string) {
	output := fmt.Sprintf(`package pinyin

// %s is data map.
// Warning: Auto-generated file, don't edit.
var %s = Dict{
`, varName, varName)
	lines := []string{}

	for _, hexCode := range data.order {
		lines = append(lines, fmt.Sprintf("\t%s: \"%s\",", hexCode, data.values[hexCode]))
	}

	output += strings.Join(lines, "\n")
	output += "\n}\n"
	outFile.WriteString(output)
	return
}

func readDictFile(path string, merge bool, data *dictData) error {
	inFp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inFp.Close()

	scanner := bufio.NewScanner(inFp)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		hexCode, pinyin, ok, err := parseDictLine(scanner.Text())
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if !ok {
			continue
		}
		if merge {
			data.merge(hexCode, pinyin)
		} else {
			data.set(hexCode, pinyin)
		}
	}
	return scanner.Err()
}

func parseDictLine(line string) (hexCode string, pinyin string, ok bool, err error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false, nil
	}

	// line: `U+4E2D: zhōng,zhòng  # 中`
	dataSlice := strings.SplitN(line, "#", 2)
	dataSlice = strings.SplitN(strings.TrimSpace(dataSlice[0]), ":", 2)
	if len(dataSlice) != 2 {
		return "", "", false, fmt.Errorf("invalid dict line %q", line)
	}

	hexCode = strings.TrimSpace(dataSlice[0])
	pinyin = strings.TrimSpace(dataSlice[1])
	if hexCode == "" || pinyin == "" {
		return "", "", false, fmt.Errorf("invalid dict line %q", line)
	}
	hexCode = strings.Replace(hexCode, "U+", "0x", 1)
	return hexCode, pinyin, true, nil
}

func mergePinyinValues(values ...string) string {
	result := make([]string, 0)
	seen := make(map[string]struct{})
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			if _, ok := seen[item]; ok {
				continue
			}
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return strings.Join(result, ",")
}

func buildDict(inputFile string, overlayFiles []string) (*dictData, error) {
	data := newDictData()
	if err := readDictFile(inputFile, false, data); err != nil {
		return nil, err
	}

	if len(overlayFiles) == 0 {
		return data, nil
	}

	overlay := newDictData()
	for _, overlayFile := range overlayFiles {
		if err := readDictFile(overlayFile, true, overlay); err != nil {
			return nil, err
		}
	}
	for _, hexCode := range overlay.order {
		data.set(hexCode, overlay.values[hexCode])
	}
	return data, nil
}

func parseCmdArgs() cmdArgs {
	varName := flag.String("name", "PinyinDict", "generated dictionary variable name")
	flag.Parse()
	args := flag.Args()
	inputFile := ""
	outputFile := ""
	overlayFiles := []string{}
	if len(args) > 0 {
		inputFile = args[0]
	}
	if len(args) > 1 {
		outputFile = args[1]
	}
	if len(args) > 2 {
		overlayFiles = args[2:]
	}
	return cmdArgs{inputFile, outputFile, *varName, overlayFiles}
}

func main() {
	args := parseCmdArgs()
	usage := "gen_pinyin_dict [-name VAR] INPUT OUTPUT [OVERLAY...]"
	inputFile := args.inputFile
	outputFile := args.outputFile
	if inputFile == "" || outputFile == "" || args.varName == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	data, err := buildDict(inputFile, args.overlayFiles)
	if err != nil {
		panic(err)
	}

	outFp, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("open file %s error", outputFile)
		panic(err)
	}
	defer outFp.Close()

	genCode(data, outFp, args.varName)
}
