// Copyright [2019] [Mark Farnan]

//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at

//        http://www.apache.org/licenses/LICENSE-2.0

//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

//  NOTICE:  This code was origionally based on sample code from https://github.com/stdiopt/gowasm-experiments though now mostly re-written
//  Along with other samples from the Go WASM Wiki https://github.com/golang/go/wiki/WebAssembly
//  It is reused under the Apache 2.0 Licence

package canvas

import (
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
)

type FontCache map[string]*truetype.Font

func (f FontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	font, ok := f[fd.Name]
	if !ok {
		return f["roboto"], nil
	}
	return font, nil
}

func (f *FontCache) Store(fd draw2d.FontData, tf *truetype.Font) {
	(*f)[fd.Name] = tf
}
