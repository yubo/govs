/*
 * Copyright 2016 Xiaomi Corporation. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 *
 * Authors:    Yu Bo <yubo@xiaomi.com>
 */
package govs

import "fmt"

const (
	EPERM   = 1  /* Operation not permitted */
	ENOENT  = 2  /* No such file or directory */
	ESRCH   = 3  /* No such process */
	EINTR   = 4  /* Interrupted system call */
	EIO     = 5  /* I/O error */
	ENXIO   = 6  /* No such device or address */
	E2BIG   = 7  /* Argument list too long */
	ENOEXEC = 8  /* Exec format error */
	EBADF   = 9  /* Bad file number */
	ECHILD  = 10 /* No child processes */
	EAGAIN  = 11 /* Try again */
	ENOMEM  = 12 /* Out of memory */
	EACCES  = 13 /* Permission denied */
	EFAULT  = 14 /* Bad address */
	ENOTBLK = 15 /* Block device required */
	EBUSY   = 16 /* Device or resource busy */
	EEXIST  = 17 /* File exists */
	EXDEV   = 18 /* Cross-device link */
	ENODEV  = 19 /* No such device */
	ENOTDIR = 20 /* Not a directory */
	EISDIR  = 21 /* Is a directory */
	EINVAL  = 22 /* Invalid argument */
	ENFILE  = 23 /* File table overflow */
	EMFILE  = 24 /* Too many open files */
	ENOTTY  = 25 /* Not a typewriter */
	ETXTBSY = 26 /* Text file busy */
	EFBIG   = 27 /* File too large */
	ENOSPC  = 28 /* No space left on device */
	ESPIPE  = 29 /* Illegal seek */
	EROFS   = 30 /* Read-only file system */
	EMLINK  = 31 /* Too many links */
	EPIPE   = 32 /* Broken pipe */
	EDOM    = 33 /* Math argument out of domain of func */
	ERANGE  = 34 /* Math result not representable */
	EMAX    = 34
)

var e2s = []string{"0",
	"EPERM",
	"ENOENT",
	"ESRCH",
	"EINTR",
	"EIO",
	"ENXIO",
	"E2BIG",
	"ENOEXEC",
	"EBADF",
	"ECHILD",
	"EAGAIN",
	"ENOMEM",
	"EACCES",
	"EFAULT",
	"ENOTBLK",
	"EBUSY",
	"EEXIST",
	"EXDEV",
	"ENODEV",
	"ENOTDIR",
	"EISDIR",
	"EINVAL",
	"ENFILE",
	"EMFILE",
	"ENOTTY",
	"ETXTBSY",
	"EFBIG",
	"ENOSPC",
	"ESPIPE",
	"EROFS",
	"EMLINK",
	"EPIPE",
	"EDOM",
	"ERANGE",
	"EMAX"}

type Ecode int

func (e Ecode) String() string {
	i := int(e)
	if i < 0 {
		i = -i
	}
	if i > EMAX {
		return fmt.Sprintf("%d", i)
	}
	return e2s[i]
}
