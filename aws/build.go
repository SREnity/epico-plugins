package main

func build() {
    baseFile := "./epico/" + apiType + "/"+ apiType + "_regional.yaml"
    cacheDir := "./epico/" + apiType + "/.cache/"
    rawYaml, err := ioutil.ReadFile(baseFile)
    if err != nil {
        fmt.Printf("--- t dump:\n%s\n\n", string("Error reading YAML API defnition" + err.Error()))
    }

    // Expand out the permutations into files and use them as a cache so this
    //    gets filled once per restart. TODO: This should be moved to a "server"
    //    option - it isn't very good form as a simple library. (Should be kept
    //    in memory and likely purged after use OR moved to build-time?)
    utils.PopulateYamlCache( cacheDir, string(rawYaml), apiType, api.VarsData )
    // TODO: Ugly hack in here for AWS global versus regional -no expansion vars
    globalFile := "./epico/" + apiType + "/"+ apiType + "_global.yaml"
    from, err := os.Open(globalFile)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
    }
    defer from.Close()
    to, err := os.OpenFile("./epico/" + apiType + "/.cache/"+ apiType + "_global.yaml", os.O_RDWR|os.O_CREATE, 0777)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
    }
    defer to.Close()
    _, err = io.Copy(to, from)
    if err != nil {
        fmt.Printf("Error copying file: %v\n", err)
    }
    // END UGLY HACK


}
