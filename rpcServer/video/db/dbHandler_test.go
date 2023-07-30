/*
*

	@author:dailinfeng
	@date:2023/7/29
	@node:

*
*/
package db

import (
	"fmt"
	"testing"
)

func TestFindVideoListBy(t *testing.T) {
	res, _ := FindVideoListBy("id", "1")
	fmt.Printf("%+v", res)
}
