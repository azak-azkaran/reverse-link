package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSimlinks(t *testing.T) {
	if _, err := os.Lstat("./random.json"); err == nil {
		os.Remove("./random.json")
	}
	_, err := os.Lstat("./random.json")
	require.Error(t, err)

	if _, err := os.Lstat("./test/random.json"); err != nil {
		d := []byte("{\n\t\"test\": \"testvalue\"\n}")
		err := ioutil.WriteFile("./test/random.json", d, 0644)
		require.NoError(t, err)
	}

	err = os.Symlink("./test/random.json", "./random.json")
	require.NoError(t, err)
}

func TestReadFile(t *testing.T) {
	fmt.Println("testing: TestReadFile")
	createSimlinks(t)

	data, err := ioutil.ReadFile("./random.json")
	require.NoError(t, err)
	require.NotZero(t, data)

	data, err = ioutil.ReadFile("./test/random.json")
	require.NoError(t, err)
	require.NotZero(t, data)

	strings := []string{"--file=./test/random.json"}
	fileName, err := ReadFile(strings)
	assert.Error(t, err)
	assert.Zero(t, fileName)

	strings = []string{"--file=./random.json"}
	fileName, err = ReadFile(strings)
	assert.NoError(t, err)
	require.NotNil(t, fileName)
	assert.Equal(t, fileName, "./random.json")
}

func TestReverseFile(t *testing.T) {
	fmt.Println("testing: TestReverseFile")
	createSimlinks(t)
	require.FileExists(t, "./random.json")
	require.FileExists(t, "./test/random.json")
	Lstat, err := os.Lstat("./random.json")
	require.NoError(t, err)
	require.False(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)

	err = ReverseFile("./random.json")
	assert.NoError(t, err)

	Lstat, err = os.Lstat("./random.json")
	require.NoError(t, err)
	require.True(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)
	_, err = os.Lstat("./test/random.json")
	require.Error(t, err)

	err = os.Remove("./random.json")
	require.NoError(t, err)
}

func TestMain(t *testing.T) {
	fmt.Println("testing: TestMain")
	createSimlinks(t)

	os.Args = append(os.Args, "--file=./random.json")
	main()
	Lstat, err := os.Lstat("./random.json")
	require.NoError(t, err)
	require.True(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)
	_, err = os.Lstat("./test/random.json")
	require.Error(t, err)

	err = os.Remove("./random.json")
	require.NoError(t, err)
}

func TestNoArgs(t *testing.T) {
	fmt.Println("testing: TestNoArgs")
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	main()
}

func TestNoFile(t *testing.T) {
	fmt.Println("testing: TestNoFile")
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	os.Args = append(os.Args, "--file=./random.json")
	main()
}
