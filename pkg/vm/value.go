package vm

// ----- ---- --- -- -
// Copyright 2019 Oneiro NA, Inc. All Rights Reserved.
//
// Licensed under the Apache License 2.0 (the "License").  You may not use
// this file except in compliance with the License.  You can obtain a copy
// in the file LICENSE in the source distribution or at
// https://www.apache.org/licenses/LICENSE-2.0.txt
// - -- --- ---- -----

// Value objects are what is managed by the VM
type Value interface {
	Equal(rhs Value) bool
	Less(rhs Value) (bool, error)
	IsScalar() bool
	String() string
	IsTrue() bool
}
