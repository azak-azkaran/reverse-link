package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func removeSimlink(){
if _, err := os.Lstat("./random.json"); err == nil {
		os.Remove("./random.json")
	}
if _, err := os.Lstat("./test/random.json"); err == nil {
		os.Remove("./random.json")
	}
}

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

func TestReadFlags(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestReadFlags")

	args := []string{"--file=./test/random.json"}
	fileName, err := ReadFlags(args)
	assert.NoError(t, err)
	assert.Equal(t, fileName, "./test/random.json")

	args = []string{"--file=./random.json"}
	fileName, err = ReadFlags(args)
	assert.NoError(t, err)
	assert.Equal(t, fileName, "./random.json")

	args = []string{"--file="}
	fileName, err = ReadFlags(args)
	assert.Error(t, err)
	assert.Equal(t, fileName, "")

	args = []string{"reverse-link","random.json"}
	fileName, err = ReadFlags(args)
	assert.NoError(t, err)
	assert.Equal(t, fileName, "random.json")
}


func TestCheckFile(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestCheckFile")
	createSimlinks(t)

	data, err := ioutil.ReadFile("./random.json")
	require.NoError(t, err)
	require.NotZero(t, data)

	data, err = ioutil.ReadFile("./test/random.json")
	require.NoError(t, err)
	require.NotZero(t, data)

	fileName, err := CheckFile("./test/random.json")
	assert.Error(t, err)
	assert.Zero(t, fileName)

	fileName, err = CheckFile("./random.json")
	assert.NoError(t, err)
	require.NotNil(t, fileName)
	assert.Equal(t, fileName, "./random.json")
}

func TestReverseFile(t *testing.T) {
	t.Cleanup(removeSimlink)
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

}

func TestReverseNoFile(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestReverseNoFile")
	createSimlinks(t)
	os.Remove("./test/random.json")
	require.FileExists(t, "./random.json")

	Lstat, err := os.Lstat("./random.json")
	require.NoError(t, err)
	require.False(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)

	err = ReverseFile("./random.json")
	assert.Error(t, err)

	Lstat, err = os.Lstat("./random.json")
	require.NoError(t, err)
	require.False(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)
}

func TestMain(t *testing.T) {
	fmt.Println("testing: TestMain")
	createSimlinks(t)

	os.Args = append(os.Args, "--file=./random.json")
	main()

	os.Args = os.Args[:len(os.Args)-1]
	Lstat, err := os.Lstat("./random.json")
	require.NoError(t, err)
	require.True(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)
	_, err = os.Lstat("./test/random.json")
	require.Error(t, err)

	createSimlinks(t)
	os.Args = append(os.Args, "./random.json")
	main()

	os.Args = os.Args[:len(os.Args)-1]
	Lstat, err = os.Lstat("./random.json")
	require.NoError(t, err)
	require.True(t, Lstat.Mode()&os.ModeSymlink != os.ModeSymlink)
	_, err = os.Lstat("./test/random.json")
	require.Error(t, err)

}


func TestEmptyFlag(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestEmptyFlag")
	defer func() {
		err := recover()
		require.NotNil(t, err)
		os.Args = os.Args[:len(os.Args)-1]
	}()
	os.Args = append(os.Args, "--file=")
	main()
}

func TestNoArgs(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestNoArgs")
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	main()
}

func TestNoFile(t *testing.T) {
	t.Cleanup(removeSimlink)
	fmt.Println("testing: TestNoFile")
	defer func() {
		err := recover()
		require.NotNil(t, err)
		os.Args = os.Args[:len(os.Args)-1]
	}()
	os.Args = append(os.Args, "--file=./random.json")
	main()
}
