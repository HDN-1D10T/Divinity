/*
The included "github.com/google/goexpect" package is used under
the BSD-3-Clause listed below.  All other code falls under the
default Creative Commons License for this project.  BSD-3-Clause
for "github.com/google/goexpect" is as follows:

Copyright (c) 2015 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package tcp

import (
	"fmt"
	"os"
	"regexp"
	"time"

	expect "github.com/google/goexpect"
)

const timeout = 500 * time.Millisecond

func filewrite(chunk, outputFile string) {
	f, err := os.OpenFile(outputFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.WriteString(chunk + "\n"); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

// Telnet - check for valid credentials
func Telnet(ip, user, pass, outputFile string) {
	userRE := regexp.MustCompile(`.*([Ll]ogin)|([Uu]sername).*`)
	passRE := regexp.MustCompile(".*[Pp]assword.*")
	promptRE := regexp.MustCompile(`.*[#\$%>].*`)
	e, _, err := expect.Spawn(fmt.Sprintf("telnet %s", ip), timeout)
	if err != nil {
		return
	}
	defer e.Close()
	e.Expect(userRE, timeout)
	e.Send(user + "\n")
	e.Expect(passRE, timeout)
	e.Send(pass + "\n")
	res, _, err := e.Expect(promptRE, timeout)
	e.Send("exit\n")
	if promptRE.MatchString(res) {
		fmt.Printf("%s:23 %s:%s\n", ip, user, pass)
		if len(outputFile) > 0 {
			filewrite(ip+":23 "+user+":"+pass, outputFile)
		}
	}
}
