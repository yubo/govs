# Copyright 2016,2017 Xiaomi Corporation. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
# 
# Authors:    Yu Bo <yubo@xiaomi.com>
# 
all: govs

govs: *.go cmd/govs/*.go
	go build -o $@ cmd/govs/*.go
