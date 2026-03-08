package main

type FileMagicNumberMapping struct {
	Name        string
	Signature   []byte
	Offset      int
	CustomCheck func(filePath string) bool
}

var FileType = []FileMagicNumberMapping{
	/* Standard Header */
	{"PNG Image", []byte{0x89, 0x50, 0x4E, 0x47}, 0, nil},
	{"PDF Document", []byte{0x25, 0x50, 0x44, 0x46, 0x2D}, 0, nil},

	/* Uses full 16bit signature */
	{"SQLite Database", []byte{0x53, 0x51, 0x4C, 0x69, 0x74, 0x65, 0x20, 0x66}, 0, nil},

	/*
		Variable signature
		All Scripts starts with '#!'
		To tell the difference between a python3 or bash file is to read the rest of the line
		'#!/usr/bin/python3' or '#!/usr/bin/bash'
	*/
	{"Script (#! environment)", []byte{0x23, 0x21}, 0, ScriptCheck},

	/* Example for larger offset */
	{"ISO Image", []byte{0x43, 0x44, 0x30, 0x30, 0x31}, 32769, nil},

	/*
	 Text - multiple formats
	 Even after that we migth not be able to find the txt file, since there is no one single standard way for handling text
	*/
	// BOMs - Byte order marks
	{"UTF-32LE Text", []byte{0xFF, 0xFE, 0x00, 0x00}, 0, nil},
	{"UTF-32BE Text", []byte{0x00, 0x00, 0xFE, 0xFF}, 0, nil},
	{"UTF-16LE Text", []byte{0xFF, 0xFE}, 0, nil},
	{"UTF-16BE Text", []byte{0xFE, 0xFF}, 0, nil},
	{"UTF-8 Plain Text", []byte{0xEF, 0xBB, 0xBF}, 0, nil},

	// UTF-7 is rare and has multiple variations
	{"UTF-7 Text (v8)", []byte{0x2B, 0x2F, 0x76, 0x38}, 0, nil},
	{"UTF-7 Text (v9)", []byte{0x2B, 0x2F, 0x76, 0x39}, 0, nil},
	{"UTF-7 Text (v+)", []byte{0x2B, 0x2F, 0x76, 0x2B}, 0, nil},
	{"UTF-7 Text (v/)", []byte{0x2B, 0x2F, 0x76, 0x2F}, 0, nil},
}
