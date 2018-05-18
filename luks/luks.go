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

// LuksInfo holds info about the old and new passwords as well as the luks
// device
type LuksInfo struct {
	OldPass, NewPass, LuksDevice string
	LuksSlot                     int
	cDevice                      *C.struct_crypt_device
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
func (cryptInfo *LuksInfo) cryptInit() (*C.struct_crypt_device, error) {
	cDevice := C.CString(cryptInfo.LuksDevice)
	defer C.free(unsafe.Pointer(cDevice))

	var cCD *C.struct_crypt_device

	err := C.crypt_init(&cCD, cDevice)
	if err < 0 {
		return nil, &Error{function: "crypt_init", code: int(err)}
	}

	return cCD, nil
}

// load populates the libcryptsetup struct with device info from disk
func (luksDevice *LuksInfo) load() error {
	cCryptType := C.CString(C.CRYPT_LUKS1)
	defer C.free(unsafe.Pointer(cCryptType))

	err := C.crypt_load(luksDevice.cDevice, cCryptType, nil)

	if err < 0 {
		return &Error{function: "crypt_load", code: int(err)}
	}

	return nil
}

// getLuksSlot returns the luks slot number for a pass. Can also be used to
// validate a password
func (luksDevice *LuksInfo) getLuksSlot(pass string) (int, error) {
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
func (luksDevice *LuksInfo) changePassword() (int, error) {
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

// freeCryptDev release the crypt device and frees memory
func (luksDevice *LuksInfo) freeCryptDev() {
	C.crypt_free(luksDevice.cDevice)
}

// PassWorks tests if a luks password is correct
func PassWorks(pass string, luksDevice string) (bool, error) {
	cryptInfo := &LuksInfo{
		LuksDevice: luksDevice,
	}

	cCryptDev, initErr := cryptInfo.cryptInit()
	if initErr != nil {
		return false, initErr
	}
	cryptInfo.cDevice = cCryptDev

	loadErr := cryptInfo.load()
	if loadErr != nil {
		return false, loadErr
	}

	luksSlot, luksSlotErr := cryptInfo.getLuksSlot(pass)
	if luksSlotErr != nil {
		return false, luksSlotErr
	}
	if luksSlot < 0 {
		return false, &Error{function: "PassWorks", code: luksSlot}
	}

	return true, nil
}

// SetRecoveryPassword changes the luks passphrase on the device
func SetRecoveryPassword(
	oldPass string,
	newPass string,
	luksDevice string,
) error {
	cryptInfo := &LuksInfo{
		OldPass:    oldPass,
		NewPass:    newPass,
		LuksDevice: luksDevice,
	}

	cCryptDev, initErr := cryptInfo.cryptInit()
	if initErr != nil {
		return initErr
	}
	cryptInfo.cDevice = cCryptDev

	loadErr := cryptInfo.load()
	if loadErr != nil {
		return loadErr
	}

	luksSlot, luksSlotErr := cryptInfo.getLuksSlot(cryptInfo.OldPass)
	if luksSlotErr != nil {
		return luksSlotErr
	}
	cryptInfo.LuksSlot = luksSlot

	_, setErr := cryptInfo.changePassword()
	if setErr != nil {
		return setErr
	}

	_, validLuksSlotErr := cryptInfo.getLuksSlot(cryptInfo.NewPass)
	if validLuksSlotErr != nil {
		return luksSlotErr
	}

	cryptInfo.freeCryptDev()

	return nil
}
