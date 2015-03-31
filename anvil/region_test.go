// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package anvil

import (
	"io"
	"os"
	"reflect"
	"testing"
)

type regionCoordTest struct {
	In  string
	X   int
	Z   int
	Err bool
}

func TestRegionCoords(t *testing.T) {
	for _, rct := range []regionCoordTest{
		{In: "", Err: true},
		{In: "r.1.2", Err: true},
		{In: "r.1.2.", X: 1, Z: 2},
		{In: "r.1.2.mca", X: 1, Z: 2},
		{In: "r.1.2.mcs", X: 1, Z: 2},
		{In: "r.-1.2.mca", X: -1, Z: 2},
		{In: "/a/b/r.-1.2.mca", X: -1, Z: 2},
		{In: "a/b/r.-1.2.mca", X: -1, Z: 2},
		{In: "/a/b/x.-1.2.mca", X: -1, Z: 2},
	} {
		testRegionCoords(t, rct)
	}
}

func testRegionCoords(t *testing.T, rc regionCoordTest) {
	x, z, err := RegionCoords(rc.In)
	if err != !rc.Err {
		t.Fatalf("success mismatch: %q\nWant: %v\nHave: %v",
			rc.In, !rc.Err, err)
	}

	if !err {
		return
	}

	if x != rc.X {
		t.Fatalf("X mismatch: %q\nWant: %v\nHave: %v",
			rc.In, rc.X, x)
	}

	if z != rc.Z {
		t.Fatalf("Z mismatch: %q\nWant: %v\nHave: %v",
			rc.In, rc.Z, z)
	}
}

// TestOverwrite ensures that we can create or load a file and overwrite
// existing bytes, as well as append new data to it.
func TestOverwrite(t *testing.T) {
	const File = "../testdata/overwritetest"

	// Create or open initial file for writing.
	fd, err := os.Create(File)
	if err != nil {
		return
	}

	// Write initial contents.
	_, err = fd.Write([]byte{1, 2, 3, 4, 5})
	if err != nil {
		t.Errorf("Write 1: %v", err)
		fd.Close()
		return
	}

	// Close and reopen file with overwritable mode.
	fd.Close()
	fd, err = os.OpenFile(File, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return
	}

	defer fd.Close()

	// Scan to 3rd byte, so we may overwrite some data.
	_, err = fd.Seek(2, 0)
	if err != nil {
		t.Errorf("Seek 2: %v", err)
		return
	}

	// This should replace 03 04 with 06 07
	_, err = fd.Write([]byte{6, 7})
	if err != nil {
		t.Errorf("Write 2: %v", err)
		return
	}

	// Scan to the end so we may append new data.
	_, err = fd.Seek(0, 2)
	if err != nil {
		t.Errorf("Seek 3: %v", err)
		return
	}

	// Adds 08 09 to make: 01 02 06 07 05 08 09
	_, err = fd.Write([]byte{8, 9})
	if err != nil {
		t.Errorf("Write 3: %v", err)
		return
	}

	// Make sure the have the correct final file size of 7 bytes.
	stat, err := fd.Stat()
	if err != nil {
		t.Errorf("Stat 1: %v", err)
		return
	}

	if stat.Size() != 7 {
		t.Errorf("expected file size 7; have %d", stat.Size())
	}
}

func TestRegionRoundtrip(t *testing.T) {
	const File1 = "../testdata/newworld/region/r.0.-1.mca"
	const File2 = "../testdata/newworld/region/r.100.-100.mca"

	if !copyFile(File2, File1) {
		return
	}

	// Load the new copy.
	ra, err := LoadRegion(File2)
	if err != nil {
		t.Errorf("Open 2: %v", err)
		return
	}

	// Don't modify ra's contents, but force re-saving of the data using the
	// region's write routines. We want to know if this works as expected by
	// comparing its contents with those in File1.
	//
	// WriteChunk() implicitely rewrites the region header as well as the
	// chunk data.
	if ra.ChunkLen() > 0 {
		var c Chunk
		xz := ra.Chunks()[0]

		if !ra.ReadChunk(xz[0], xz[1], &c) {
			t.Errorf("ReadChunk failed")
			return
		}

		if !ra.WriteChunk(xz[0], xz[1], &c) {
			t.Errorf("WriteChunk failed")
			return
		}
	}

	err = ra.Save()
	if err != nil {
		t.Errorf("Save: %v", err)
		return
	}

	// Reload ra so we get up-to-date contents. Then compare it to ra.
	ra, err = LoadRegion(File2)
	if err != nil {
		t.Errorf("Open 2: %v", err)
		return
	}

	// Load the original region file.
	rb, err := LoadRegion(File1)
	if err != nil {
		t.Errorf("Open 1: %v", err)
		return
	}

	if ra.ChunkLen() != rb.ChunkLen() {
		t.Errorf("roundtrip mismatch:\nWant: %+v\nHave: %+v", ra, rb)
	}

	if ra.ChunkLen() == 0 {
		return
	}

	var ca, cb Chunk
	xza := ra.Chunks()[0]
	xzb := rb.Chunks()[0]

	if !ra.ReadChunk(xza[0], xza[1], &ca) {
		t.Errorf("ReadChunk ra failed")
		return
	}

	if !rb.ReadChunk(xzb[0], xzb[1], &cb) {
		t.Errorf("ReadChunk rb failed")
		return
	}

	if !reflect.DeepEqual(ca, cb) {
		t.Errorf("roundtrip mismatch:\nWant: %+v\nHave: %+v", ca, cb)
	}
}

// copyFile copies file src to file dst.
func copyFile(dst, src string) bool {
	fs, err := os.Open(src)
	if err != nil {
		return false
	}

	defer fs.Close()

	fd, err := os.Create(dst)
	if err != nil {
		return false
	}

	defer fd.Close()

	_, err = io.Copy(fd, fs)
	return err == nil
}
