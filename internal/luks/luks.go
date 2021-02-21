// Copyright 2018 Square Inc.
//
// Use of this source code is governed by a GNU
// General Public License license version 3 that
// can be found in the LICENSE file.

package luks

// #cgo pkg-config: libcryptsetup
// #include <stdlib.h>
// #include <libcryptsetup.h>
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

// Settings holds info about the old and new passwords as well as the luks
// device
type Settings struct {
	OldPass, NewPass, LuksDevice string
	LuksSlot                     int
	cDevice                      *C.struct_crypt_device
	LuksVersion                  int
}

// Error holds error messages from this package
type Error struct {
	code     int
	function string
}

// Error prints error messages from this package
func (e *Error) Error() string {
	return fmt.Sprintf("%s returned error with code %d", e.function, e.code)
}

// cryptInit creates the libcryptsetup struct
func (luksDevice *Settings) cryptInit() (*C.struct_crypt_device, error) {
	cDevice := C.CString(luksDevice.LuksDevice)
	defer C.free(unsafe.Pointer(cDevice))

	var cCD *C.struct_crypt_device

	err := C.crypt_init(&cCD, cDevice)
	if err < 0 {
		return nil, &Error{function: "crypt_init", code: int(err)}
	}

	return cCD, nil
}

// load populates the libcryptsetup struct with device info from disk
func (luksDevice *Settings) load() error {
	cCryptType := C.CString(C.CRYPT_LUKS2)
	if luksDevice.LuksVersion == 1 {
		cCryptType = C.CString(C.CRYPT_LUKS1)
	}

	defer C.free(unsafe.Pointer(cCryptType))

	err := C.crypt_load(luksDevice.cDevice, cCryptType, nil)

	if err < 0 {
		return &Error{function: "crypt_load", code: int(err)}
	}

	return nil
}

// getLuksSlot returns the luks slot number for a pass. Can also be used to
// validate a password
func (luksDevice *Settings) getLuksSlot(pass string) (int, error) {
	cPass := C.CString(pass)
	defer C.free(unsafe.Pointer(cPass))

	cCryptSlot := C.crypt_activate_by_passphrase(
		luksDevice.cDevice,
		nil,
		C.CRYPT_ANY_SLOT,
		cPass,
		C.size_t(len(pass)),
		C.uint32_t(0),
	)
	cryptSlot := int(cCryptSlot)

	if cryptSlot < 0 {
		return -1, &Error{function: "crypt_activate_by_passphrase", code: int(cryptSlot)}
	}

	log.Printf("password is valid and stored in slot %d\n", cryptSlot)

	return cryptSlot, nil
}

// changePassword uses an existing password and updates it to a new password
func (luksDevice *Settings) changePassword() (int, error) {
	var cExistingPass *C.char
	if luksDevice.OldPass == "" {
		cExistingPass = nil
	} else {
		cExistingPass = C.CString(luksDevice.OldPass)
		defer C.free(unsafe.Pointer(cExistingPass))
	}

	cNewPass := C.CString(luksDevice.NewPass)
	defer C.free(unsafe.Pointer(cNewPass))

	log.Printf(
		"updating password on luks device: \"%s\"",
		C.GoString(C.crypt_get_device_name(luksDevice.cDevice)),
	)

	cChangePassRes := C.crypt_keyslot_change_by_passphrase(
		luksDevice.cDevice,
		C.int(luksDevice.LuksSlot),
		C.int(luksDevice.LuksSlot),
		cExistingPass,
		C.size_t(len(luksDevice.OldPass)),
		cNewPass,
		C.size_t(len(luksDevice.NewPass)),
	)
	changePassRes := int(cChangePassRes)
	if changePassRes < 0 {
		return -1, &Error{
			function: "crypt_keyslot_change_by_passphrase",
			code:     int(changePassRes),
		}
	}

	log.Printf("wrote new luks passphrase to slot %d\n", changePassRes)

	return changePassRes, nil
}

// format uses libcryptsetup device format method. This primarily used by tests
// to create a virtual luks disk for testing
func (luksDevice *Settings) format() (int, error) {
	cHash := C.CString("sha256")
	defer C.free(unsafe.Pointer(cHash))

	luksParams := C.struct_crypt_params_luks1{
		hash:           cHash,
		data_alignment: C.size_t(0),
		data_device:    nil,
	}

	cLuksType := C.CString(C.CRYPT_LUKS2)
	if luksDevice.LuksVersion == 1 {
		cLuksType = C.CString(C.CRYPT_LUKS1)
	}

	defer C.free(unsafe.Pointer(cLuksType))

	cLuksCipher := C.CString("aes")
	defer C.free(unsafe.Pointer(cLuksCipher))

	cLuksCipherMode := C.CString("xts-plain64")
	defer C.free(unsafe.Pointer(cLuksCipherMode))

	log.Printf(
		"formating luks device: '%s'",
		C.GoString(C.crypt_get_device_name(luksDevice.cDevice)),
	)
	cFormatRes := C.crypt_format(
		luksDevice.cDevice,
		cLuksType,
		cLuksCipher,
		cLuksCipherMode,
		nil,
		nil,
		C.ulong(256/8),
		unsafe.Pointer(&luksParams),
	)
	formatRes := int(cFormatRes)
	if formatRes < 0 {
		return -1, &Error{
			function: "crypt_format",
			code:     int(formatRes),
		}
	}

	log.Println("formated luks device")
	return 0, nil
}

// addKeyslotByVolumeKey uses the existing luks device context to add a passphrase
// to the next available slot. Used to create test volumes
func (luksDevice *Settings) addKeyslotByVolumeKey() error {
	cPass := C.CString(luksDevice.NewPass)
	defer C.free(unsafe.Pointer(cPass))

	err := C.crypt_keyslot_add_by_volume_key(
		luksDevice.cDevice,
		C.CRYPT_ANY_SLOT,
		nil,
		C.size_t(0),
		cPass,
		C.size_t(len(luksDevice.NewPass)),
	)
	if err < 0 {
		return &Error{
			function: "crypt_keyslot_add_by_volume_key",
			code:     int(err),
		}
	}

	return nil
}

// freeCryptDev release the crypt device and frees memory
func (luksDevice *Settings) freeCryptDev() {
	C.crypt_free(luksDevice.cDevice)
}

// formatSetPassword formats a device with luks and adds a passphrase to device
// This is used by tests to create a virtual disk
func formatSetPassword(pass string, luksDevice string, luksVersion int) (bool, error) {
	cryptInfo := &Settings{
		NewPass:     pass,
		LuksDevice:  luksDevice,
		LuksVersion: luksVersion,
	}

	cCryptDev, err := cryptInfo.cryptInit()
	if err != nil {
		return false, err
	}
	cryptInfo.cDevice = cCryptDev

	cryptInfo.format()
	cryptInfo.addKeyslotByVolumeKey()

	_, err = cryptInfo.getLuksSlot(cryptInfo.NewPass)
	if err != nil {
		return false, err
	}

	cryptInfo.freeCryptDev()

	return true, nil
}

// PassWorks tests if a luks password is correct
func PassWorks(pass string, luksDevice string, luksVersion int) (bool, error) {
	cryptInfo := &Settings{
		LuksDevice:  luksDevice,
		LuksVersion: luksVersion,
	}

	cCryptDev, err := cryptInfo.cryptInit()
	if err != nil {
		return false, err
	}
	cryptInfo.cDevice = cCryptDev

	err = cryptInfo.load()
	if err != nil {
		return false, err
	}

	luksSlot, err := cryptInfo.getLuksSlot(pass)
	if err != nil {
		return false, err
	}
	if luksSlot < 0 {
		return false, &Error{function: "PassWorks", code: luksSlot}
	}

	return true, nil
}

// SetRecoveryPassword changes the luks passphrase on the device
func SetRecoveryPassword(oldPass string, newPass string, luksDevice string, luksVersion int) error {
	cryptInfo := &Settings{
		OldPass:     oldPass,
		NewPass:     newPass,
		LuksDevice:  luksDevice,
		LuksVersion: luksVersion,
	}

	cCryptDev, err := cryptInfo.cryptInit()
	if err != nil {
		return err
	}
	cryptInfo.cDevice = cCryptDev

	err = cryptInfo.load()
	if err != nil {
		return err
	}

	luksSlot, err := cryptInfo.getLuksSlot(cryptInfo.OldPass)
	if err != nil {
		return err
	}
	cryptInfo.LuksSlot = luksSlot

	_, err = cryptInfo.changePassword()
	if err != nil {
		return err
	}

	_, err = cryptInfo.getLuksSlot(cryptInfo.NewPass)
	if err != nil {
		return err
	}

	cryptInfo.freeCryptDev()

	return nil
}
