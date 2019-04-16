package main

// essay generator using n-gram markov chains
// NOTE: dumb implementation with no map optimizations

import (
    "fmt"
    "io"
    "os"
    "bufio"
    "strconv"
    "math/rand"
)

// TYPES -- there are no classes in go, only types with methods
type Gram struct {
    gram string
    // p float32 // probability of gram FIXME: implement this optimization
}

func (self *Gram) String() string{
    return fmt.Sprintf(self.gram)
}


type MarkovMap struct {
    gramMap map[string][]Gram
    n int
}

func (self *MarkovMap) String() string{
    return fmt.Sprintf("Map: %v\n", self.gramMap)
}

func (self *MarkovMap) Initialize(filename string) {
    // open file and check for failure
    file, err := os.Open(filename)
    check(err)
    defer file.Close() // defer close to clean up code + increase locality

    reader := bufio.NewReader(file)
    firstNGram := make([]byte, self.n)
    _, err = reader.Read(firstNGram)
    if err != nil {
        fmt.Println("File is too small for that n-gram")
        os.Exit(1)
    }

    // check if n is length of file and increments reader another NGram
    secondNGram := make([]byte, self.n)
    _, err = reader.Read(secondNGram)
    if err == io.EOF {
        fmt.Println("File is too small for that n-gram")
    }

    // stores the previous nGram (key in markov map)
    nGram, gram := firstNGram, secondNGram
    buffer := make([]byte, 1)
    for { // infinite loop in go
        // SAVE
        key := string(nGram)
        value := Gram{string(gram)}
        if elem, ok := self.gramMap[key]; ok {
            self.gramMap[key] = append(elem, value)
        } else {
            self.gramMap[key] = []Gram{value}
        }

        // READ
        _, err = reader.Read(buffer)

        // CHECK
        // deal with when the reader reads exactly all the file
        if _, err :=  reader.Peek(self.n); err != nil {
             // breaking because there is already enough data
             // -- not technically reading the whole file
             break
        }
        if err == io.EOF {
            break
        }

        // SET
        nGram = append(nGram[1:], gram[0])
        gram = append(gram[1:], buffer[0])
    }
}

// generate 'length' n-grams text + save to file 'filename'
func (self *MarkovMap) Generate(filename string, length int) {
    // open file and check for failure
    file, err := os.Create(filename)
    check(err)
    defer file.Close() // defer close to clean up code + increase locality

    // randomly choose starting n-gram from maps' keys
    keys := make([]string, 0, len(self.gramMap))
    for k := range self.gramMap {
        keys = append(keys, k)    
    }
    key := keys[rand.Intn(len(keys))]

    for i := 0; i < length; i++ {
        // write key to file
        file.WriteString(key)
        file.Sync()

        // find new key
        values := self.gramMap[key]
        key = values[rand.Intn(len(values))].gram
    }
}

// HELPER FUNCTIONS
// catch-all error checker
func check(err error) {
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }
}


// ENTRY POINT
func main() {
    // get cl arguments
    if len(os.Args[1:]) < 1 || len(os.Args[1:]) > 4 {
        fmt.Println("4 arguments required!")
        fmt.Println("\t1: input filename")
        fmt.Println("\t2: length of n-gram")
        fmt.Println("\t3: output filename")
        fmt.Println("\t4: number of n-grams to generate")
        os.Exit(1)
    }
    n, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("length of n-gram must be an int")
        os.Exit(1)
    }
    length, err := strconv.Atoi(os.Args[4])
    if err != nil {
        fmt.Println("number of n-grams generated must be an int")
        os.Exit(1)
    }

    markovMap := MarkovMap{make(map[string][]Gram), n}
    markovMap.Initialize(os.Args[1]) // populate markov map with data from file
    //fmt.Println(markovMap)
    markovMap.Generate(os.Args[3], length)
}
