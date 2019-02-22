package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
)

const badEnvVariable = "BAD_ENV"

func convertMap(envMap map[string]string) []string {
	retVal := []string{}
	for key, val := range envMap {
		retVal = append(retVal, key+"="+val)
	}
	return retVal
}

func getEnv() []string {
	envMap := map[string]string{}
	for _, variable := range os.Environ() {
		kv := strings.SplitN(variable, "=", 2)
		envMap[kv[0]] = kv[1]
	}

	val, ok := os.LookupEnv(badEnvVariable)
	if !ok {
		return convertMap(envMap)
	}

	decodedVal, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		panic(err)
	}

	tmpEnvMap := map[string]string{}
	err = json.Unmarshal(decodedVal, &tmpEnvMap)
	if err != nil {
		panic(err)
	}
	for key, val := range tmpEnvMap {
		envMap[key] = val
	}

	return convertMap(envMap)
}

func main() {
	fmt.Println(os.Args[1:])
	variables := getEnv()
	err := syscall.Exec(os.Args[1], os.Args[1:], variables)
	if err != nil {
		panic(err)
	}
}
